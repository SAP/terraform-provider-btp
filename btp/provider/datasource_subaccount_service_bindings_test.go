package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceBindings(t *testing.T) {

	t.Parallel()
	t.Run("happy path - all service bindings of a subaccount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_bindings")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceBindings("uut", "integration-test-services-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_bindings.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_bindings.uut", "values.#", "4"),
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
					Config:      `data "btp_subaccount_service_bindings" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

}

func hclDatasourceSubaccountServiceBindings(resourceName string, subaccountName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_bindings" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName)
}
