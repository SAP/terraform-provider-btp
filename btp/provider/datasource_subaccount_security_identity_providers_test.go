package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountIdentityProviders(t *testing.T) {
	t.Parallel()
	t.Run("happy path - list all idps", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_identity_providers.list_all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountIdentityProvidersList("uut", "integration-test-security-settings"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_identity_providers.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_identity_providers.uut", "values.#", "12"),
					),
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
					Config:      hclDatasourceSubaccountIdentityProvidersNoParams("uut"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - subaccount not found", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_identity_providers.subaccount_not_found")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceSubaccountIdentityProvidersWithSubaccountId("uut", "00000000-0000-0000-0000-000000000000"),
					ExpectError: regexp.MustCompile(`API Error Reading Available IdPs`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountIdentityProvidersList(resourceName string, subaccountName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_identity_providers" "%s" {
    subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName)
}

func hclDatasourceSubaccountIdentityProvidersNoParams(resourceName string) string {
	template := `data "btp_subaccount_identity_providers" "%s" {}`
	return fmt.Sprintf(template, resourceName)
}

func hclDatasourceSubaccountIdentityProvidersWithSubaccountId(resourceName string, subaccountId string) string {
	return fmt.Sprintf(`
data "btp_subaccount_identity_providers" "%s" {
    subaccount_id = "%s"
}`, resourceName, subaccountId)
}
