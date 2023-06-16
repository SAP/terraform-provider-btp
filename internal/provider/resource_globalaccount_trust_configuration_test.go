package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGlobalaccountTrustConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("happy path - complete configuration", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration.complete")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceGlobalaccountTrustConfigurationSimple("uut", "terraformint.accounts400.ondemand.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_globalaccount_trust_configuration.uut", "subaccount_id", regexpValidUUID),
						//resource.TestMatchResourceAttr("btp_globalaccount_trust_configuration.uut", "created_date", regexpValidRFC3999Format),
						//resource.TestMatchResourceAttr("btp_globalaccount_trust_configuration.uut", "last_modified", regexpValidRFC3999Format),
						//resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "id", "sap.custom"),

					),
				},
			},
		})
	})

}

func hclResourceGlobalaccountTrustConfigurationSimple(resourceName string, identityProvider string) string {
	template := `
resource "btp_globalaccount_trust_configuration" "%s" {
    identity_provider  	= "%s"
}`

	return fmt.Sprintf(template, resourceName, identityProvider)
}
