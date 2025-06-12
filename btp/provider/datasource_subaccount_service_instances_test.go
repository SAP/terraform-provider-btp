package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceInstances(t *testing.T) {
	t.Parallel()
	t.Run("happy path - service instances for subaccount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_instances.all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceInstanceBySubaccount("uut", "integration-test-services-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instances.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "values.#", "4"),
					),
				},
			},
		})
	})

	t.Run("happy path - service instances for subaccount with fields filter", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_instances.namefilter")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountInstancesBySubaccountAndFields("uut", "integration-test-services-static", "false"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instances.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "values.#", "0"),
					),
				},
			},
		})
	})

	t.Run("happy path - service instances for subaccount with fields filter (variant)", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_instances.namefilter_variant")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountInstancesBySubaccountAndFields("uut", "integration-test-services-static", "true"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instances.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instances.uut", "values.#", "4"),
					),
				},
			},
		})
	})

	t.Run("happy path - service instances for subaccount with labels filter", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_instances.labelsfilter")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountInstancesBySubaccountAndLabels("uut", "integration-test-services-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instances.uut", "subaccount_id", regexpValidUUID),
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
					Config:      hclDatasourceSubaccountServiceInstanceBySubaccountId("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})

}

func hclDatasourceSubaccountServiceInstanceBySubaccountId(resourceName string, subaccountId string) string {
	template := `
data "btp_subaccount_service_instances" "%s" {
     subaccount_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId)
}

func hclDatasourceSubaccountServiceInstanceBySubaccount(resourceName string, subaccountName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_instances" "%s" {
     subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName)
}

func hclDatasourceSubaccountInstancesBySubaccountAndFields(resourceName string, subaccountName string, usable string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_instances" "%s" {
     subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	 fields_filter = "usable eq '%s'"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, usable)
}

func hclDatasourceSubaccountInstancesBySubaccountAndLabels(resourceName string, subaccountName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_instances" "%s" {
     subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	 labels_filter = "org eq 'testvalue'"
}`
	return fmt.Sprintf(template, resourceName, subaccountName)
}
