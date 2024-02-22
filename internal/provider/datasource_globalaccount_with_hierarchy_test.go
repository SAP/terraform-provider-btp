package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountWithHierarchy(t *testing.T) {
	t.Parallel()
	t.Run("happy path -- empty globalaccount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_with_hierarchy.empty_globalaccount")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountWithHierarchy("globalaccount_canary"),
					//Config: hclProvider(user) + hclDatasourceGlobalaccountWithHierarchy("globalaccount_canary"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_with_hierarchy.globalaccount_canary", "subaccounts.#", "4"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_with_hierarchy.globalaccount_canary", "directories.#", "5"),
					),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountWithHierarchy(resourceName string) string {
	template := `data "btp_globalaccount_with_hierarchy" "%s" {}`
	return fmt.Sprintf(template, resourceName)
}
