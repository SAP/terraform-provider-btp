package btpclisession

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	keychainService    = "cli.btp.cloud.sap:session"
	defaultConfigFile  = "config.json"
	cachedSessionFile  = "session"
	secureStoreKey     = "login.securestore" // persistence key in config.json Settings map
	secureStoreDefault = true                // CLI default is true
)

// Relevant fields from the BTP CLI config.json.
type Config struct {
	ServerURL      string `json:"ServerURL"`
	UserName       string `json:"UserName"`
	Authentication struct {
		Mail   string `json:"Mail"`
		Issuer string `json:"Issuer"`
	} `json:"Authentication"`
	Settings map[string]string `json:"Settings"`
}

type Result struct {
	SessionID  string
	SessionKey string // SHA-256 hash of the config path; used as keychain account key
	Config     Config
	ConfigPath string
	Source     string // "keychain" | "file" | "not found"
}

// ReadSession loads the BTP CLI config from configPath (pass "" for the
// default location) and returns the session ID together with the parsed config.
//
// configPath mirrors the --config flag of the btp CLI: it may be:
//   - empty          → use the platform default directory
//   - a directory    → use config.json inside that directory
//   - a file path    → use that exact file
func ReadSession(configPath string) (Result, error) {
	configPath, err := resolveConfigPath(configPath)
	if err != nil {
		return Result{}, fmt.Errorf("resolve config path: %w", err)
	}

	cfg, err := readConfig(configPath)
	if err != nil {
		return Result{}, fmt.Errorf("read config %s: %w", configPath, err)
	}

	hash, err := hashOfPath(configPath)
	if err != nil {
		return Result{}, fmt.Errorf("hash config path: %w", err)
	}

	var sessionID, source string

	if useSecureStore(cfg) {
		sessionID, err = readFromKeychain(hash)
		if err != nil {
			// keychain unavailable – file fallback
			sessionID, err = readFromFile(hash)
			if err != nil {
				return Result{}, fmt.Errorf("read session from file (keychain fallback): %w", err)
			}
			source = "file (keychain fallback)"
		} else {
			source = "keychain"
		}
	} else {
		sessionID, err = readFromFile(hash)
		if err != nil {
			return Result{}, fmt.Errorf("read session from file: %w", err)
		}
		source = "file"
	}

	if sessionID == "" {
		source = "not found"
	}

	return Result{
		SessionID:  sessionID,
		SessionKey: hash,
		Config:     cfg,
		ConfigPath: configPath,
		Source:     source,
	}, nil
}

func resolveConfigPath(configPath string) (string, error) {
	if configPath == "" {
		dir, err := defaultConfigDirectory()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, defaultConfigFile), nil
	}

	info, statErr := os.Stat(configPath)
	if statErr == nil && info.IsDir() {
		return filepath.Join(filepath.Clean(configPath), defaultConfigFile), nil
	}

	// treat as file path (may or may not exist yet)
	last := configPath[len(configPath)-1]
	if last == os.PathSeparator {
		return filepath.Join(filepath.Clean(configPath), defaultConfigFile), nil
	}
	return configPath, nil
}

func defaultConfigDirectory() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", errors.New("cannot determine user config directory")
	}
	return joinWithBTPSubDir(userConfigDir), nil
}

func joinWithBTPSubDir(base string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(base, "SAP", "btp")
	}
	return filepath.Join(base, ".btp")
}

func readConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config JSON: %w", err)
	}
	return cfg, nil
}

func hashOfPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(abs))
	return hex.EncodeToString(sum[:16]), nil
}

func useSecureStore(cfg Config) bool {
	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		return false // Linux: keychain not supported
	}
	val, ok := cfg.Settings[secureStoreKey]
	if !ok || val == "default" {
		return secureStoreDefault // respect the CLI default (true)
	}
	return val == "true"
}

func readFromKeychain(accountKey string) (string, error) {
	sessionID, err := keyring.Get(keychainService, accountKey)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", nil
	}
	return sessionID, err
}

func readFromFile(hash string) (string, error) {
	path, err := cachedSessionPath(hash)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func cachedSessionPath(hash string) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", errors.New("cannot determine user cache directory")
	}
	return filepath.Join(joinWithBTPSubDir(cacheDir), hash, cachedSessionFile), nil
}
