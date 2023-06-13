package provider

import (
	"encoding/json"
	"fmt"
	"os"
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
	t.Run("happy path - simple CF creation", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount_environment_instance.landscape_label")
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

	t.Run("error path - wrong plan name", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount_environment_instance.wrong_plan")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + hclResourceSubaccountEnvironmentInstanceCF("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "cloudFoundry-from-terraform", "default", "cf-eu12", "cf-terraform-org"),
					ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 500; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})

	t.Run("error path - no parameters", func(t *testing.T) {

		rec := setupVCR(t, "fixtures/resource_subaccount_environment_instance.no_parameters")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + hclResourceSubaccountEnvironmentInstanceCFWoParameters("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "cloudFoundry-from-terraform", "default", "cf-eu12"),
					ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 500; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})

}

func hclResourceSubaccountEnvironmentInstanceCF(resourceName string, subaccountId string, name string, planName string, landscapeLabel string, orgName string) string {

	cfParameters := cfOrgParameters{
		InstanceName: orgName,
		Users: []cfUsers{
			{
				Id:    os.Getenv("BTP_USERNAME"),
				Email: os.Getenv("BTP_USERNAME"),
			},
		},
	}

	jsonCfParameters, _ := json.Marshal(cfParameters)

	return fmt.Sprintf(`resource "btp_subaccount_environment_instance" "%s"{
		subaccount_id      = "%s"
		name               = "%s"
		environment_type   = "cloudfoundry"
		plan_name          = "%s"
		service_name       = "cloudfoundry"
		landscape_label    = "%s"
		parameters         = %q
	}`, resourceName, subaccountId, name, planName, landscapeLabel, string(jsonCfParameters))
}

func hclResourceSubaccountEnvironmentInstanceCFWoParameters(resourceName string, subaccountId string, name string, planName string, landscapeLabel string) string {

	return fmt.Sprintf(`resource "btp_subaccount_environment_instance" "%s"{
		subaccount_id      = "%s"
		name               = "%s"
		environment_type   = "cloudfoundry"
		plan_name          = "%s"
		service_name       = "cloudfoundry"
		landscape_label    = "%s"
		parameters         = ""
	}`, resourceName, subaccountId, name, planName, landscapeLabel)
}
