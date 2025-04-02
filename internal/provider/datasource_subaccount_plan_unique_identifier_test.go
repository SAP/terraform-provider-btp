package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountEntitlementUniqueIdentifier(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_entitlement_unique_identifier")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountEntitlementUniqueIdentifier("uut", "integration-test-acc-static", "hana-cloud", "hana"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_entitlement_unique_identifier.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_entitlement_unique_identifier.uut", "plan_unique_identifier"),
					),
				},
			},
		})
	})

	t.Run("error path - subaccount_id is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: `
data "btp_subaccount_entitlement_unique_identifier" "uut" {
  service_name = "hana-cloud"
  plan_name    = "hana"
}
`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required`),
				},
			},
		})
	})

	t.Run("error path - invalid service or plan", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_entitlement_unique_identifier_invalid")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceSubaccountEntitlementUniqueIdentifier("uut", "integration-test-acc-static", "invalid-service", "invalid-plan"),
					ExpectError: regexp.MustCompile(`Could not find service 'invalid-service' with plan 'invalid-plan'`),
				},
			},
		})
	})
}
func hclDatasourceSubaccountEntitlementUniqueIdentifier(resourceName, subaccountId, serviceName, planName string) string {
	return fmt.Sprintf(`
data "btp_subaccount_entitlement_unique_identifier" "%s" {
  subaccount_id = "%s"
  service_name  = "%s"
  plan_name     = "%s"
}
`, resourceName, subaccountId, serviceName, planName)
}
