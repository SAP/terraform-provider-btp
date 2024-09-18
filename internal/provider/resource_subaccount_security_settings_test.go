package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceSubaccountSecuritySettings(t *testing.T) {
	t.Parallel()

	t.Run("happy path - complete configuration", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_security_settings.complete")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountSecuritySettings("uut", "integration-test-security-settings", "terraformint-platform", 3601, 3602, true, "[\"domain1.test\",\"domain2.test\"]", "iframedomain.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_security_settings.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "access_token_validity", "3601"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "refresh_token_validity", "3602"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "default_identity_provider", "terraformint-platform"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.#", "2"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.1", "domain2.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "access_token_validity", "4000"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_doamins", "iframedomain.test"),
					),
				},

				{
					Config: hclProviderFor(user) + hclResourceSubaccountSecuritySettings("uut", "integration-test-security-settings", "terraformint-platform", 4000, 3602, false, "[\"domain1.test\"]", "updatediframedomain.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_security_settings.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "access_token_validity", "4000"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "refresh_token_validity", "3602"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "default_identity_provider", "terraformint-platform"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "false"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.#", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_doamins", "updatediframedomain.test"),
					),
				},
			},
		})
	})
}

func hclResourceSubaccountSecuritySettings(resourceName string, subaccountName string, defaultIdp string, accessTokenValidity int, refreshTokenValidity int, treatUsersWithSameEmailAsSameUser bool, customEmailDomains string, iFrameDomain string) string {
	template := `
data "btp_subaccounts" "all" {}
resource "btp_subaccount_security_settings" "%s" {
    subaccount_id            = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]

    default_identity_provider = "%s"

    access_token_validity  = %v
    refresh_token_validity = %v

    treat_users_with_same_email_as_same_user = %v

    custom_email_domains = %v

		iframe_domain = %s
}`
	return fmt.Sprintf(template, resourceName, subaccountName, defaultIdp, accessTokenValidity, refreshTokenValidity, treatUsersWithSameEmailAsSameUser, customEmailDomains, iFrameDomain)
}
