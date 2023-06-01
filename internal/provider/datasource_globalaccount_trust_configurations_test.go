package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountTrustConfigurations(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_trust_configurations")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountTrustConfigurations("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configurations.uut", "values.#", "2"),
					),
				},
			},
		})
	})

}

func hclDatasourceGlobalaccountTrustConfigurations(resourceName string) string {
	template := `data "btp_globalaccount_trust_configurations" "%s" {}`
	return fmt.Sprintf(template, resourceName)
}
