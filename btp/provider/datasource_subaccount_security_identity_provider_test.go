package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountIdentityProvider(t *testing.T) {
	t.Parallel()
	t.Run("happy path - get idp by host", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_identity_provider.by_host")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{

					Config: hclProviderFor(user) + hclDatasourceSubaccountIdentityProviderByHost("uut", "integration-test-security-settings", "a2dynbhnd.accounts400.ondemand.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_identity_provider.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_identity_provider.uut", "host", "a2dynbhnd.accounts400.ondemand.com"),
						resource.TestCheckResourceAttr("data.btp_subaccount_identity_provider.uut", "tenant_id", "a2dynbhnd"),
						resource.TestCheckResourceAttr("data.btp_subaccount_identity_provider.uut", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("data.btp_subaccount_identity_provider.uut", "tenant_type", "trial"),
						resource.TestCheckResourceAttr("data.btp_subaccount_identity_provider.uut", "common_host", "a2dynbhnd.accounts400.cloud.sap"),
						resource.TestCheckResourceAttr("data.btp_subaccount_identity_provider.uut", "data_center_id", "EU2_QA"),
						resource.TestCheckResourceAttr("data.btp_subaccount_identity_provider.uut", "cost_center_id", "101000297"),
					),
				},
			},
		})
	})

	t.Run("error path - subaccount not found", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_identity_provider.subaccount_not_found")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceSubaccountIdentityProviderWithSubaccountId("uut", "00000000-0000-0000-0000-000000000000", "a2dynbhnd.accounts400.ondemand.com"),
					ExpectError: regexp.MustCompile(`API Error Reading IdP Details`),
				},
			},
		})
	})

	t.Run("error path - mandatory fields missing", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountIdentityProviderNoParams("uut"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountIdentityProviderByHost(resourceName string, subaccountName string, host string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_identity_provider" "%s" {
    subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
    host          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, host)
}

func hclDatasourceSubaccountIdentityProviderNoParams(resourceName string) string {
	template := `data "btp_subaccount_identity_provider" "%s" {}`
	return fmt.Sprintf(template, resourceName)
}

func hclDatasourceSubaccountIdentityProviderWithSubaccountId(resourceName string, subaccountId string, host string) string {
	return fmt.Sprintf(`
data "btp_subaccount_identity_provider" "%s" {
    subaccount_id = "%s"
    host          = "%s"
}`, resourceName, subaccountId, host)
}
