package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceSubaccountEntitlement(t *testing.T) {
	t.Parallel()

	t.Run("happy path - no amount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_entitlement.no_amount")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEntitlementBySubaccount("uut", "integration-test-acc-static", "hana-cloud", "hana"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "id", "hana-cloud-hana"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_name", "hana"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_id", "hana-cloud-hana"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "service_name", "hana-cloud"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "amount", "3"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "state", "OK"),
					),
				},
				{
					ResourceName:      "btp_subaccount_entitlement.uut",
					ImportStateIdFunc: getImportStateIdForSubaccountEntitlement("btp_subaccount_entitlement.uut", "hana-cloud", "hana"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - directory hierarchy", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_entitlement.dir_hierarchy")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEntitlementBySubaccount("uut", "integration-test-acc-entitlements-stacked", "hana-cloud", "hana"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "id", "hana-cloud-hana"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_name", "hana"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_id", "hana-cloud-hana"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "service_name", "hana-cloud"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "amount", "3"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "state", "OK"),
					),
				},
				{
					ResourceName:      "btp_subaccount_entitlement.uut",
					ImportStateIdFunc: getImportStateIdForSubaccountEntitlement("btp_subaccount_entitlement.uut", "hana-cloud", "hana"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - with amount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_entitlement.amount_set")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEntitlementWithAmountBySubaccount("uut", "integration-test-acc-static", "data-privacy-integration-service", "standard", "3"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_name", "standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "service_name", "data-privacy-integration-service"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "amount", "3"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "state", "OK"),
					),
				},
				{
					ResourceName:      "btp_subaccount_entitlement.uut",
					ImportStateIdFunc: getImportStateIdForSubaccountEntitlement("btp_subaccount_entitlement.uut", "data-privacy-integration-service", "standard"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_entitlement.update")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEntitlementWithAmountBySubaccount("uut", "integration-test-acc-static", "data-privacy-integration-service", "standard", "1"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_name", "standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "service_name", "data-privacy-integration-service"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "amount", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "state", "OK"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEntitlementWithAmountBySubaccount("uut", "integration-test-acc-static", "data-privacy-integration-service", "standard", "2"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_name", "standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "service_name", "data-privacy-integration-service"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "amount", "2"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "state", "OK"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountEntitlementBySubaccount("uut", "integration-test-acc-static", "data-privacy-integration-service", "standard"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_entitlement.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_name", "standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "plan_id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "service_name", "data-privacy-integration-service"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "amount", "2"),
						resource.TestCheckResourceAttr("btp_subaccount_entitlement.uut", "state", "OK"),
					),
				},
			},
		})
	})

	t.Run("error path - zero amount", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountEntitlementWithAmountBySubaccount("uut", "integration-test-acc-static", "data-privacy-integration-service", "standard", "0"),
					ExpectError: regexp.MustCompile(`Attribute amount value must be between 1 and 2000000000, got: 0`),
				},
			},
		})
	})

}

func hclResourceSubaccountEntitlementBySubaccount(resourceName string, subaccountName string, serviceName string, planName string) string {
	template := `
data "btp_subaccounts" "all" {}
resource "btp_subaccount_entitlement" "%s" {
    subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
    service_name  = "%s"
    plan_name     = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, serviceName, planName)
}

func hclResourceSubaccountEntitlementWithAmountBySubaccount(resourceName string, subaccountName string, serviceName string, planName string, amount string) string {
	return fmt.Sprintf(`
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_entitlement" "%s" {
        subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
        service_name  = "%s"
        plan_name     = "%s"
        amount        = %s
    }`, resourceName, subaccountName, serviceName, planName, amount)
}

func getImportStateIdForSubaccountEntitlement(resourceName string, serviceName string, planName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s,%s", rs.Primary.Attributes["subaccount_id"], serviceName, planName), nil
	}
}
