package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceSubaccountServiceBinding(t *testing.T) {
	// Using the alert notification service as ID for the service instance
	t.Run("happy path - simple service_binding", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_binding")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBinding("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "df532d07-57a7-415e-a261-23a398ef068a", "tfint-test-alert-sb"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_binding.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{
					ResourceName:      "btp_subaccount_service_binding.uut",
					ImportStateIdFunc: getServiceBindingImportStateId("btp_subaccount_service_binding.uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})
	t.Run("error path - subacount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountServiceBindingNoSubaccountId("uut", "df532d07-57a7-415e-a261-23a398ef068a", "tfint-test-alert-sb"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - service instance id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountServiceBindingNoServiceInstanceId("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "tfint-test-alert-sb"),
					ExpectError: regexp.MustCompile(`The argument "service_instance_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - service plan ID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountServiceBindingNoName("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "df532d07-57a7-415e-a261-23a398ef068a"),
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - import failure", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_binding_import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBinding("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "df532d07-57a7-415e-a261-23a398ef068a", "tfint-test-alert-sb"),
				},
				{
					ResourceName:      "btp_subaccount_service_binding.uut",
					ImportStateId:     "59cd458e-e66e-4b60-b6d8-8f219379f9a5",
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Unexpected Import Identifier`),
				},
			},
		})
	})

}

func hclResourceSubaccountServiceBinding(resourceName string, subaccountId string, serviceInstanceId string, name string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_binding" "%s"{
		    subaccount_id       = "%s"
			service_instance_id = "%s"
			name                = "%s"
		}`, resourceName, subaccountId, serviceInstanceId, name)
}

func hclResourceSubaccountServiceBindingNoSubaccountId(resourceName string, serviceInstanceId string, name string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_binding" "%s"{
			service_instance_id = "%s"
			name                = "%s"
		}`, resourceName, serviceInstanceId, name)
}

func hclResourceSubaccountServiceBindingNoServiceInstanceId(resourceName string, subaccountId string, name string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_binding" "%s"{
		    subaccount_id       = "%s"
			name                = "%s"
		}`, resourceName, subaccountId, name)
}

func hclResourceSubaccountServiceBindingNoName(resourceName string, subaccountId string, serviceInstanceId string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_binding" "%s"{
		    subaccount_id       = "%s"
			service_instance_id = "%s"
		}`, resourceName, subaccountId, serviceInstanceId)
}

func getServiceBindingImportStateId(resourceName string, subaccountId string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s,%s", subaccountId, rs.Primary.ID), nil
	}
}
