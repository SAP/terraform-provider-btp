package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceDirectoryApiCredential(t *testing.T) {
	t.Parallel()

	t.Run("happy path - api-credential with client secret", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_api_credential.with_secret")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryApiCredential("uut", "directory-api-credential-with-secret", "test_dir", false),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "name", "directory-api-credential-with-secret"),
						resource.TestMatchResourceAttr("btp_directory_api_credential.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "credential_type", "secret"),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "read_only", "false"),
					),
				},
			},
		})
	})

	t.Run("happy path - api-credential with certificate", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_api_credential.with_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryApiCredentialWithCertificate("uut", "directory-api-credential-with-certificate", "test_dir", rec.IsRecording()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "name", "directory-api-credential-with-certificate"),
						resource.TestMatchResourceAttr("btp_directory_api_credential.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "credential_type", "client certificate"),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "read_only", "false"),
					),
				},
			},
		})
	})

	t.Run("happy path - api-credential with read-only set to true", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_api_credential.read_only_credentials")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryApiCredential("uut", "directory-api-credential-read-only", "test_dir", true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "name", "directory-api-credential-read-only"),
						resource.TestMatchResourceAttr("btp_directory_api_credential.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "credential_type", "secret"),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "read_only", "true"),
					),
				},
			},
		})
	})

	t.Run("error path - invalid certificate", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_api_credential.error_invalid_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclResourceDirectoryApiCredentialWithInvalidCertificate("uut", "directory-api-credential-invalid-certificate", "test_dir", rec.IsRecording()),
					ExpectError: regexp.MustCompile(`The certificate is not valid PEM format`),
				},
			},
		})
	})

	t.Run("error path - directory id is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceDirectoryApiCredentialWithMissingDirectoryId("uut", "directory-api-credential-no-directory-id"),
					ExpectError: regexp.MustCompile(`The argument "directory_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclResourceDirectoryApiCredential(resourceName string, apiCredentialName string, directoryName string, readOnly bool) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_api_credential" "%s"{
	name = "%s"
	directory_id = [for sa in data.btp_directories.all.values : sa.id if sa.name == "%s"][0]
	read_only = %t
}
	`, resourceName, apiCredentialName, directoryName, readOnly)
}

func hclResourceDirectoryApiCredentialWithCertificate(resourceName string, apiCredentialName string, directoryName string, recording bool) string {

	var directoryCertificate string
	if recording {
		directoryCertificate, _ = tfutils.ReadCertificate()
	} else {
		directoryCertificate = "redacted"
	}

	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_api_credential" "%s"{
	name = "%s"
	directory_id = [for sa in data.btp_directories.all.values : sa.id if sa.name == "%s"][0]
	certificate_passed = "%s"
}
	`, resourceName, apiCredentialName, directoryName, directoryCertificate)
}

func hclResourceDirectoryApiCredentialWithInvalidCertificate(resourceName string, apiCredentialName string, directoryName string, recording bool) string {

	var directoryCertificate string
	if recording {
		directoryCertificate = "Invalid-PEM-Certificate"
	} else {
		directoryCertificate = "redacted"
	}

	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_api_credential" "%s"{
	name = "%s"
	directory_id = [for sa in data.btp_directories.all.values : sa.id if sa.name == "%s"][0]
	certificate_passed = "%s"
}
	`, resourceName, apiCredentialName, directoryName, directoryCertificate)
}

func hclResourceDirectoryApiCredentialWithMissingDirectoryId(resourceName string, apiCredentialName string) string {
	return fmt.Sprintf(`
resource "btp_directory_api_credential" "%s"{
	name = "%s"
}
	`, resourceName, apiCredentialName)
}
