package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGlobalaccountTrustConfiguration(t *testing.T) {

	var testIdp = getenv("BTP_TEST_IDP", "terraformtest.accounts400.ondemand.com")

	t.Run("happy path - minimal configuration with update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration.minimal")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationMinimal("uut", testIdp),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "identity_provider", testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "domain", testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "name", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "description", "Custom Platform Identity Provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "origin", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "id", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "type", "Platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "read_only", "false"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationComplete("uut", testIdp, testIdp, "Custom platform IAS tenant", "terraformtest-platform", "Description for "+testIdp),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "identity_provider", testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "domain", testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "name", "Custom platform IAS tenant"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "description", "Description for "+testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "origin", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "id", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "type", "Platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "read_only", "false"),
					),
				},
				{
					ResourceName:      "btp_globalaccount_trust_configuration.uut",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - complete configuration without update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration.complete")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationComplete("uut", testIdp, testIdp, "Custom platform IAS tenant", "terraformtest-platform", "Description for "+testIdp),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "identity_provider", testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "domain", testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "name", "Custom platform IAS tenant"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "description", "Description for "+testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "origin", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "id", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "type", "Platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "read_only", "false"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationMinimal("uut", testIdp),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "identity_provider", testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "domain", testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "name", "Custom platform IAS tenant"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "description", "Description for "+testIdp),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "origin", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "id", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "type", "Platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "read_only", "false"),
					),
				},
				{
					ResourceName:      "btp_globalaccount_trust_configuration.uut",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("error path - idp does not exist", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration.invalid_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationMinimal("uut", "this.is.not.a.valid.idp.accounts400.ondemand.com"),
					ExpectError: regexp.MustCompile(`the backend responded with an unknown error: 404`),
					// TODO: The actual error message from XSUAA is swallowed by the error handling in client.go
				},
			},
		})
	})

	t.Run("error path - malformed origin", func(t *testing.T) {

		rec, user := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration.invalid_origin")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationComplete("uut", testIdp, testIdp, "Custom platform IAS tenant", "invalid-platform-origin", "Description for "+testIdp),
					ExpectError: regexp.MustCompile(`Attribute origin must end with '-platform' and not exceed 36 characters`),
				},
			},
		})
	})

	// TODO: add "happy path - recreate when origin is changed", see NGPBUG-366076
	// TODO: add "error path - trust config does already exist", see NGPBUG-362964
}

func hclResourceGlobalaccountTrustConfigurationMinimal(resourceName string, identityProvider string) string {
	return fmt.Sprintf(`resource "btp_globalaccount_trust_configuration" "%s" { identity_provider = "%s" }`, resourceName, identityProvider)
}

func hclResourceGlobalaccountTrustConfigurationComplete(resourceName string, identityProvider string, domain string, name string, origin string, description string) string {
	return fmt.Sprintf(`resource "btp_globalaccount_trust_configuration" "%s" {
		identity_provider        = "%s"
	    domain                   = "%s"
		name                     = "%s"
		origin                   = "%s"
		description              = "%s"
}`, resourceName, identityProvider, domain, name, origin, description)
}
