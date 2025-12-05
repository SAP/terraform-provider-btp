package tfutils

import (
	"fmt"
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

func TestGenerateP12FromPEM(t *testing.T) {

	t.Run("Generate P12 Certificate from PEM - Should create P12 Certificate File", func(t *testing.T) {

		if err := GeneratePEMCertificate(); err != nil {
			t.Fatalf("Generate Certificate failed: %v", err)
		}

		err := GenerateP12FromPEM("")

		assert.NoError(t, err)
		assert.FileExists(t, p12File)

		// Clean up the file after test
		err = os.Remove(pemFile)
		if err != nil {
			t.Fatalf("Failed to remove certificate file after test: %v", err)
		}

		err = os.Remove(p12File)
		if err != nil {
			t.Fatalf("Failed to remove certificate file after test: %v", err)
		}

		err = os.Remove(keyFile)
		if err != nil {
			t.Fatalf("Failed to remove key file after test: %v", err)
		}

	})
}

func TestGetBase64EncodedCertificate(t *testing.T) {

	certTypes := []string{"pem", "pfx", "p12", "invalid-type"}

	tests := []struct {
		certType      string
		encodedString bool
		err           error
	}{
		{
			certType:      certTypes[0],
			encodedString: true,
			err:           nil,
		},
		{
			certType:      certTypes[1],
			encodedString: true,
			err:           nil,
		},
		{
			certType:      certTypes[2],
			encodedString: true,
			err:           nil,
		},
		{
			certType:      certTypes[3],
			encodedString: false,
			err:           fmt.Errorf("unsupported certificate type: %s", certTypes[3]),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Generate Base64 Encoded Certificate. for type %s", test.certType), func(t *testing.T) {
			encoded, err := GetBase64EncodedCertificate(test.certType)

			if test.encodedString {
				assert.NotEqual(t, len(encoded), 0)
			} else {
				assert.Equal(t, len(encoded), 0)
			}

			assert.Equal(t, err, test.err)
		})
	}

}
