package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountEntitlement(t *testing.T) {
	t.Parallel()

	t.Run("happy path - plan_unique_identifier found", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_entitlement")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountEntitlement("uut", "integration-test-acc-static", "hana-cloud", "hana"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_entitlement.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_entitlement.uut", "service_name", "hana-cloud"),
						resource.TestCheckResourceAttr("data.btp_subaccount_entitlement.uut", "plan_name", "hana"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_entitlement.uut", "plan_unique_identifier"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_entitlement.uut", "plan_id"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_entitlement.uut", "quota_assigned"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_entitlement.uut", "quota_remaining"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_entitlement.uut", "category"),
					),
				},
			},
		})
	})

	t.Run("error path - subaccount_id is required", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_entitlement_required")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + `
						data "btp_subaccount_entitlement" "uut" {
						  service_name = "hana-cloud"
						  plan_name    = "hana"
						}
					`,
					ExpectError: regexp.MustCompile(`The argument \"subaccount_id\" is required`),
				},
			},
		})
	})

	t.Run("error path - plan not found", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_entitlement_invalid")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountEntitlement("uut", "integration-test-acc-static", "invalid-service", "invalid-plan"),
					ExpectError: regexp.MustCompile(
						`No plan found for service 'invalid-service' with plan_name = 'invalid-plan'`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountEntitlement(resourceName, subaccountId, serviceName, planName string) string {
	return fmt.Sprintf(`
	data "btp_subaccounts" "all" {}

	data "btp_subaccount_entitlement" "%s" {
	  subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	  service_name  = "%s"
	  plan_name     = "%s"
	}
	`, resourceName, subaccountId, serviceName, planName)
}
