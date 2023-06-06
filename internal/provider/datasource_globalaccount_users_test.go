package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountUsers(t *testing.T) {
	t.Parallel()
	t.Run("happy path - default idp", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_users.default_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountUsers("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_users.uut", "values.#", "11"),
					),
				},
			},
		})
	})
	t.Run("happy path with custom idp", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_users.custom_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountUsersWithCustomIdp("uut", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_users.uut", "values.#", "2"),
					),
				},
			},
		})
	})
	// TODO: error path with non existing idp
}

func hclDatasourceGlobalaccountUsers(resourceName string) string {
	template := `data "btp_globalaccount_users" "%s" {}`
	return fmt.Sprintf(template, resourceName)
}

func hclDatasourceGlobalaccountUsersWithCustomIdp(resourceName string, origin string) string {
	template := `data "btp_globalaccount_users" "%s" {
  origin    = "%s"
}`
	return fmt.Sprintf(template, resourceName, origin)
}
