package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountWithHierarchy(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_with_hierarchy")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountWithHierarchy("globalaccount_canary"),
					Check: resource.ComposeAggregateTestCheckFunc(
						//directory: integration-test-dir-entitlements 's parent -->  terraform-integration-canary
						resource.TestCheckResourceAttr("data.btp_globalaccount_with_hierarchy.globalaccount_canary", "directories.2.parent_name", "terraform-integration-canary"),
						//directory: integration-test-dir-entitlements-stacked 's parent --> integration-test-dir-entitlements
						resource.TestCheckResourceAttr("data.btp_globalaccount_with_hierarchy.globalaccount_canary", "directories.2.directories.0.parent_name", "integration-test-dir-entitlements"),
						//subaccount: integration-test-acc-entitlements-stacked 's parent --> integration-test-dir-entitlements-stacked
						resource.TestCheckResourceAttr("data.btp_globalaccount_with_hierarchy.globalaccount_canary", "directories.2.directories.0.subaccounts.0.parent_name", "integration-test-dir-entitlements-stacked"),
						//subaccount: integration-test-acc-static 's parent --> terraform-integration-canary
						resource.TestCheckResourceAttr("data.btp_globalaccount_with_hierarchy.globalaccount_canary", "subaccounts.0.parent_name", "terraform-integration-canary"),
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
