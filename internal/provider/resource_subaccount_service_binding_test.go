package provider

import (
	"fmt"
	"regexp"
	//"regexp"
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
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBindingByServiceInstanceBySubaccount("uut", "integration-test-services-static", "tf-testacc-alertnotification-instance", "tfint-test-alert-sb"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{
					ResourceName:      "btp_subaccount_service_binding.uut",
					ImportStateIdFunc: getServiceBindingImportStateId("btp_subaccount_service_binding.uut"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - simple service_binding with labels", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_binding.with_labels")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBindingWithLabelsByServiceInstanceBySubaccount("uut", "integration-test-services-static", "tf-testacc-alertnotification-instance", "tfint-test-alert-sb"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_binding.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_binding.uut", "labels.foo.0", "bar"),
					),
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
					Config:      hclResourceSubaccountServiceBindingNoSubaccountId("uut", "00000000-0000-0000-0000-000000000000", "tfint-test-alert-sb"),
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
					Config:      hclResourceSubaccountServiceBindingNoServiceInstanceId("uut", "00000000-0000-0000-0000-000000000000", "tfint-test-alert-sb"),
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
					Config:      hclResourceSubaccountServiceBindingNoName("uut", "00000000-0000-0000-0000-00000000000", "00000000-0000-0000-0000-00000000000"),
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - import failure", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_binding.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBindingByServiceInstanceBySubaccount("uut", "integration-test-services-static", "tf-testacc-alertnotification-instance", "tfint-test-alert-sb"),
				},
				{
					ResourceName:      "btp_subaccount_service_binding.uut",
					ImportStateIdFunc: getServiceBindingImportStateIdNoServiceBindingId("btp_subaccount_service_binding.uut"),
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Unexpected Import Identifier`),
				},
			},
		})
	})

}

func hclResourceSubaccountServiceBindingByServiceInstanceBySubaccount(resourceName string, subaccountId string, serviceInstanceId string, name string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		data "btp_subaccount_service_instances" "all" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
		}
		resource "btp_subaccount_service_binding" "%[1]s"{
		    subaccount_id       = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			service_instance_id = [for ssi in data.btp_subaccount_service_instances.all.values : ssi.id if ssi.name == "%[3]s"][0]
			name                = "%[4]s"
		}`, resourceName, subaccountId, serviceInstanceId, name)
}

func hclResourceSubaccountServiceBindingWithLabelsByServiceInstanceBySubaccount(resourceName string, subaccountName string, serviceInstanceName string, name string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		data "btp_subaccount_service_instances" "all" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
		}
		resource "btp_subaccount_service_binding" "%[1]s"{
		    subaccount_id       = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			service_instance_id = [for ssi in data.btp_subaccount_service_instances.all.values : ssi.id if ssi.name == "%[3]s"][0]
			name                = "%[4]s"
			labels              = {"foo" = ["bar"]}
		}`, resourceName, subaccountName, serviceInstanceName, name)
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

func getServiceBindingImportStateId(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["subaccount_id"], rs.Primary.ID), nil
	}
}

func getServiceBindingImportStateIdNoServiceBindingId(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s", rs.Primary.Attributes["subaccount_id"]), nil
	}
}
