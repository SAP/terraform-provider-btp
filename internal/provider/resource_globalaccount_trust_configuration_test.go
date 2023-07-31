package provider

/* TODO This test is capable of deleting the trust to the terraformint IAS tenant. It was unexpectedly able to create
        the trust when the limitation to one trust was dropped. This was anyway not the error it was targeting at.
        The test was disabled and needs to be adapted to work with an IAS tenant that is not required for development.

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
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration.exists")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationSimple("uut", "terraformint.accounts400.ondemand.com"),
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

*/
