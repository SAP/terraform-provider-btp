package tfutils

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateCertificate(t *testing.T) {

	t.Run("Generate Certificate - Should Create Certificate File", func(t *testing.T) {

		err := GenerateCertificate()
		if err != nil {
			t.Fatalf("Generate Certificate failed: %v", err)
		}

		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Fatalf("Certificate file %s not created", file)
		}

		// Clean up the file after test
		err = os.Remove(file)
		if err != nil {
			t.Fatalf("Failed to remove certificate file after test: %v", err)
		}
	})
}

func TestReadCertificate(t *testing.T) {
	t.Run("Read Certificate - Should Return Valid Pem String", func(t *testing.T) {

		pemString, err := ReadCertificate()
		if err != nil {
			t.Fatalf("Read Certificate failed: %v", err)
		}

		// Check if the pemString is not empty
		if pemString == "" {
			t.Fatalf("Read Certificate returned an empty PEM string")
		}

		// Check if the string starts with the PEM certificate header
		if !strings.HasPrefix(pemString, "-----BEGIN CERTIFICATE-----") {
			t.Fatalf("PEM string does not have the correct header")
		}

		// Check if the string ends with the PEM certificate footer
		if !strings.HasSuffix(pemString, "-----END CERTIFICATE-----\\n") {
			t.Fatalf("PEM string does not have the correct footer")
		}
	})
}
