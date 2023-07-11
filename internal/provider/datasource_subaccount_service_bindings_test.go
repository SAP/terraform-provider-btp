package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceBindings(t *testing.T) {

	t.Parallel()
	t.Run("happy path - all service bindings of a subaccount", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_subaccount_service_bindings")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceSubaccountSubscriptions("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_bindings.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_bindings.uut", "values.#", "3"),
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
					Config:      hclProvider() + `data "btp_subaccount_service_bindings" "uut" {}`,
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
					Config:      hclProvider() + hclDatasourceSubaccountSubscriptions("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountSubscriptions(resourceName string, subaccountId string) string {
	template := `data "btp_subaccount_service_bindings" "%s" {
	subaccount_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId)
}
