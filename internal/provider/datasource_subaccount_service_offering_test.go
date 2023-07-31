package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceoffering(t *testing.T) {
	t.Parallel()
	t.Run("happy path - service offering by id", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_offering_by_id")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountOfferingById("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "d67ff82d-9bfe-43e3-abd2-f2e21a5362c5"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "id", "d67ff82d-9bfe-43e3-abd2-f2e21a5362c5"),
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
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_offering_by_name")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountOfferingByName("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "xsuaa"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offering.uut", "id", "d67ff82d-9bfe-43e3-abd2-f2e21a5362c5"),
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
					Config:      hclDatasourceSubaccountServiceOfferingIdName("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "standard"),
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

func hclDatasourceSubaccountOfferingById(resourceName string, subaccountId string, offeringId string) string {
	template := `
data "btp_subaccount_service_offering" "%s" {
     subaccount_id = "%s"
	 id            = "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId, offeringId)
}

func hclDatasourceSubaccountOfferingByName(resourceName string, subaccountId string, offeringName string) string {
	template := `
data "btp_subaccount_service_offering" "%s" {
    subaccount_id = "%s"
    name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, offeringName)
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
