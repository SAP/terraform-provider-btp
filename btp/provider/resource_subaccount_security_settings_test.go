package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceSubaccountSecuritySettings(t *testing.T) {
	t.Run("happy path - complete configuration", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_security_settings.complete")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountSecuritySettings("uut", "integration-test-security-settings", "terraformint-platform", 3601, 3602, true, "[\"domain1.test\",\"domain2.test\"]", "https://iframedomain.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_security_settings.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "access_token_validity", "3601"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "refresh_token_validity", "3602"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "default_identity_provider", "terraformint-platform"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.#", "2"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.1", "domain2.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains", "https://iframedomain.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains_list.#", "1"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountSecuritySettings("uut", "integration-test-security-settings", "terraformint-platform", 4000, 3602, false, "[\"domain1.test\"]", "https://iframedomain.test https://updated.iframedomain.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_security_settings.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "access_token_validity", "4000"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "refresh_token_validity", "3602"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "default_identity_provider", "terraformint-platform"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "false"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.#", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains", "https://iframedomain.test https://updated.iframedomain.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains_list.#", "2"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains_list.1", "https://updated.iframedomain.test"),
					),
				},
				{
					ResourceName:      "btp_subaccount_security_settings.uut",
					ImportStateIdFunc: getSecuritySettingsImportStateId("btp_subaccount_security_settings.uut"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - complete configuration with iframe domains list", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_security_settings_with_iframe_domains_list.complete")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountSecuritySettingsWithIFrameDomainsList("uut", "integration-test-security-settings", "terraformint-platform", 4000, 3602, false, "[\"domain1.test\"]", "[\"https://iframedomainlist.test\"]"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_security_settings.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "access_token_validity", "4000"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "refresh_token_validity", "3602"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "default_identity_provider", "terraformint-platform"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "false"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.#", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains", "https://iframedomainlist.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains_list.#", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains_list.0", "https://iframedomainlist.test"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountSecuritySettingsWithIFrameDomainsList("uut", "integration-test-security-settings", "terraformint-platform", 4000, 3602, false, "[\"domain1.test\"]", "[\"https://iframedomainlist.test\",\"https://updated.iframedomainlist.test\"]"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_security_settings.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "access_token_validity", "4000"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "refresh_token_validity", "3602"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "default_identity_provider", "terraformint-platform"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "false"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.#", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains", "https://iframedomainlist.test https://updated.iframedomainlist.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains_list.#", "2"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains_list.0", "https://iframedomainlist.test"),
					),
				},
				{
					ResourceName:      "btp_subaccount_security_settings.uut",
					ImportStateIdFunc: getSecuritySettingsImportStateId("btp_subaccount_security_settings.uut"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - IFrame deletion", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_security_settings.destroy")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountSecuritySettings("uut", "integration-test-security-settings", "terraformint-platform", 3601, 3602, true, "[\"domain1.test\",\"domain2.test\"]", "https://iframedomain.test"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_security_settings.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "access_token_validity", "3601"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "refresh_token_validity", "3602"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "default_identity_provider", "terraformint-platform"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.#", "2"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.1", "domain2.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains", "https://iframedomain.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains_list.#", "1"),
					),
				},

				{
					Config: hclProviderFor(user) + hclResourceSubaccountSecuritySettings("uut", "integration-test-security-settings", "terraformint-platform", 4000, 3602, false, "[\"domain1.test\"]", ""),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_security_settings.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "access_token_validity", "4000"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "refresh_token_validity", "3602"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "default_identity_provider", "terraformint-platform"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "false"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.0", "domain1.test"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "custom_email_domains.#", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains", ""),
						resource.TestCheckResourceAttr("btp_subaccount_security_settings.uut", "iframe_domains_list.#", "0"),
					),
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
					Config:      hclResourceSubaccountSecuritySettings("uut", "integration-test-security-settings", "terraformint-platform", 4000, 3602, false, "[\"domain1.test\"]", " "),
					ExpectError: regexp.MustCompile(`Attribute iframe_domains The attribute iframe_domains must be empty`),
				},
			},
		})
	})

	t.Run("error path - both iframe_domains and iframe_domains_list is defined in configuration", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountSecuritySettingsWithIFrameDomainsListAndIframeDomains("uut", "integration-test-security-settings", "terraformint-platform", 4000, 3602, false, "[\"domain1.test\"]", "https://iframedomain.test", "[\"https://iframedomain.test\"]"),
					ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
				},
			},
		})
	})
}

func hclResourceSubaccountSecuritySettings(resourceName string, subaccountName string, defaultIdp string, accessTokenValidity int, refreshTokenValidity int, treatUsersWithSameEmailAsSameUser bool, customEmailDomains string, iFrameDomains string) string {
	template := `
data "btp_subaccounts" "all" {}
resource "btp_subaccount_security_settings" "%s" {
    subaccount_id            = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]

    default_identity_provider = "%s"

    access_token_validity  = %v
    refresh_token_validity = %v

    treat_users_with_same_email_as_same_user = %v

    custom_email_domains = %v

		iframe_domains = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, defaultIdp, accessTokenValidity, refreshTokenValidity, treatUsersWithSameEmailAsSameUser, customEmailDomains, iFrameDomains)
}

func hclResourceSubaccountSecuritySettingsWithIFrameDomainsListAndIframeDomains(resourceName string, subaccountName string, defaultIdp string, accessTokenValidity int, refreshTokenValidity int, treatUsersWithSameEmailAsSameUser bool, customEmailDomains string, iFrameDomains string, iFrameDomainsList string) string {
	template := `
data "btp_subaccounts" "all" {}
resource "btp_subaccount_security_settings" "%s" {
    subaccount_id            = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]

    default_identity_provider = "%s"

    access_token_validity  = %v
    refresh_token_validity = %v

    treat_users_with_same_email_as_same_user = %v

    custom_email_domains = %v

		iframe_domains = "%s"
		iframe_domains_list = %v
}`
	return fmt.Sprintf(template, resourceName, subaccountName, defaultIdp, accessTokenValidity, refreshTokenValidity, treatUsersWithSameEmailAsSameUser, customEmailDomains, iFrameDomains, iFrameDomainsList)
}

func hclResourceSubaccountSecuritySettingsWithIFrameDomainsList(resourceName string, subaccountName string, defaultIdp string, accessTokenValidity int, refreshTokenValidity int, treatUsersWithSameEmailAsSameUser bool, customEmailDomains string, iFrameDomainsList string) string {
	template := `
data "btp_subaccounts" "all" {}
resource "btp_subaccount_security_settings" "%s" {
    subaccount_id            = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]

    default_identity_provider = "%s"

    access_token_validity  = %v
    refresh_token_validity = %v

    treat_users_with_same_email_as_same_user = %v

    custom_email_domains = %v
		iframe_domains_list = %v
}`
	return fmt.Sprintf(template, resourceName, subaccountName, defaultIdp, accessTokenValidity, refreshTokenValidity, treatUsersWithSameEmailAsSameUser, customEmailDomains, iFrameDomainsList)
}

func getSecuritySettingsImportStateId(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.Attributes["subaccount_id"], nil
	}
}
