package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type cfUsers struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type cfOrgParameters struct {
	InstanceName string    `json:"instance_name"`
	Users        []cfUsers `json:"users"`
}

func TestResourceSubaccountEnvironmentInstance(t *testing.T) {
	t.Parallel()
	t.Run("happy path - simple CF creation", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount_environment_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceSubaccountEnvironmentInstanceCF("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "cloudFoundry-from-terraform", "standard", "cf-eu12", "cf-terraform-org"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_environment_instance.uut", "type", "Provision"),
						resource.TestMatchResourceAttr("btp_subaccount_environment_instance.uut", "labels", regexp.MustCompile(`"API Endpoint":"https:\/\/api\.cf\.eu12\.hana\.ondemand\.com"`)),
					),
				},
			},
		})
	})

	// Error cases for CREATE lead to errors as no resource was created, but plugin test framework tries to delete the non existent resources
	// See also: https://github.com/hashicorp/terraform-plugin-testing/issues/85
}

func hclResourceSubaccountEnvironmentInstanceCF(resourceName string, subaccountId string, name string, planName string, landscapeLabel string, orgName string) string {

	cfParameters := cfOrgParameters{
		InstanceName: orgName,
		Users: []cfUsers{
			{
				Id:    "john.doe@int.test",
				Email: "john.doe@int.test",
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
}`, resourceName, subaccountId, name, planName, landscapeLabel, string(jsonCfParameters))
}
