package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceSubaccountTrustConfiguration(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount_trust_configuration.simple_default_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceSubaccountTrustConfigurationSimple("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "sap.default"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "parent_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", "john.doe@int.test"),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "UNSET"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "beta_enabled", "false"),
					),
				},
				{
					Config: hclProvider() + hclResourceSubaccountTrustConfigurationSimple("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "parent_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", "john.doe@int.test"),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "UNSET"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "beta_enabled", "false"),
					),
				},
				{
					Config: hclProvider() + hclResourceSubaccountTrustConfigurationSimple("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "sap.default"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
					),
				},
				{
					ResourceName:      "btp_subaccount.uut",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

}

func hclResourceSubaccountTrustConfigurationSimple(resourceName string, subaccountId string, identityProvider string) string {
	return fmt.Sprintf(`resource "btp_subaccount_trust_configuration" "%s" {
        subaccount_id		= "%s"
        identity_provider	= "%s"
    }`, resourceName, subaccountId, identityProvider)
}

/*
func hclResourceSubaccountTrustConfigurationFullyCustomized(resourceName string, subaccountId string, identityProvider string, name string, description string, origin string) string {
	return fmt.Sprintf(`resource "btp_subaccount_trust_configuration" "%s" {
        subaccount_id      	= "%s"
        identity_provider   = "%s"
        name				= "%s"
        description    		= "%s"
        origin 				= "%s"
		}`, resourceName, subaccountId, identityProvider, name, description, origin)
}
*/
