package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountSecuritySettings(t *testing.T) {

	t.Parallel()
	t.Run("happy path - security settings by id", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_security_settings")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountSecuritySettingbyId("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_security_settings.uut", "custom_email_domains.#", "0"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_security_settings.uut", "default_identity_provider", "sap.default"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_security_settings.uut", "treat_users_with_same_email_as_same_user", "false"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_security_settings.uut", "access_token_validity", "-1"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_security_settings.uut", "refresh_token_validity", "-1"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_security_settings.uut", "iframe_domains", ""),
					),
				},
			},
		})

	})
}

func hclDatasourceGlobalaccountSecuritySettingbyId(resourceName string) string {
	template := `data "btp_globalaccount_security_settings" "%s" { }`
	return fmt.Sprintf(template, resourceName)
}
