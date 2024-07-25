package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceSubaccountSubscription(t *testing.T) {
	t.Run("happy path - simple subscription", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_subscription")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountSubscriptionBySubaccount("uut", "integration-test-services-static", "auditlog-viewer", "free"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_subscription.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_subscription.uut", "subaccount_id", regexpValidUUID),
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
				{
					ResourceName:      "btp_subaccount_subscription.uut",
					ImportStateIdFunc: getSubscriptionImportStateId("btp_subaccount_subscription.uut", "auditlog-viewer", "free"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - simple subscription with timeouts", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_subscription_with_timeouts")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountSubscriptionBySubaccountWithTimeout("uut", "integration-test-services-static", "auditlog-viewer", "free"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_subscription.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_subscription.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "app_name", "auditlog-viewer"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "plan_name", "free"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "app_id", "auditlog-viewer!t49"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "state", "SUBSCRIBED"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "quota", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "customer_developed", "false"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "authentication_provider", "XSUAA"),
						resource.TestMatchResourceAttr("btp_subaccount_subscription.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_subscription.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "timeouts.create", "25m"),
						resource.TestCheckResourceAttr("btp_subaccount_subscription.uut", "timeouts.delete", "15m"),
					),
				},
				{
					ResourceName:      "btp_subaccount_subscription.uut",
					ImportStateIdFunc: getSubscriptionImportStateId("btp_subaccount_subscription.uut", "auditlog-viewer", "free"),
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
					Config:      hclResourceSubaccountSubscriptionNoSubaccountId("uut", "auditlog-viewer", "free"),
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
					Config:      hclResourceSubaccountSubscriptionNoAppName("uut", "auditlog-viewer", "free"),
					ExpectError: regexp.MustCompile(`The argument "app_name" is required, but no definition was found`),
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
					Config:      hclResourceSubaccountSubscriptionNoPlan("uut", "00000000-0000-0000-0000-000000000000", "auditlog-viewer"),
					ExpectError: regexp.MustCompile(`The argument "plan_name" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - import failure", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_subscription.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountSubscriptionBySubaccount("uut", "integration-test-services-static", "auditlog-viewer", "free"),
				},
				{
					ResourceName:      "btp_subaccount_subscription.uut",
					ImportStateIdFunc: getSubscriptionImportStateIdNoAppNameNoPlanName("btp_subaccount_subscription.uut"),
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Unexpected Import Identifier`),
				},
			},
		})
	})

}

func hclResourceSubaccountSubscriptionBySubaccount(resourceName string, subaccountName string, appName string, planName string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		resource "btp_subaccount_subscription" "%s"{
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
			app_name         = "%s"
			plan_name        = "%s"
		}`, resourceName, subaccountName, appName, planName)
}

func hclResourceSubaccountSubscriptionBySubaccountWithTimeout(resourceName string, subaccountName string, appName string, planName string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		resource "btp_subaccount_subscription" "%s"{
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
			app_name         = "%s"
			plan_name        = "%s"
			timeouts = {
				create = "25m"
				delete = "15m"
			  }
		}`, resourceName, subaccountName, appName, planName)
}

func hclResourceSubaccountSubscriptionNoSubaccountId(resourceName string, appName string, planName string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_subscription" "%s"{
		    app_name         = "%s"
			plan_name        = "%s"
		}`, resourceName, appName, planName)
}

func hclResourceSubaccountSubscriptionNoAppName(resourceName string, subaccountId string, planName string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_subscription" "%s"{
		    subaccount_id    = "%s"
			plan_name        = "%s"
		}`, resourceName, subaccountId, planName)
}

func hclResourceSubaccountSubscriptionNoPlan(resourceName string, subaccountId string, appName string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_subscription" "%s"{
		    subaccount_id    = "%s"
			app_name         = "%s"
		}`, resourceName, subaccountId, appName)
}

func getSubscriptionImportStateId(resourceName string, appName string, planName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s,%s", rs.Primary.Attributes["subaccount_id"], appName, planName), nil
	}
}

func getSubscriptionImportStateIdNoAppNameNoPlanName(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.Attributes["subaccount_id"], nil
	}
}
