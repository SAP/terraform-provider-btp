package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServicePlans(t *testing.T) {
	t.Parallel()
	t.Run("happy path - service plans for subaccount", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_subaccount_service_plans_all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceSubaccountPlansBySubaccount("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plans.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					),
				},
			},
		})
	})

}

func hclDatasourceSubaccountPlansBySubaccount(resourceName string, subaccountId string) string {
	template := `
data "btp_subaccount_service_plans" "%s" { 
     subaccount_id = "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId)
}
