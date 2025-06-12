package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceOffering(t *testing.T) {
	t.Parallel()

	t.Run("happy path - service offering by id", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_offering.by_id")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceOfferingByIdBySubaccountIdFilteredFromList("uut", "integration-test-services-static", "xsuaa"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_offering.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_offering.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "name", "xsuaa"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "ready", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "bindable", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "instances_retrievable", "false"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "bindings_retrievable", "false"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "plan_updateable", "false"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "allow_context_updates", "false"),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_offering.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_offering.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
			},
		})
	})

	t.Run("happy path service offering by name", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_offering.by_name")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceOfferingByNameBySubaccountIdFilteredFromList("uut", "integration-test-services-static", "xsuaa"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_offering.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_offering.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "name", "xsuaa"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "ready", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "bindable", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "instances_retrievable", "false"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "bindings_retrievable", "false"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "plan_updateable", "false"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "allow_context_updates", "false"),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_offering.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_offering.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
			},
		})
	})

	t.Run("error path - specify ID and name", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountServiceOfferingIdName("uut", "00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", "standard"),
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Combination`),
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
					Config:      hclDatasourceSubaccountOfferingWoSubaccount("uut", "lite"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})

}

func hclDatasourceSubaccountServiceOfferingByIdBySubaccountIdFilteredFromList(resourceName string, subaccountName string, serviceOfferingName string) string {
	template := `
data "btp_subaccounts" "allsas" {}
data "btp_subaccount_service_offerings" "allssos" {
  subaccount_id = [for sa in data.btp_subaccounts.allsas.values : sa.id if sa.name == "%[2]s"][0]
}
data "btp_subaccount_service_offering" "%[1]s" {
     subaccount_id = [for sa in data.btp_subaccounts.allsas.values : sa.id if sa.name == "%[2]s"][0]
	 id            = [for sso in data.btp_subaccount_service_offerings.allssos.values : sso.id if sso.name == "%[3]s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName, serviceOfferingName)
}

func hclDatasourceSubaccountServiceOfferingByNameBySubaccountIdFilteredFromList(resourceName string, subaccountName string, offeringName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_offering" "%s" {
    subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
    name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, offeringName)
}

func hclDatasourceSubaccountServiceOfferingIdName(resourceName string, subaccountId string, offeringId string, offeringName string) string {
	template := `
data "btp_subaccount_service_offering" "%s" {
    subaccount_id = "%s"
	id            = "%s"
    name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, offeringId, offeringName)
}

func hclDatasourceSubaccountOfferingWoSubaccount(resourceName string, offeringName string) string {
	template := `
data "btp_subaccount_service_offering" "%s" {
	offering_name = "%s"
}`
	return fmt.Sprintf(template, resourceName, offeringName)
}
