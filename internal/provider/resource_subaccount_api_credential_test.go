package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

var (
	certificate, _ = tfutils.ReadCertificate()
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
					Config: hclProviderFor(user) + hclResourceSubaccountApiCredential("uut", "subaccount-api-credential", "integration-test-acc-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "name", "subaccount-api-credential"),
						resource.TestMatchResourceAttr("btp_subaccount_api_credential.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "credential_type", "secret"),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "read_only", "false"),
					),
				},
				{
					ResourceName: "btp_subaccount_api_credential.uut",
					ImportStateIdFunc: getImportStateIdForSubaccountApiCredential("btp_subaccount_api_credential.uut", "subaccount-api-credential"),
					ImportState: true,
				},
			},
		})
	})

	t.Run("happy path - api-credential with certificate", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_subaccount_api_credential.with_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountApiCredentialWithCertificate("uut", "subaccount-api-credential", "integration-test-acc-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "name", "subaccount-api-credential"),
						resource.TestMatchResourceAttr("btp_subaccount_api_credential.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "credential_type", "client certificate"),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "read_only", "false"),
					),
				},
				{
					ResourceName: "btp_subaccount_api_credential.uut",
					ImportStateIdFunc: getImportStateIdForSubaccountApiCredential("btp_subaccount_api_credential.uut", "subaccount-api-credential"),
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
					Config: hclProviderFor(user) + hclResourceSubaccountApiCredentialWithReadOnly("uut", "subaccount-api-credential", "integration-test-acc-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "name", "subaccount-api-credential"),
						resource.TestMatchResourceAttr("btp_subaccount_api_credential.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "credential_type", "secret"),
						resource.TestCheckResourceAttr("btp_subaccount_api_credential.uut", "read_only", "true"),
					),
				},
				{
					ResourceName: "btp_subaccount_api_credential.uut",
					ImportStateIdFunc: getImportStateIdForSubaccountApiCredential("btp_subaccount_api_credential.uut", "subaccount-api-credential"),
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
					Config: hclProviderFor(user) + hclResourceSubaccountApiCredentialWithInvalidCertificate("uut", "subaccount-api-credential", "integration-test-acc-static"),
					ExpectError: regexp.MustCompile(`The certificate is not valid PEM format`),
				},
			},
		})
	})

	t.Run("error path - missing subaccount id", func(t *testing.T){
		rec, user := setupVCR(t, "fixtures/resource_subaccount_api_credential.error_missing_subaccount_id")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest: true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountApiCredentialWithMissingSubaccountId("uut", "subaccount-api-credential"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclResourceSubaccountApiCredential (resourceName string, apiCredentialName string, subaccountName string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_api_credential" "%s"{
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
}
	`, resourceName, apiCredentialName, subaccountName)
}

func hclResourceSubaccountApiCredentialWithCertificate (resourceName string, apiCredentialName string, subaccountName string) string {

	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_api_credential" "%s"{
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	certificate_passed = "%s"
}
	`, resourceName, apiCredentialName, subaccountName, certificate)
}

func hclResourceSubaccountApiCredentialWithReadOnly (resourceName string, apiCredentialName string, subaccountName string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_api_credential" "%s"{
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	read_only = true
}
	`, resourceName, apiCredentialName, subaccountName)
}

func hclResourceSubaccountApiCredentialWithInvalidCertificate (resourceName string, apiCredentialName string, subaccountName string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_api_credential" "%s"{
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	certificate_passed = "Invalid-PEM-Certificate"
}
	`,resourceName, apiCredentialName, subaccountName)
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