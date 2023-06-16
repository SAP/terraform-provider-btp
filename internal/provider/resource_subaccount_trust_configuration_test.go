package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceSubaccountTrustConfiguration(t *testing.T) {
	t.Parallel()
	/*
		t.Run("happy path - complete configuration 1", func(t *testing.T) {
			rec := setupVCR(t, "fixtures/resource_subaccount_trust_configuration.complete1")
			defer stopQuietly(rec)

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
				Steps: []resource.TestStep{
					{
						Config: hclProvider() + hclResourceSubaccountTrustConfigurationComplete("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "terraformint.accounts400.ondemand.com", "terraform1-platform", "Custom Platform Identity Provider 1 for testing", "terraform1-tst-platform"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestMatchResourceAttr("btp_subaccount_trust_configuration.uut", "subaccount_id", regexpValidUUID),
							resource.TestMatchResourceAttr("btp_subaccount_trust_configuration.uut", "created_date", regexpValidRFC3999Format),
							resource.TestMatchResourceAttr("btp_subaccount_trust_configuration.uut", "last_modified", regexpValidRFC3999Format),
							resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "id", "hana-cloud-hana"),
							resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "identity_provider", "hana-cloud-hana"),
							resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "name", "hana-cloud-hana"),
							resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "description", "hana-cloud-hana"),
							resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "origin", "hana-cloud-hana"),
						),
					},
				},
			})
		}) */

	t.Run("happy path - minimal configuration 1", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount_trust_configuration.minimum1")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceSubaccountTrustConfigurationMinimum("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "terraformint.accounts400.ondemand.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_trust_configuration.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_trust_configuration.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_trust_configuration.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "id", "sap.custom"),
						resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "identity_provider", "terraformint.accounts400.ondemand.com"),
						resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "name", "Custom IAS tenant"),
						resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "description", "IAS tenant terraformint.accounts400.ondemand.com (OpenID Connect)"),
						resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "origin", "sap.custom"),
					),
				},
			},
		})
	})

	/*
		t.Run("happy path - minimal configuration 2", func(t *testing.T) {
			rec := setupVCR(t, "fixtures/resource_subaccount_trust_configuration.minimum2")
			defer stopQuietly(rec)

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
				Steps: []resource.TestStep{
					{
						Config: hclProvider() + hclResourceSubaccountTrustConfigurationMinimum("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "terraformint.accounts400.ondemand.com"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestMatchResourceAttr("btp_subaccount_trust_configuration.uut", "subaccount_id", regexpValidUUID),
							resource.TestMatchResourceAttr("btp_subaccount_trust_configuration.uut", "created_date", regexpValidRFC3999Format),
							resource.TestMatchResourceAttr("btp_subaccount_trust_configuration.uut", "last_modified", regexpValidRFC3999Format),
							resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "id", "hana-cloud-hana"),
							resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "identity_provider", "hana-cloud-hana"),
							resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "name", "hana-cloud-hana"),
							resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "description", "hana-cloud-hana"),
							resource.TestCheckResourceAttr("btp_subaccount_trust_configuration.uut", "origin", "hana-cloud-hana"),
						),
					},
				},
			})
		})
	*/
}

func hclResourceSubaccountTrustConfigurationComplete(resourceName string, subaccountId string, identityProvider string, name string, description string, origin string) string {
	template := `
resource "btp_subaccount_trust_configuration" "%s" {
    subaccount_id 		= "%s"
    identity_provider  	= "%s"
    name     			= "%s"
    description     	= "%s"
    origin     			= "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId, identityProvider, name, description, origin)
}

func hclResourceSubaccountTrustConfigurationMinimum(resourceName string, subaccountId string, identityProvider string) string {
	template := `
resource "btp_subaccount_trust_configuration" "%s" {
    subaccount_id 		= "%s"
    identity_provider  	= "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId, identityProvider)
}
