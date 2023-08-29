package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceDirectoryEntitlement(t *testing.T) {
	t.Parallel()
	t.Run("happy path - no amount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_entitlement.no_amount")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryEntitlement("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "hana-cloud", "hana"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "directory_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "id", "hana-cloud-hana"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "plan_name", "hana"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "plan_id", "hana-cloud-hana"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "service_name", "hana-cloud"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "amount", "1"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "state", "OK"),
					),
				},
				{
					ResourceName:      "btp_directory_entitlement.uut",
					ImportStateId:     "05368777-4934-41e8-9f3c-6ec5f4d564b9,hana-cloud,hana",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - with amount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_entitlement.amount_set")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryEntitlementWithAmount("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "data-privacy-integration-service", "standard", "3"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "directory_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "plan_name", "standard"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "plan_id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "service_name", "data-privacy-integration-service"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "amount", "3"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "state", "OK"),
					),
				},
				{
					ResourceName:      "btp_directory_entitlement.uut",
					ImportStateId:     "05368777-4934-41e8-9f3c-6ec5f4d564b9,data-privacy-integration-service,standard",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_entitlement.update")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryEntitlementWithAmount("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "data-privacy-integration-service", "standard", "1"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "directory_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "plan_name", "standard"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "plan_id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "service_name", "data-privacy-integration-service"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "amount", "1"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "state", "OK"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceDirectoryEntitlementWithAmount("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "data-privacy-integration-service", "standard", "2"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "directory_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory_entitlement.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "plan_name", "standard"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "plan_id", "data-privacy-integration-service-standard"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "service_name", "data-privacy-integration-service"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "amount", "2"),
						resource.TestCheckResourceAttr("btp_directory_entitlement.uut", "state", "OK"),
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
					Config:      hclResourceDirectoryEntitlementWithAmount("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "data-privacy-integration-service", "standard", "0"),
					ExpectError: regexp.MustCompile(`Attribute amount value must be between 1 and 2000000000, got: 0`),
				},
			},
		})
	})
}

func hclResourceDirectoryEntitlement(resourceName string, directoryId string, serviceName string, planName string) string {
	template := `
resource "btp_directory_entitlement" "%s" {
    directory_id = "%s"
    service_name  = "%s"
    plan_name     = "%s"
}`

	return fmt.Sprintf(template, resourceName, directoryId, serviceName, planName)
}

func hclResourceDirectoryEntitlementWithAmount(resourceName string, directoryId string, serviceName string, planName string, amount string) string {
	return fmt.Sprintf(`resource "btp_directory_entitlement" "%s" {
        directory_id      = "%s"
        service_name    = "%s"
        plan_name = "%s"
        amount = %s
    }`, resourceName, directoryId, serviceName, planName, amount)
}
