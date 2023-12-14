package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountSubscription(t *testing.T) {

	t.Parallel()

	t.Run("happy path - get subscriptions by id and plan", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_subscription.by_id_and_plan")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountSubscriptionByIdByPlanBySubaccountIdFromFilteredList("uut", "integration-test-services-static", "content-agent-ui", "free"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_subscription.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_subscription.uut", "app_name", "content-agent-ui"),
						resource.TestCheckResourceAttr("data.btp_subaccount_subscription.uut", "plan_name", "free"),
						resource.TestCheckResourceAttr("data.btp_subaccount_subscription.uut", "state", "SUBSCRIBED"),
						resource.TestCheckResourceAttr("data.btp_subaccount_subscription.uut", "quota", "1"),
						resource.TestMatchResourceAttr("data.btp_subaccount_subscription.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_subscription.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
			},
		})
	})

	t.Run("error path - subaccount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountSubscriptionNoSubaccountId("uut", "content-agent-ui", "free"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - app_name mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountSubscriptionNoAppName("uut", "00000000-0000-0000-0000-0000000000000", "free"),
					ExpectError: regexp.MustCompile(`The argument "app_name" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - subaccount_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountSubscriptionByIdAndPlan("uut", "this-is-not-a-uuid", "content-agent-ui", "free"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountSubscriptionByIdAndPlan(resourceName string, subaccountId string, appName string, planName string) string {
	template := `
data "btp_subaccount_subscription" "%s" {
	subaccount_id = "%s"
	app_name      = "%s"
	plan_name     = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, appName, planName)
}

func hclDatasourceSubaccountSubscriptionByIdByPlanBySubaccountIdFromFilteredList(resourceName string, subaccountName string, appName string, planName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_subscription" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	app_name      = "%s"
	plan_name     = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, appName, planName)
}

func hclDatasourceSubaccountSubscriptionNoSubaccountId(resourceName string, appName string, planName string) string {
	template := `data "btp_subaccount_subscription" "%s" {
	app_name      = "%s"
	plan_name     = "%s"
}`
	return fmt.Sprintf(template, resourceName, appName, planName)
}

func hclDatasourceSubaccountSubscriptionNoAppName(resourceName string, subaccountId string, planName string) string {
	template := `data "btp_subaccount_subscription" "%s" {
    subaccount_id = "%s"
	plan_name     = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, planName)
}
