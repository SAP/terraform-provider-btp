package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceDirectoryEntitlement(t *testing.T) {
	t.Parallel()

	t.Run("happy path - data source fetches correct values", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/data_source_directory_entitlement_success")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDirectoryEntitlementDataSource("uut", "integration-test-dir-se-static", "alert-notification", "lite"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_directory_entitlement.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_directory_entitlement.uut", "service_name", "alert-notification"),
						resource.TestCheckResourceAttr("data.btp_directory_entitlement.uut", "plan_name", "lite"),
						resource.TestCheckResourceAttr("data.btp_directory_entitlement.uut", "plan_id", "alert-notification-lite"),
						resource.TestCheckResourceAttr("data.btp_directory_entitlement.uut", "plan_unique_identifier", "alert-notification-lite"),
					),
				},
			},
		})
	})
}

func hclDirectoryEntitlementDataSource(resourceName string, directoryName string, serviceName string, planName string) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}

data "btp_directory_entitlement" "%s" {
  directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
  service_name = "%s"
  plan_name    = "%s"
}
`, resourceName, directoryName, serviceName, planName)
}
