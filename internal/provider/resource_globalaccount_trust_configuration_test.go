package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGlobalaccountTrustConfiguration(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationSimple("uut", "terraformtest.accounts400.ondemand.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "origin", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "name", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "description", "Custom Platform Identity Provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "type", "Platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "identity_provider", "terraformtest.accounts400.ondemand.com"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "read_only", "false"),
					),
				},
			},
		})
	})

	t.Run("happy path - full spec", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration.full_spec")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationFullSpec("uut", "terraformtest.accounts400.ondemand.com", "terraformtest-name", "terraformtest-platform", "Terraform Test Identity Provider"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "origin", "terraformtest-platform"), // TODO: use an alternate origin after NGPBUG-363131 is resolved
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "name", "terraformtest-name"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "description", "Terraform Test Identity Provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "type", "Platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "identity_provider", "terraformtest.accounts400.ondemand.com"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "read_only", "false"),
					),
				},
			},
		})
	})

	t.Run("happy path - update", func(t *testing.T) {

		t.Skip("test is skipped until update is implemented")

		rec, user := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationSimple("uut", "terraformtest.accounts400.ondemand.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "origin", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "name", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "description", "Custom Platform Identity Provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "type", "Platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "identity_provider", "terraformtest.accounts400.ondemand.com"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "read_only", "false"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationFullSpec("uut", "terraformtest.accounts400.ondemand.com", "terraformtest-name", "terraformtest-platform", "Terraform Test Identity Provider"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "origin", "terraformtest-platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "name", "terraformtest-name"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "description", "Terraform Test Identity Provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "type", "Platform"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "identity_provider", "terraformtest.accounts400.ondemand.com"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("btp_globalaccount_trust_configuration.uut", "read_only", "false"),
					),
				},
			},
		})
	})

	t.Run("error path - idp does not exist", func(t *testing.T) {

		t.Skip("test is skipped until NGPBUG-363098 is resolved")

		rec, user := setupVCR(t, "fixtures/resource_globalaccount_trust_configuration.invalid_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationSimple("uut", "this.is.not.a.valid.idp.accounts400.ondemand.com"),
					ExpectError: regexp.MustCompile(`tbd`),
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
					Config:      hclProviderFor(user) + hclResourceGlobalaccountTrustConfigurationFullSpec("uut", "terraformtest.accounts400.ondemand.com", "terraformtest-name", "invalid-platform-origin", "Terraform Test Identity Provider"),
					ExpectError: regexp.MustCompile(`Attribute origin must end with '-platform' and not exceed 36 characters`),
				},
			},
		})
	})

	// TODO: add "error path - trust config does already exist", see NGPBUG-362964
}

func hclResourceGlobalaccountTrustConfigurationSimple(resourceName string, identityProvider string) string {
	return fmt.Sprintf(`resource "btp_globalaccount_trust_configuration" "%s" { identity_provider = "%s" }`, resourceName, identityProvider)
}

func hclResourceGlobalaccountTrustConfigurationFullSpec(resourceName string, identityProvider string, name string, origin string, description string) string {
	return fmt.Sprintf(`resource "btp_globalaccount_trust_configuration" "%s" {
		identity_provider = "%s"
		name              = "%s"
		origin            = "%s"
		description       = "%s"
}`, resourceName, identityProvider, name, origin, description)
}
