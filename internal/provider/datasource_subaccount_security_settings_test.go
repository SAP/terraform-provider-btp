package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountSecuritySettings(t *testing.T) {

	t.Parallel()
	t.Run("happy path - security settings by id", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_security_settings.by_id")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountSecuritySettingbyId("uut", "integration-test-services-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_security_settings.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_security_settings.uut", "custom_email_domains.#", "0"),
						resource.TestCheckResourceAttr("data.btp_subaccount_security_settings.uut", "default_identity_provider", "sap.default"),
						resource.TestCheckResourceAttr("data.btp_subaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "false"),
						resource.TestCheckResourceAttr("data.btp_subaccount_security_settings.uut", "access_token_validity", "-1"),
						resource.TestCheckResourceAttr("data.btp_subaccount_security_settings.uut", "refresh_token_validity", "-1"),
					),
				},
			},
		})

	})

	t.Run("error path - subaccount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountSecuritySettingNoSubaccount("uut"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountSecuritySettingbyId(resourceName string, subaccountName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_security_settings" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName)
}

func hclDatasourceSubaccountSecuritySettingNoSubaccount(resourceName string) string {
	template := `data "btp_subaccount_security_settings" "%s" {
	}`
	return fmt.Sprintf(template, resourceName)
}
