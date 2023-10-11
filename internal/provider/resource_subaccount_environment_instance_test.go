package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type cfUsers struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type cfOrgParameters struct {
	InstanceName string    `json:"instance_name"`
	Users        []cfUsers `json:"users,omitempty"`
}

func TestResourceSubaccountEnvironmentInstance(t *testing.T) {
	t.Parallel()
	t.Run("happy path - simple CF creation", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_environment_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEnvironmentInstanceCF("uut",
						"ef23ace8-6ade-4d78-9c1f-8df729548bbf",
						"cloudFoundry-from-terraform",
						"standard",
						"cf-eu12",
						"cf-terraform-org",
						"john.doe@int.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "type", "Provision"),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "landscape_label", "cf-eu12"),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "labels", regexp.MustCompile(`"API Endpoint":"https:\/\/api\.cf\.eu12\.hana\.ondemand\.com"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"instance_name":"cf-terraform-org"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"id":"john.doe@int.test"`)),
					),
				},
				{
					ResourceName:      "btp_subaccount_environment_instance.uut",
					ImportStateIdFunc: getEnvironmentInstanceIdForImport("btp_subaccount_environment_instance.uut"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - CF creation with timeout", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_environment_instance.with_timeout")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEnvironmentInstanceCFTimeout("uut",
						"ef23ace8-6ade-4d78-9c1f-8df729548bbf",
						"cloudFoundry-from-terraform",
						"standard",
						"cf-eu12",
						"cf-terraform-org",
						"john.doe@int.test",
						"60m",
						"20m",
						"30m"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "type", "Provision"),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "landscape_label", "cf-eu12"),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "labels", regexp.MustCompile(`"API Endpoint":"https:\/\/api\.cf\.eu12\.hana\.ondemand\.com"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"instance_name":"cf-terraform-org"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"id":"john.doe@int.test"`)),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "timeouts.create", "60m"),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "timeouts.update", "20m"),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "timeouts.delete", "30m"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEnvironmentInstanceCFTimeout("uut",
						"ef23ace8-6ade-4d78-9c1f-8df729548bbf",
						"cloudFoundry-from-terraform",
						"standard",
						"cf-eu12",
						"cf-terraform-org-updated",
						"john.doe@int.test",
						"60m",
						"20m",
						"30m"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "type", "Update"),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "landscape_label", "cf-eu12"),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "labels", regexp.MustCompile(`"API Endpoint":"https:\/\/api\.cf\.eu12\.hana\.ondemand\.com"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"instance_name":"cf-terraform-org-updated"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"id":"john.doe@int.test"`)),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "timeouts.create", "60m"),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "timeouts.update", "20m"),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "timeouts.delete", "30m"),
					),
				},
			},
		})
	})
	t.Run("happy path - update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_environment_instance.update")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEnvironmentInstanceCF("uut",
						"ef23ace8-6ade-4d78-9c1f-8df729548bbf",
						"cloudFoundry-from-terraform",
						"standard",
						"cf-eu12",
						"cf-terraform-org",
						"john.doe@int.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "type", "Provision"),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "landscape_label", "cf-eu12"),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "labels", regexp.MustCompile(`"API Endpoint":"https:\/\/api\.cf\.eu12\.hana\.ondemand\.com"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"instance_name":"cf-terraform-org"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"id":"john.doe@int.test"`)),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEnvironmentInstanceCF("uut",
						"ef23ace8-6ade-4d78-9c1f-8df729548bbf",
						"cloudFoundry-from-terraform",
						"standard",
						"cf-eu12",
						"cf-terraform-org",
						"jane.doe@int.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "type", "Update"),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "landscape_label", "cf-eu12"),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "labels", regexp.MustCompile(`"API Endpoint":"https:\/\/api\.cf\.eu12\.hana\.ondemand\.com"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"instance_name":"cf-terraform-org"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", notContainsCheckFunc(`"id":"john.doe@int.test"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"id":"jane.doe@int.test"`)),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEnvironmentInstanceCFWithOrgParamsWoLandscapeLabel("uut",
						"ef23ace8-6ade-4d78-9c1f-8df729548bbf",
						"cloudFoundry-from-terraform",
						"standard",
						cfOrgParameters{
							InstanceName: "cf-terraform-org",
						}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "type", "Update"),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "landscape_label", "cf-eu12"),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "labels", regexp.MustCompile(`"API Endpoint":"https:\/\/api\.cf\.eu12\.hana\.ondemand\.com"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", containsCheckFunc(`"instance_name":"cf-terraform-org"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", notContainsCheckFunc(`"id":"john.doe@int.test"`)),
						resource.TestCheckResourceAttrWith("btp_subaccount_environment_instance.uut", "parameters", notContainsCheckFunc(`"id":"jane.doe@int.test"`)),
					),
				},
			},
		})
	})

	// Error cases for CREATE lead to errors as no resource was created, but plugin test framework tries to delete the non existent resources
	// See also: https://github.com/hashicorp/terraform-plugin-testing/issues/85
}

