package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceInstancess(t *testing.T) {
	t.Parallel()
	t.Run("happy path - service instances for subaccount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_instances_all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceInstanceBySubaccount("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "values.#", "2"),
					),
				},
			},
		})
	})

	t.Run("happy path - service instances for subaccount with fields filter", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_instances_namefilter")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountInstancesBySubaccountAndFields("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "false"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "values.#", "0"),
					),
				},
			},
		})
	})

	t.Run("happy path - service instances for subaccount with fields filter (variant)", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_instances_namefilter_variant")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountInstancesBySubaccountAndFields("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "true"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "values.#", "2"),
					),
				},
			},
		})
	})

	t.Run("happy path - service instances for subaccount with labels filter", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_instances_labelsfilter")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountInstancesBySubaccountAndLabels("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "values.#", "1"),
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
					Config:      `data "btp_subaccount_service_instances" "uut" {}`,
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
					Config:      hclDatasourceSubaccountServiceInstanceBySubaccount("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})

}

func hclDatasourceSubaccountServiceInstanceBySubaccount(resourceName string, subaccountId string) string {
	template := `
data "btp_subaccount_service_instances" "%s" {
     subaccount_id = "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId)
}

func hclDatasourceSubaccountInstancesBySubaccountAndFields(resourceName string, subaccountId string, usable string) string {
	template := `
data "btp_subaccount_service_instances" "%s" {
     subaccount_id = "%s"
	 fields_filter = "usable eq '%s'"
}`

	return fmt.Sprintf(template, resourceName, subaccountId, usable)
}

func hclDatasourceSubaccountInstancesBySubaccountAndLabels(resourceName string, subaccountId string) string {
	template := `
data "btp_subaccount_service_instances" "%s" {
     subaccount_id = "%s"
	 labels_filter = "org eq 'testvalue'"
}`

	return fmt.Sprintf(template, resourceName, subaccountId)
}
