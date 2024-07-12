package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func TestResourceDirectoryApiCredential(t *testing.T) {
	t.Parallel()

	t.Run("happy path - api-credential with client secret", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_directory_api_credential.with_secret")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryApiCredential("uut", "directory-api-credential-with-secret", "test-with_um", false),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "name", "directory-api-credential-with-secret"),
						resource.TestMatchResourceAttr("btp_directory_api_credential.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "credential_type", "secret"),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "read_only", "false"),
					),
				},
				{
					ResourceName: "btp_directory_api_credential.uut",
					ImportStateIdFunc: getImportStateIdForDirectoryApiCredential("btp_directory_api_credential.uut", "directory-api-credential-with-secret"),
					ImportState: true,
				},
			},
		})
	})

	t.Run("happy path - api-credential with certificate", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_directory_api_credential.with_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryApiCredentialWithCertificate("uut", "directory-api-credential-with-certificate", "test-with_um", rec.IsRecording()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "name", "directory-api-credential-with-certificate"),
						resource.TestMatchResourceAttr("btp_directory_api_credential.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "credential_type", "client certificate"),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "read_only", "false"),
					),
				},
				{
					ResourceName: "btp_directory_api_credential.uut",
					ImportStateIdFunc: getImportStateIdForDirectoryApiCredential("btp_directory_api_credential.uut", "directory-api-credential-with-certificate"),
					ImportState: true,
				},
			},
		})
	})

	//Please note that the following test case must not be re-recorded within the same 24-hour period,
	//as the deletion of credentials with certificates takes some time to reflect within the directory.
	t.Run("happy path - api-credential with read-only set to true", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_directory_api_credential.read_only_credentials")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryApiCredential("uut", "directory-api-credential-read-only", "test-with_um", true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "name", "directory-api-credential-read-only"),
						resource.TestMatchResourceAttr("btp_directory_api_credential.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "credential_type", "secret"),
						resource.TestCheckResourceAttr("btp_directory_api_credential.uut", "read_only", "true"),
					),
				},
				{
					ResourceName: "btp_directory_api_credential.uut",
					ImportStateIdFunc: getImportStateIdForDirectoryApiCredential("btp_directory_api_credential.uut", "directory-api-credential-read-only"),
					ImportState: true,
				},
			},
		})
	})

	t.Run("error path - invalid certificate", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_directory_api_credential.error_invalid_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryApiCredentialWithInvalidCertificate("uut", "directory-api-credential-invalid-certificate", "test-with_um", rec.IsRecording()),
					ExpectError: regexp.MustCompile(`The certificate is not valid PEM format`),
				},
			},
		})
	})

	t.Run("error path - directory id is mandatory", func(t *testing.T){
		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: hclResourceDirectoryApiCredentialWithMissingDirectoryId("uut", "directory-api-credential-no-directory-id"),
					ExpectError: regexp.MustCompile(`The argument "directory_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclResourceDirectoryApiCredential (resourceName string, apiCredentialName string, directoryName string, readOnly bool) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_api_credential" "%s"{
	name = "%s"
	directory_id = [for sa in data.btp_directories.all.values : sa.id if sa.name == "%s"][0]
	read_only = %t
}
	`, resourceName, apiCredentialName, directoryName, readOnly)
}

func hclResourceDirectoryApiCredentialWithCertificate (resourceName string, apiCredentialName string, directoryName string, recording bool) string {

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

func hclResourceDirectoryApiCredentialWithInvalidCertificate (resourceName string, apiCredentialName string, directoryName string, recording bool) string {
	
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
	`,resourceName, apiCredentialName, directoryName, directoryCertificate)
}

func hclResourceDirectoryApiCredentialWithMissingDirectoryId (resourceName string, apiCredentialName string) string {
	return fmt.Sprintf(`
resource "btp_directory_api_credential" "%s"{
	name = "%s"
}
	`,resourceName, apiCredentialName)
}

func getImportStateIdForDirectoryApiCredential(resourceName string, apiCredentialName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error){
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["directory_id"], apiCredentialName), nil
	}
}