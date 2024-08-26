package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func TestResourceGlobalaccountApiCredential(t *testing.T) {
	t.Parallel()

	t.Run("happy path - api-credential with client secret", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_api_credential.with_secret")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountApiCredential("uut", "globalaccount-api-credential-with-secret", false),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_api_credential.uut", "name", "globalaccount-api-credential-with-secret"),
						resource.TestMatchResourceAttr("btp_globalaccount_api_credential.uut", "globalaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_globalaccount_api_credential.uut", "credential_type", "secret"),
						resource.TestCheckResourceAttr("btp_globalaccount_api_credential.uut", "read_only", "false"),
					),
				},
			},
		})
	})

	t.Run("happy path - api-credential with certificate", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_api_credential.with_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountApiCredentialWithCertificate("uut", "globalaccount-api-credential-with-certificate", rec.IsRecording()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_api_credential.uut", "name", "globalaccount-api-credential-with-certificate"),
						resource.TestMatchResourceAttr("btp_globalaccount_api_credential.uut", "globalaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_globalaccount_api_credential.uut", "credential_type", "client certificate"),
						resource.TestCheckResourceAttr("btp_globalaccount_api_credential.uut", "read_only", "false"),
					),
				},
			},
		})
	})

	t.Run("happy path - api-credential with read-only set to true", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_api_credential.read_only_credentials")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountApiCredential("uut", "globalaccount-api-credential-read-only", true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_api_credential.uut", "name", "globalaccount-api-credential-read-only"),
						resource.TestMatchResourceAttr("btp_globalaccount_api_credential.uut", "globalaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_globalaccount_api_credential.uut", "credential_type", "secret"),
						resource.TestCheckResourceAttr("btp_globalaccount_api_credential.uut", "read_only", "true"),
					),
				},
			},
		})
	})

	t.Run("error path - invalid certificate", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_api_credential.error_invalid_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountApiCredentialWithInvalidCertificate("uut", "globalaccount-api-credential-invalid-certificate", rec.IsRecording()),
					ExpectError: regexp.MustCompile(`The certificate is not valid PEM format`),
				},
			},
		})
	})
}

func hclResourceGlobalaccountApiCredential (resourceName string, apiCredentialName string, readOnly bool) string {
	return fmt.Sprintf(`
resource "btp_globalaccount_api_credential" "%s"{
	name = "%s"
	read_only = %t
}
	`, resourceName, apiCredentialName, readOnly)
}

func hclResourceGlobalaccountApiCredentialWithCertificate (resourceName string, apiCredentialName string, recording bool) string {

	var globalaccountCertificate string
	if recording {
		globalaccountCertificate, _ = tfutils.ReadCertificate()
	} else {
		globalaccountCertificate = "redacted"
	}

	return fmt.Sprintf(`
resource "btp_globalaccount_api_credential" "%s"{
	name = "%s"
	certificate_passed = "%s"
}
	`, resourceName, apiCredentialName, globalaccountCertificate)
}

func hclResourceGlobalaccountApiCredentialWithInvalidCertificate (resourceName string, apiCredentialName string, recording bool) string {
	
	var globalaccountCertificate string
	if recording {
		globalaccountCertificate = "Invalid-PEM-Certificate"
	} else {
		globalaccountCertificate = "redacted"
	}
	
	return fmt.Sprintf(`
resource "btp_globalaccount_api_credential" "%s"{
	name = "%s"
	certificate_passed = "%s"
}
	`,resourceName, apiCredentialName, globalaccountCertificate)
}