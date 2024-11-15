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

var file = "cert.pem" 

func GenerateCertificate() error {

	// Generate a new RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
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
		return err
	}

	// Save the certificate to a file
	certFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer certFile.Close()

	certPEM := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}
	if err := pem.Encode(certFile, &certPEM); err != nil {
		return err
	}

	return nil
}

func ReadCertificate() (string, error) {
	err := GenerateCertificate()

	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a string and replace all the newline characters with the escaped newline characters
	pemString := strings.ReplaceAll(string(data), "\n", "\\n")

	// Delete the contents of the file
	err = os.Remove(file)
	if err != nil {
		fmt.Println("Unable to delete PEM file")
	}

	return pemString, nil
}
