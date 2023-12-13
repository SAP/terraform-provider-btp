package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceBinding(t *testing.T) {

	t.Parallel()
	t.Run("happy path - service bindings by id", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_binding.by_id")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceBindingBySubaccountNameByBindingName("uut", "integration-test-services-static", "test-service-binding-malware-scanner"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_binding.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_binding.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_binding.uut", "name", "test-service-binding-malware-scanner"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_binding.uut", "ready", "true"),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_binding.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_binding.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
			},
		})

	})

	t.Run("happy path - service bindings by name", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_binding.by_name")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceBindingByNameBySubaccountNameByBindingName("uut", "integration-test-services-static", "test-service-binding-malware-scanner"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_binding.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_binding.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_binding.uut", "name", "test-service-binding-malware-scanner"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_binding.uut", "ready", "true"),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_binding.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_binding.uut", "last_modified", regexpValidRFC3999Format),
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
					Config:      hclDatasourceSubaccountServiceBindingNoSubaccount("uut", "any-sb-name"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - no ID or name", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountServiceBindingNoIdOrName("uut", "00000000-0000-0000-0000-000000000000"),
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Combination`),
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
					Config:      hclDatasourceSubaccountServiceBindingByNameBySubaccountIdByBindingName("uut", "this-is-not-a-uuid", "any-sb-name"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountServiceBindingBySubaccountNameByBindingName(resourceName string, subaccountName string, bindingName string) string {
	template := `
data "btp_subaccounts" "allsas" {}
data "btp_subaccount_service_bindings" "allsbs" {
  subaccount_id = [for sa in data.btp_subaccounts.allsas.values : sa.id if sa.name == "%[2]s"][0]
}
data "btp_subaccount_service_binding" "%[1]s" {
	subaccount_id = [for sa in data.btp_subaccounts.allsas.values : sa.id if sa.name == "%[2]s"][0]
	id            = [for sb in data.btp_subaccount_service_bindings.allsbs.values : sb.id if sb.name == "%[3]s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName, bindingName)
}

func hclDatasourceSubaccountServiceBindingByNameBySubaccountIdByBindingName(resourceName string, subaccountId string, bindingName string) string {
	template := `data "btp_subaccount_service_binding" "%s" {
	subaccount_id = "%s"
	name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, bindingName)
}

func hclDatasourceSubaccountServiceBindingByNameBySubaccountNameByBindingName(resourceName string, subaccountName string, bindingName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_binding" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, bindingName)
}

func hclDatasourceSubaccountServiceBindingNoSubaccount(resourceName string, bindingName string) string {
	template := `data "btp_subaccount_service_binding" "%s" {
	name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, bindingName)
}

func hclDatasourceSubaccountServiceBindingNoIdOrName(resourceName string, subaccountId string) string {
	template := `data "btp_subaccount_service_binding" "%s" {
	subaccount_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId)
}
