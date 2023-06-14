package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalAccount(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalAccount("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "id", "03760ecf-9d89-4189-a92a-1c7efed09298"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "name", "terraform-integration-canary"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "contract_status", "ACTIVE"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "geo_access", "STANDARD"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "license_type", "SAPDEV"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "costobject_type", "COST_CENTER"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "usage", "Testing"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "commercial_model", "Subscription"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "consumption_based", "true"),
					),
				},
			},
		})
	})
}

func hclDatasourceGlobalAccount(resourceName string) string {
	return fmt.Sprintf(`data "btp_globalaccount" "%s" {}`, resourceName)
}