func hclResourceSubaccountEnvironmentInstanceCF(resourceName string, subaccountId string, name string, planName string, landscapeLabel string, orgName string, user string) string {
	cfParameters := cfOrgParameters{
		InstanceName: orgName,
		Users: []cfUsers{
			{
				Id:    user,
				Email: user,
			},
		},
	}

	return hclResourceSubaccountEnvironmentInstanceCFWithOrgParams(resourceName, subaccountId, name, planName, landscapeLabel, cfParameters)
}

func hclResourceSubaccountEnvironmentInstanceCFWithOrgParams(resourceName string, subaccountId string, name string, planName string, landscapeLabel string, cfParameters cfOrgParameters) string {
	jsonCfParameters, _ := json.Marshal(cfParameters)

	return fmt.Sprintf(`
resource "btp_subaccount_environment_instance" "%s"{
    subaccount_id    = "%s"
	name             = "%s"
	environment_type = "cloudfoundry"
	plan_name        = "%s"
	service_name     = "cloudfoundry"
	landscape_label  = "%s"
	parameters       = %q
}`, resourceName, subaccountId, name, planName, landscapeLabel, string(jsonCfParameters))
}

func hclResourceSubaccountEnvironmentInstanceCFWithOrgParamsWoLandscapeLabel(resourceName string, subaccountId string, name string, planName string, cfParameters cfOrgParameters) string {
	jsonCfParameters, _ := json.Marshal(cfParameters)

	return fmt.Sprintf(`
resource "btp_subaccount_environment_instance" "%s"{
    subaccount_id    = "%s"
	name             = "%s"
	environment_type = "cloudfoundry"
	plan_name        = "%s"
	service_name     = "cloudfoundry"
	parameters       = %q
}`, resourceName, subaccountId, name, planName, string(jsonCfParameters))
}

func hclResourceSubaccountEnvironmentInstanceCFTimeout(resourceName string, subaccountId string, name string, planName string, landscapeLabel string, orgName string, user string, createTimeout string, updateTimeout string, deleteTimeout string) string {
	cfParameters := cfOrgParameters{
		InstanceName: orgName,
		Users: []cfUsers{
			{
				Id:    user,
				Email: user,
			},
		},
	}

	jsonCfParameters, _ := json.Marshal(cfParameters)

	return fmt.Sprintf(`
resource "btp_subaccount_environment_instance" "%s"{
    subaccount_id    = "%s"
	name             = "%s"
	environment_type = "cloudfoundry"
	plan_name        = "%s"
	service_name     = "cloudfoundry"
	landscape_label  = "%s"
	parameters       = %q
	timeouts = {
		create = "%s"
		update = "%s"
		delete = "%s"
	  }
}`, resourceName, subaccountId, name, planName, landscapeLabel, string(jsonCfParameters), createTimeout, updateTimeout, deleteTimeout)
}

func getEnvironmentInstanceIdForImport(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s,%s", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", rs.Primary.ID), nil
	}
}
