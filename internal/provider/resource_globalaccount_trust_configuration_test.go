package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGlobalaccountTrustConfiguration(t *testing.T) {
	t.Parallel()

	// TODO: we need an additional test in a global account
	// without a configured trust configuration

	t.Run("happy path - trust config exists", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration.exists")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceGlobalaccountTrustConfigurationSimple("uut", "terraformint.accounts400.ondemand.com"),
					// TODO: we need to work on a good error message for an account that already has a trust configuration
					ExpectError: regexp.MustCompile(`invalid character 'L' looking for beginning of value`),
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
