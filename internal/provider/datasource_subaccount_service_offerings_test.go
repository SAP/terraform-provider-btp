package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceOfferings(t *testing.T) {
	t.Parallel()
	t.Run("happy path - service offerings for subaccount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_offerings_all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountOfferingsBySubaccount("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offerings.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offerings.uut", "values.#", "17"),
					),
				},
			},
		})
	})

	t.Run("happy path - service offerings for subaccount and environment", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_offerings_by_environment")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountOfferingsBySubaccountAndEnvironment("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offerings.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offerings.uut", "values.#", "0"),
					),
				},
			},
		})
	})

	t.Run("happy path - service plans for subaccount with fields filter", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_offerings_namefilter")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountOfferingsBySubaccountAndFields("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offerings.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_offerings.uut", "values.#", "1"),
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
					Config:      `data "btp_subaccount_service_offerings" "uut" {}`,
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
					Config:      hclDatasourceSubaccountOfferingsBySubaccount("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})

}

func hclDatasourceSubaccountOfferingsBySubaccount(resourceName string, subaccountId string) string {
	template := `
data "btp_subaccount_service_offerings" "%s" {
     subaccount_id = "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId)
}

func hclDatasourceSubaccountOfferingsBySubaccountAndEnvironment(resourceName string, subaccountId string) string {
	template := `
data "btp_subaccount_service_offerings" "%s" {
     subaccount_id = "%s"
	 environment   = "cloudfoundry"
}`

	return fmt.Sprintf(template, resourceName, subaccountId)
}

func hclDatasourceSubaccountOfferingsBySubaccountAndFields(resourceName string, subaccountId string) string {
	template := `
data "btp_subaccount_service_offerings" "%s" {
     subaccount_id = "%s"
     fields_filter = "name eq 'html5-apps-repo'"
}`

	return fmt.Sprintf(template, resourceName, subaccountId)
}
