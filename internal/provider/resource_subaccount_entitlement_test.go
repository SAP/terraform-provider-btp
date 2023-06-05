package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceSubaccountEntitlement(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount_entitlement")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceSubaccountEntitlement("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "hana-cloud", "hana"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_name", "hana"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_id", "hana-cloud-hana"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "service_name", "hana-cloud"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "state", "OK"),
					),
				},
				{
					ResourceName:      "btp_subaccount_entitlement.uut",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})
}

func hclResourceSubaccountEntitlement(resourceName string, subaccountId string, serviceName string, planName string) string {
	return fmt.Sprintf(`resource "btp_subaccount_entitlement" "%s" {
        subaccount_id      = "%s"
        service_name    = "%s"
        plan_name = "%s"
    }`, resourceName, subaccountId, serviceName, planName)
}

/*
func hclResourceSubaccountEntitlementWithAmount(resourceName string, subaccountId string, serviceName string, planName string, amount string) string {
	return fmt.Sprintf(`resource "btp_subaccount_entitlement" "%s" {
        subaccount_id      = "%s"
        service_name    = "%s"
        plan_name = "%s"
        amount = %s
    }`, resourceName, subaccountId, serviceName, planName, amount)
}
*/
