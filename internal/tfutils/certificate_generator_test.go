package tfutils

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCertificate(t *testing.T) {
	t.Run("Generate PEM Certificate - Should Create PEM Certificate File and Key", func(t *testing.T) {

		err := GeneratePEMCertificate()

		assert.NoError(t, err)
		assert.FileExists(t, pemFile)

		// Clean up the file after test
		err = os.Remove(pemFile)
		if err != nil {
			t.Fatalf("Failed to remove certificate file after test: %v", err)
		}

		err = os.Remove(keyFile)
		if err != nil {
			t.Fatalf("Failed to remove key file after test: %v", err)
		}
	})
}

func TestReadCertificate(t *testing.T) {
	t.Run("Read PEM Certificate - Should Return Valid Pem String", func(t *testing.T) {

		pemString, err := ReadPEMCertificate()

		assert.NoError(t, err)
		assert.NotEqual(t, len(pemString), 0)
		assert.Equal(t, strings.HasPrefix(pemString, "-----BEGIN CERTIFICATE-----"), true)
		assert.Equal(t, strings.HasSuffix(pemString, "-----END CERTIFICATE-----\\n"), true)

	})
}
