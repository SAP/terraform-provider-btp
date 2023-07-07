package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceSubaccountSubscription(t *testing.T) {
	t.Run("happy path - simple subscription", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount_subscription")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceSubaccountSubscription("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "auditlog-viewer", "free"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_subscription.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "app_name", "auditlog-viewer"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "plan_name", "free"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "app_id", "auditlog-viewer!t49"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "state", "SUBSCRIBED"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "quota", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "customer_developed", "false"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "authentication_provider", "XSUAA"),
						resource.TestMatchResourceAttr("btp_subaccount_subscription.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_subscription.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				/*
					{
						ResourceName:      "btp_subaccount_subscription.uut",
						ImportStateId:     "ef23ace8-6ade-4d78-9c1f-8df729548bbf,auditlog-viewer,free",
						ImportState:       true,
						ImportStateVerify: true,
					},
				*/
			},
		})
	})
}

func hclResourceSubaccountSubscription(resourceName string, subaccountId string, appName string, planName string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_subscription" "%s"{
		    subaccount_id    = "%s"
			app_name         = "%s"
			plan_name        = "%s"
		}`, resourceName, subaccountId, appName, planName)
}
