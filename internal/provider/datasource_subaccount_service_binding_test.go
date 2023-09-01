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
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceBindingbyId("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "b02e4b22-906b-40c5-9c5e-dbb6a9068444"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_binding.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_binding.uut", "id", "b02e4b22-906b-40c5-9c5e-dbb6a9068444"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_binding.uut", "name", "test-service-binding-iban"),
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
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceBindingbyName("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "test-service-binding-iban"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_binding.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_binding.uut", "id", "b02e4b22-906b-40c5-9c5e-dbb6a9068444"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_binding.uut", "name", "test-service-binding-iban"),
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
					Config:      hclDatasourceSubaccountServiceBindingNoSubaccount("uut", "test-service-binding-iban"),
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
					Config:      hclDatasourceSubaccountServiceBindingNoIdOrName("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
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
					Config:      hclDatasourceSubaccountServiceBindingbyName("uut", "this-is-not-a-uuid", "test-service-binding-iban"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountServiceBindingbyId(resourceName string, subaccountId string, bindingId string) string {
	template := `data "btp_subaccount_service_binding" "%s" {
	subaccount_id = "%s"
	id            = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, bindingId)
}

func hclDatasourceSubaccountServiceBindingbyName(resourceName string, subaccountId string, bindingName string) string {
	template := `data "btp_subaccount_service_binding" "%s" {
	subaccount_id = "%s"
	name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, bindingName)
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
