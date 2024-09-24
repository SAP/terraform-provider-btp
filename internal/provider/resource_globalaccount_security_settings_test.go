package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGlobalaccountSecuritySettings(t *testing.T) {
	t.Run("happy path - complete configuration", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_security_settings.complete")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountSecuritySettings("uut", "terraformint-platform", 4500, 4500, true, "[\"domain1.test\",\"domain2.test\"]", "https://iframedomain.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "access_token_validity", "4500"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "refresh_token_validity", "4500"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "true"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "custom_email_domains.#", "2"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "custom_email_domains.1", "domain2.test"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "iframe_domains", "https://iframedomain.test"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountSecuritySettings("uut", "terraformint-platform", 4500, 4500, false, "[\"domain1.test\"]", "https://iframedomain.test https://updated.iframedomain.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "access_token_validity", "4500"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "refresh_token_validity", "4500"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "false"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "custom_email_domains.#", "1"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "iframe_domains", "https://iframedomain.test https://updated.iframedomain.test"),
					),
				},
				{
					ResourceName:      "btp_globalaccount_security_settings.uut",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - IFrame deletion", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_security_settings.destroy")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountSecuritySettings("uut", "terraformint-platform", 4500, 4500, true, "[\"domain1.test\",\"domain2.test\"]", "https://iframedomain.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "access_token_validity", "4500"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "refresh_token_validity", "4500"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "true"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "custom_email_domains.#", "2"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "custom_email_domains.1", "domain2.test"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "iframe_domains", "https://iframedomain.test"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountSecuritySettings("uut", "sap.default", 4500, 4500, false, "[]", ""),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "access_token_validity", "4500"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "refresh_token_validity", "4500"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "false"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "custom_email_domains.#", "0"),
						resource.TestCheckResourceAttr("btp_globalaccount_security_settings.uut", "iframe_domains", ""),
					),
				},
				{
					ResourceName:      "btp_globalaccount_security_settings.uut",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("error path - invalid iframe value", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceGlobalaccountSecuritySettings("uut", "sap.default", 4500, 4500, false, "[]", " "),
					ExpectError: regexp.MustCompile(`Attribute iframe_domains The attribute iframe_domains must be empty`),
				},
			},
		})
	})
}

func hclResourceGlobalaccountSecuritySettings(resourceName string, defaultIdp string, accessTokenValidity int, refreshTokenValidity int, treatUsersWithSameEmailAsSameUser bool, customEmailDomains string, iFrameDomains string) string {
	template := `
resource "btp_globalaccount_security_settings" "%s" {
    default_identity_provider = "%s"

    access_token_validity  = %v
    refresh_token_validity = %v

    treat_users_with_same_email_as_same_user = %v

    custom_email_domains = %v

		iframe_domains = "%s"
}`

	return fmt.Sprintf(template, resourceName, defaultIdp, accessTokenValidity, refreshTokenValidity, treatUsersWithSameEmailAsSameUser, customEmailDomains, iFrameDomains)
}
