package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountUser(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_user")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountUser("uut", "jenny.doe@test.com", "sap.default"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "origin", "sap.default"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "active", "true"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "family_name", "unknown"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "given_name", "unknown"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "id", "86535387-54aa-4282-af13-67dd50cdd13c"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "role_collections.#", "2"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "verified", "false"),
					),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountUser(resourceName string, userName string, origin string) string {
	template := `data "btp_globalaccount_user" "%s" {
  user_name = "%s"
  origin    = "%s"
}`

	return fmt.Sprintf(template, resourceName, userName, origin)
}
