package tfutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"software.sslmate.com/src/go-pkcs12"

	"strings"
	"time"
)

var (
	pemFile = "cert.pem"
	keyFile = "key.pem"
	p12File = "cert.p12"
)

func GeneratePEMCertificate() error {

	// Generate a new RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("unable to generate RSA key : %v", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128) // 128-bit serial number
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	// Create a template for the certificate
	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,

		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0), // 10 years validity
		Subject: pkix.Name{
			Organization: []string{"My Company"},
		},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create the self-signed certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("unable to create self-signed certificate : %v", err)
	}

	certPEM := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}

	privateKeyPEM := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	fileTemplateMap := map[string]pem.Block{
		pemFile: certPEM,
		keyFile: privateKeyPEM,
	}

	for file, pemBlock := range fileTemplateMap {

		certFile, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("unable to write to file %s : %v ", file, err)
		}
		defer func() {
			if tempErr := certFile.Close(); tempErr != nil {
				err = tempErr
			}
		}()

		if err := pem.Encode(certFile, &pemBlock); err != nil {
			return fmt.Errorf("unable to encode pem file : %v", err)
		}

	}

	return nil
}

func ReadPEMCertificate() (string, error) {
	if err := GeneratePEMCertificate(); err != nil {
		return "", fmt.Errorf("unable to generate pem certificate and key : %v", err)
	}

	data, err := os.ReadFile(pemFile)
	if err != nil {
		return "", fmt.Errorf("unable to read certificate file : %s :%v", pemFile, err)
	}

	// Convert the byte slice to a string and replace all the newline characters with the escaped newline characters
	pemString := strings.ReplaceAll(string(data), "\n", "\\n")

	// Delete the contents of the file
	err = os.Remove(pemFile)
	if err != nil {
		return "", fmt.Errorf("unable to delete file : %s : %v", pemFile, err)
	}

	err = os.Remove(keyFile)
	if err != nil {
		return "", fmt.Errorf("unable to delete file : %s : %v", keyFile, err)
	}

	return pemString, nil
}

func GenerateP12FromPEM(password string) error {
	// Read certificate PEM file
	certPEMData, err := os.ReadFile(pemFile)
	if err != nil {
		return fmt.Errorf("unable to read certificate file: %v", err)
	}

	// Read private key PEM file
	keyPEMData, err := os.ReadFile(keyFile)
	if err != nil {
		return fmt.Errorf("unable to read key file: %v", err)
	}

	// Decode certificate PEM
	certBlock, _ := pem.Decode(certPEMData)
	if certBlock == nil {
		return fmt.Errorf("unable to decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return fmt.Errorf("unable to parse certificate: %v", err)
	}

	// Decode private key PEM
	keyBlock, _ := pem.Decode(keyPEMData)
	if keyBlock == nil {
		return fmt.Errorf("unable to decode private key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("unable to parse private key: %v", err)
	}

	// Encode to P12 format
	p12Data, err := pkcs12.LegacyRC2.Encode(privateKey, cert, nil, password)
	if err != nil {
		return fmt.Errorf("unable to encode P12: %v", err)
	}

	// Write P12 file
	if err = os.WriteFile(p12File, p12Data, 0600); err != nil {
		return fmt.Errorf("unable to write P12 file: %v", err)
	}

	return nil
}

func GetBase64EncodedCertificate(certType string) (string, error) {

	var files []string

	if err := GeneratePEMCertificate(); err != nil {
		return "", err
	}

	switch certType {

	case "pem":
		files = []string{pemFile, keyFile}

	case "pfx":
		fallthrough
	case "p12":
		p12Password := ""

		if err := GenerateP12FromPEM(p12Password); err != nil {
			return "", fmt.Errorf("unable to generate p12 file from pem certificate : %v", err)
		}

		files = []string{p12File, pemFile, keyFile}

	default:

		certTypeErr := fmt.Errorf("unsupported certificate type: %s", certType)

		if err := os.Remove(pemFile); err != nil {
			return "", fmt.Errorf("%v\n unable to delete file : %s : %v", certTypeErr, pemFile, err)
		}

		if err := os.Remove(keyFile); err != nil {
			return "", fmt.Errorf("%v\n unable to delete file : %s : %v", certTypeErr, keyFile, err)
		}

		return "", certTypeErr
	}

	data, err := os.ReadFile(files[0])
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(string(data)))

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return encoded, fmt.Errorf("unable to delete file : %s : %v", file, err)
		}
	}

	return encoded, nil
}
