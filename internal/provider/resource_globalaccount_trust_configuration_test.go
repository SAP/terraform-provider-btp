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

	t.Run("error path - trust config does already exist", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration.exists")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + hclResourceGlobalaccountTrustConfigurationSimple("uut", "terraformint.accounts400.ondemand.com"),
					ExpectError: regexp.MustCompile(`the backend responded with an unknown error: 400`), //FIXME NGPBUG-350117
				},
			},
		})
	})

}

func hclResourceGlobalaccountTrustConfigurationSimple(resourceName string, identityProvider string) string {
	template := `resource "btp_globalaccount_trust_configuration" "%s" { identity_provider = "%s" }`

	return fmt.Sprintf(template, resourceName, identityProvider)
}
