package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountRoleCollection(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_role_collection.role_collection_exists")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountRoleCollection("uut", "Global Account Administrator"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "description", "Administrative access to the global account"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "read_only", "true"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "roles.#", "4"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection not available", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_role_collection.role_collection_not_available")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountRoleCollection("uut", "fuh"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "description", ""),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "read_only", "false"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "roles.#", "0"),
					),
				},
			},
		})
	})

}

func hclDatasourceGlobalaccountRoleCollection(resourceName string, name string) string {
	template := `data "btp_globalaccount_role_collection" "%s" {
		name               = "%s"
	  }`

	return fmt.Sprintf(template, resourceName, name)
}
