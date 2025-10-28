package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountApps(t *testing.T) {

	t.Parallel()
	t.Run("happy path - all apps of subaccount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_apps")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountApps("uut", "integration-test-services-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_apps.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_apps.uut", "values.#", "18"),
					),
				},
			},
		})

	})
	t.Run("error path - subaccount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_subaccount_apps" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountApps(resourceName string, subaccountName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_apps" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName)
}
