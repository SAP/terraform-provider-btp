package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServicePlan(t *testing.T) {
	t.Parallel()
	t.Run("happy path - service plan by id", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_plan.by_id")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountPlanById("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "cdf9c103-ef56-43e5-ac1d-4f1c5b15e05c"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "id", "cdf9c103-ef56-43e5-ac1d-4f1c5b15e05c"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "name", "lite"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "ready", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "catalog_name", "lite"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "free", "true"),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_plan.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_plan.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
			},
		})
	})
	t.Run("happy path service plan  by name", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_plan.by_name")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountPlanByNameAndOffering("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "lite", "destination"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "id", "cdf9c103-ef56-43e5-ac1d-4f1c5b15e05c"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "name", "lite"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "ready", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "catalog_name", "lite"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_plan.uut", "free", "true"),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_plan.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_plan.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
			},
		})
	})

	t.Run("error path - offering name mandatory in case of name", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountPlanWoOffering("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "standard"),
					ExpectError: regexp.MustCompile(`Attribute "offering_name" must be specified when "name" is specified`),
				},
			},
		})
	})

	t.Run("error path - subaccount id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountPlanWoSubaccount("uut", "lite", "destination"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})

}

func hclDatasourceSubaccountPlanById(resourceName string, subaccountId string, planId string) string {
	template := `
data "btp_subaccount_service_plan" "%s" {
     subaccount_id = "%s"
	 id            = "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId, planId)
}

func hclDatasourceSubaccountPlanByNameAndOffering(resourceName string, subaccountId string, planName string, offeringName string) string {
	template := `
data "btp_subaccount_service_plan" "%s" {
    subaccount_id = "%s"
    name          = "%s"
	offering_name = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, planName, offeringName)
}

func hclDatasourceSubaccountPlanWoOffering(resourceName string, subaccountId string, planName string) string {
	template := `
data "btp_subaccount_service_plan" "%s" {
    subaccount_id = "%s"
    name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, planName)
}

func hclDatasourceSubaccountPlanWoSubaccount(resourceName string, planName string, offeringName string) string {
	template := `
data "btp_subaccount_service_plan" "%s" {
    name          = "%s"
	offering_name = "%s"
}`
	return fmt.Sprintf(template, resourceName, planName, offeringName)
}
