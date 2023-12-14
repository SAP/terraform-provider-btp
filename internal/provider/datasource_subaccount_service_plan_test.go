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
					Config: hclProviderFor(user) + hclDatasourceSubaccountPlanByIdBySubaccountIdFromFilteredList("uut", "integration-test-services-static", "lite"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_plan.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_plan.uut", "id", regexpValidUUID),
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
					Config: hclProviderFor(user) + hclDatasourceSubaccountPlanByNameByOfferingBySubaccountIdFromFilteredList("uut", "integration-test-services-static", "lite", "destination"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_plan.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_plan.uut", "id", regexpValidUUID),
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
					Config:      hclDatasourceSubaccountPlanWoOffering("uut", "00000000-0000-0000-0000-000000000000", "standard"),
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

func hclDatasourceSubaccountPlanByIdBySubaccountIdFromFilteredList(resourceName string, subaccountName string, planName string) string {
	template := `
data "btp_subaccounts" "allsas" {}
data "btp_subaccount_service_plans" "allssps" {
  subaccount_id = [for sa in data.btp_subaccounts.allsas.values : sa.id if sa.name == "%[2]s"][0]
}
data "btp_subaccount_service_plan" "%[1]s" {
     subaccount_id = [for sa in data.btp_subaccounts.allsas.values : sa.id if sa.name == "%[2]s"][0]
	 id            = [for ssp in data.btp_subaccount_service_plans.allssps.values : ssp.id if ssp.name == "%[3]s"][0]
}`

	return fmt.Sprintf(template, resourceName, subaccountName, planName)
}

func hclDatasourceSubaccountPlanByNameByOfferingBySubaccountIdFromFilteredList(resourceName string, subaccountName string, planName string, offeringName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_plan" "%s" {
    subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
    name          = "%s"
	offering_name = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, planName, offeringName)
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
