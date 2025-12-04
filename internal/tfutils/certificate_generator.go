package tfutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"

	"strings"
	"time"
)

var (
	pemFile = "cert.pem"
	keyFile = "key.pem"
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
