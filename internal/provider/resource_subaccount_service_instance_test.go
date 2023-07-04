package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type testDestinationEntry struct {
	Name           string `json:"Name"`
	Type           string `json:"Type"`
	Url            string `json:"URL"`
	Authentication string `json:"Authentication"`
	ProxyType      string `json:"ProxyType"`
	Description    string `json:"Description"`
}

type testDestinationSubaccountData struct {
	ExistingDestinationPolicy string                 `json:"existing_destinations_policy"`
	Destination               []testDestinationEntry `json:"destinations"`
}

type testDestinationInitData struct {
	Subaccount testDestinationSubaccountData `json:"subaccount"`
}

type testParamsDestination struct {
	HTML5RuntimeEnabled string                  `json:"HTML5Runtime_enable"`
	InitData            testDestinationInitData `json:init_data`
}

func TestResourceSubaccountServiceInstance(t *testing.T) {
	t.Run("happy path - simple service creation wo parameters", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount_service_instance_wo_parameters")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceSubaccountServiceInstanceWoParameters("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "tf-test-alert-notification", "4bf8a2c4-6277-4bb1-b80d-2e46e87bd1a5"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", "4bf8a2c4-6277-4bb1-b80d-2e46e87bd1a5"),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "tf-test-alert-notification"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "platform_id", "service-manager"),
					),
				},
				{
					ResourceName:      "btp_subaccount_service_instance.uut",
					ImportStateIdFunc: getServiceInstanceIdForImport("btp_subaccount_service_instance.uut"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - simple service creation with parameters", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount_service_instance_with_parameters")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceSubaccountServiceInstanceWithParameters("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "tf-test-destintion", "cdf9c103-ef56-43e5-ac1d-4f1c5b15e05c"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", "cdf9c103-ef56-43e5-ac1d-4f1c5b15e05c"),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "tf-test-destintion"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "platform_id", "service-manager"),
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
					Config:      hclProvider() + hclResourceSubaccountServiceInstanceNoSubaccountId("uut", "tf-test-alert-notification", "4bf8a2c4-6277-4bb1-b80d-2e46e87bd1a5"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - service name mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + hclResourceSubaccountServiceInstanceNoServicName("uut", "tf-test-alert-notification", "4bf8a2c4-6277-4bb1-b80d-2e46e87bd1a5"),
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found`),
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
					Config:      hclProvider() + hclResourceSubaccountServiceInstanceNoPlan("uut", "this-is-not-a-uuid", "tf-test-alert-notification"),
					ExpectError: regexp.MustCompile(`The argument "serviceplan_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - import failure", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount_service_instance_import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceSubaccountServiceInstanceWoParameters("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "tf-test-alert-notification", "4bf8a2c4-6277-4bb1-b80d-2e46e87bd1a5"),
				},
				{
					ResourceName:      "btp_subaccount_service_instance.uut",
					ImportStateId:     "59cd458e-e66e-4b60-b6d8-8f219379f9a5",
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Unexpected Import Identifier`),
				},
			},
		})
	})
}

func hclResourceSubaccountServiceInstanceWoParameters(resourceName string, subaccountId string, name string, servicePlanId string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_instance" "%s"{
		    subaccount_id    = "%s"
			name             = "%s"
			serviceplan_id   = "%s"
		}`, resourceName, subaccountId, name, servicePlanId)
}

func hclResourceSubaccountServiceInstanceWithParameters(resourceName string, subaccountId string, name string, servicePlanId string) string {

	destinationInitData := testDestinationInitData{
		Subaccount: testDestinationSubaccountData{
			ExistingDestinationPolicy: "fail",
			Destination: []testDestinationEntry{
				{
					Name:           "Task_Center_global_settings",
					Type:           "HTTP",
					Url:            "http://sap.com",
					Authentication: "NoAuthentication",
					ProxyType:      "Internet",
					Description:    "SAP Task Center Global Settings",
				},
			},
		},
	}

	destParameters := testParamsDestination{
		HTML5RuntimeEnabled: "true",
		InitData:            destinationInitData,
	}

	destParametersJson, _ := json.Marshal(destParameters)

	return fmt.Sprintf(`
		resource "btp_subaccount_service_instance" "%s"{
		    subaccount_id    = "%s"
			name             = "%s"
			serviceplan_id   = "%s"
			parameters       = %q
		}`, resourceName, subaccountId, name, servicePlanId, string(destParametersJson))
}

func hclResourceSubaccountServiceInstanceNoSubaccountId(resourceName string, name string, servicePlanId string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_instance" "%s"{
		    name             = "%s"
			serviceplan_id   = "%s"
		}`, resourceName, name, servicePlanId)
}

func hclResourceSubaccountServiceInstanceNoServicName(resourceName string, subaccountId string, servicePlanId string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_instance" "%s"{
		    subaccount_id    = "%s"
			serviceplan_id   = "%s"
		}`, resourceName, subaccountId, servicePlanId)
}

func hclResourceSubaccountServiceInstanceNoPlan(resourceName string, subaccountId string, name string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_instance" "%s"{
		    subaccount_id    = "%s"
			name             = "%s"
		}`, resourceName, subaccountId, name)
}

func getServiceInstanceIdForImport(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s,%s", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", rs.Primary.ID), nil
	}
}
