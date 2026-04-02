package btpclisession

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// joinWithBTPSubDir
// ---------------------------------------------------------------------------

func TestJoinWithBTPSubDir(t *testing.T) {
	t.Parallel()

	base := "/home/user/.config"

	result := joinWithBTPSubDir(base)

	if runtime.GOOS == "windows" {
		want := filepath.Join(base, "SAP", "btp")
		if result != want {
			t.Errorf("windows: got %q, want %q", result, want)
		}
	} else {
		want := filepath.Join(base, ".btp")
		if result != want {
			t.Errorf("non-windows: got %q, want %q", result, want)
		}
	}
}

func TestJoinWithBTPSubDir_EmptyBase(t *testing.T) {
	t.Parallel()

	result := joinWithBTPSubDir("")

	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(result, filepath.Join("SAP", "btp")) {
			t.Errorf("windows: expected suffix SAP/btp, got %q", result)
		}
	} else {
		if result != ".btp" {
			t.Errorf("non-windows: got %q, want %q", result, ".btp")
		}
	}
}

// ---------------------------------------------------------------------------
// useSecureStore
// ---------------------------------------------------------------------------

func TestUseSecureStore_Linux(t *testing.T) {
	t.Parallel()

	if runtime.GOOS != "linux" {
		t.Skip("only relevant on linux")
	}

	cfg := Config{}
	if useSecureStore(cfg) {
		t.Error("expected false on linux, got true")
	}
}

func TestUseSecureStore_WindowsDarwin(t *testing.T) {
	t.Parallel()

	if runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		t.Skip("only relevant on windows/darwin")
	}

	tests := []struct {
		name     string
		settings map[string]string
		want     bool
	}{
		{
			name:     "no settings key – use default (true)",
			settings: nil,
			want:     secureStoreDefault,
		},
		{
			name:     "settings key absent – use default (true)",
			settings: map[string]string{"other": "value"},
			want:     secureStoreDefault,
		},
		{
			name:     "settings key = default – use default (true)",
			settings: map[string]string{secureStoreKey: "default"},
			want:     secureStoreDefault,
		},
		{
			name:     "settings key = true",
			settings: map[string]string{secureStoreKey: "true"},
			want:     true,
		},
		{
			name:     "settings key = false",
			settings: map[string]string{secureStoreKey: "false"},
			want:     false,
		},
		{
			name:     "settings key = unrecognised value",
			settings: map[string]string{secureStoreKey: "yes"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg := Config{Settings: tt.settings}
			got := useSecureStore(cfg)
			if got != tt.want {
				t.Errorf("useSecureStore(%v) = %v, want %v", tt.settings, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// hashOfPath
// ---------------------------------------------------------------------------

func TestHashOfPath_DeterministicForAbsolutePath(t *testing.T) {
	t.Parallel()

	// Use an absolute path so filepath.Abs is a no-op and no I/O occurs.
	path := "/absolute/path/to/config.json"

	h1, err := hashOfPath(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h2, err := hashOfPath(path)
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}

	if h1 != h2 {
		t.Errorf("hashOfPath is not deterministic: %q != %q", h1, h2)
	}
}

func TestHashOfPath_DifferentPathsProduceDifferentHashes(t *testing.T) {
	t.Parallel()

	h1, err := hashOfPath("/path/one/config.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h2, err := hashOfPath("/path/two/config.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h1 == h2 {
		t.Error("different paths produced identical hashes")
	}
}

func TestHashOfPath_IsHexString(t *testing.T) {
	t.Parallel()

	hash, err := hashOfPath("/some/absolute/config.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// sha256 truncated to 16 bytes → 32 hex chars
	const wantLen = 32
	if len(hash) != wantLen {
		t.Errorf("hash length = %d, want %d; hash = %q", len(hash), wantLen, hash)
	}
	for _, c := range hash {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Errorf("hash contains non-hex character %q; hash = %q", c, hash)
		}
	}
}

// ---------------------------------------------------------------------------
// resolveConfigPath
// ---------------------------------------------------------------------------

func TestResolveConfigPath_NonExistentFilePath(t *testing.T) {
	t.Parallel()

	// A non-existent file path (no trailing separator): os.Stat will fail,
	// so the function returns the path unchanged.
	input := "/does/not/exist/config.json"
	got, err := resolveConfigPath(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != input {
		t.Errorf("got %q, want %q", got, input)
	}
}

func TestResolveConfigPath_TrailingPathSeparator(t *testing.T) {
	t.Parallel()

	// A path ending with the OS separator should append config.json.
	input := "/does/not/exist/somedir" + string(filepath.Separator)
	got, err := resolveConfigPath(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := filepath.Join("/does/not/exist/somedir", defaultConfigFile)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
