package provider

import (
	"fmt"
	"regexp"
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
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plans.uut", "values.#", "30"),
					),
				},
			},
		})
	})

	t.Run("happy path - service plans for subaccount and environment", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_subaccount_service_plans_cloudfoundry")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceSubaccountPlansBySubaccountAndEnvironment("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "cloudfoundry"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plans.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plans.uut", "values.#", "0"),
					),
				},
			},
		})
	})

	t.Run("happy path - service plans for subaccount with fields filter", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_subaccount_service_plans_namefilter")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceSubaccountPlansBySubaccountAndFields("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plans.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plans.uut", "values.#", "2"),
					),
				},
			},
		})
	})

	t.Run("error path - subaccount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + `data "btp_subaccount_service_plans" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})
	t.Run("error path - subaccount_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + hclDatasourceSubaccountPlansBySubaccount("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
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

func hclDatasourceSubaccountPlansBySubaccountAndEnvironment(resourceName string, subaccountId string, environment string) string {
	template := `
data "btp_subaccount_service_plans" "%s" { 
     subaccount_id = "%s"
     environment   = "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId, environment)
}

func hclDatasourceSubaccountPlansBySubaccountAndFields(resourceName string, subaccountId string) string {
	template := `
data "btp_subaccount_service_plans" "%s" { 
     subaccount_id = "%s"
     fields_filter = "name eq 'standard'"
}`

	return fmt.Sprintf(template, resourceName, subaccountId)
}
