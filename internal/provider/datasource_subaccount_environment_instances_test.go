package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountEnvironmentInstances(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_subaccount_environment_instances")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceSubaccountEnvironmentInstances("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_environment_instances.uut", "subaccount_id", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
						resource.TestCheckResourceAttr("data.btp_subaccount_environment_instances.uut", "values.#", "0"),
					),
				},
			},
		})
	})

}

func hclDatasourceSubaccountEnvironmentInstances(resourceName string, subaccountId string) string {
	template := `
data "btp_subaccount_environment_instances" "%s" {
  subaccount_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId)
}
