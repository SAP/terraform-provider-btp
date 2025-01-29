package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountResourceProviders(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_resource_providers")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountResourceProviders("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_providers.uut", "values.#", "1"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_providers.uut", "values.0.provider_type", "AWS"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_providers.uut", "values.0.technical_name", "tf_test_resource_provider"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_providers.uut", "values.0.display_name", "Test AWS Resource Provider"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_providers.uut", "values.0.description", "Description of the resource provider"),
					),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountResourceProviders(resourceName string) string {
	return fmt.Sprintf(`data "btp_globalaccount_resource_providers" "%s" {}`, resourceName)
}
