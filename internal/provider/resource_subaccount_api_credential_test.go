package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func TestResourceSubaccountApiCredential(t *testing.T) {
	t.Parallel()

	t.Run("happy path - api-credential with client secret", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_subaccount_api_credential.with_secret")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountApiCredential("uut", "subaccount-api-credential-with-secret", "integration-test-acc-static", false),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "name", "subaccount-api-credential-with-secret"),
						resource.TestMatchResourceAttr("btp_subaccount_api_credential.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "credential_type", "secret"),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "read_only", "false"),
					),
				},
				{
					ResourceName: "btp_subaccount_api_credential.uut",
					ImportStateIdFunc: getImportStateIdForSubaccountApiCredential("btp_subaccount_api_credential.uut", "subaccount-api-credential-with-secret"),
					ImportState: true,
				},
			},
		})
	})

	//Please note that the following test case must not be re-recorded within the same 24-hour period,
	//as the deletion of credentials with certificates takes some time to reflect within the subaccount.
	t.Run("happy path - api-credential with certificate", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_subaccount_api_credential.with_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountApiCredentialWithCertificate("uut", "subaccount-api-credential-with-certificate", "integration-test-acc-static", rec.IsRecording()),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "name", "subaccount-api-credential-with-certificate"),
						resource.TestMatchResourceAttr("btp_subaccount_api_credential.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "credential_type", "client certificate"),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "read_only", "false"),
					),
				},
				{
					ResourceName: "btp_subaccount_api_credential.uut",
					ImportStateIdFunc: getImportStateIdForSubaccountApiCredential("btp_subaccount_api_credential.uut", "subaccount-api-credential-with-certificate"),
					ImportState: true,
				},
			},
		})
	})

	t.Run("happy path - api-credential with read-only set to true", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_subaccount_api_credential.read_only_credentials")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountApiCredential("uut", "subaccount-api-credential-read-only", "integration-test-acc-static", true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "name", "subaccount-api-credential-read-only"),
						resource.TestMatchResourceAttr("btp_subaccount_api_credential.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "credential_type", "secret"),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "read_only", "true"),
					),
				},
				{
					ResourceName: "btp_subaccount_api_credential.uut",
					ImportStateIdFunc: getImportStateIdForSubaccountApiCredential("btp_subaccount_api_credential.uut", "subaccount-api-credential-read-only"),
					ImportState: true,
				},
			},
		})
	})

	t.Run("error path - invalid certificate", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_subaccount_api_credential.error_invalid_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountApiCredentialWithInvalidCertificate("uut", "subaccount-api-credential-invalid-certificate", "integration-test-acc-static", rec.IsRecording()),
					ExpectError: regexp.MustCompile(`The certificate is not valid PEM format`),
				},
			},
		})
	})

	t.Run("error path - subaccount id is mandatory", func(t *testing.T){
		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: hclResourceSubaccountApiCredentialWithMissingSubaccountId("uut", "subaccount-api-credential-no-subaccount-id"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclResourceSubaccountApiCredential (resourceName string, apiCredentialName string, subaccountName string, readOnly bool) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_api_credential" "%s"{
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	read_only = %t
}
	`, resourceName, apiCredentialName, subaccountName, readOnly)
}

func hclResourceSubaccountApiCredentialWithCertificate (resourceName string, apiCredentialName string, subaccountName string, recording bool) string {

	var subaccountCertificate string

	if recording {
		subaccountCertificate, _ = tfutils.ReadCertificate()
	} else {
		subaccountCertificate = "redacted"
	}

	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_api_credential" "%s"{
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	certificate_passed = "%s"
}
	`, resourceName, apiCredentialName, subaccountName, subaccountCertificate)
}

func hclResourceSubaccountApiCredentialWithInvalidCertificate (resourceName string, apiCredentialName string, subaccountName string, recording bool) string {
	
	var subaccountCertificate string
	if recording {
		subaccountCertificate = "Invalid-PEM-Certificate"
	} else {
		subaccountCertificate = "redacted"
	}

	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_api_credential" "%s"{
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	certificate_passed = "%s"
}
	`,resourceName, apiCredentialName, subaccountName, subaccountCertificate)
}

func hclResourceSubaccountApiCredentialWithMissingSubaccountId (resourceName string, apiCredentialName string) string {
	return fmt.Sprintf(`
resource "btp_subaccount_api_credential" "%s"{
	name = "%s"
}
	`,resourceName, apiCredentialName)
}

func getImportStateIdForSubaccountApiCredential(resourceName string, apiCredentialName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error){
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["subaccount_id"], apiCredentialName), nil
	}
}