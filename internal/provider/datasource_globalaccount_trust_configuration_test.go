package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountTrustConfiguration(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_trust_configuration.default")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountTrustConfiguration("uut", "sap.default"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "id", "sap.default"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "description", ""),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "identity_provider", ""),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "name", "sap.default"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "read_only", "false"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "type", "Application"),
					),
				},
			},
		})
	})
	t.Run("happy path - custom idp - existing", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_trust_configuration.custom_idp_exists")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountTrustConfiguration("uut", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "id", "terraformint-platform"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "description", "Custom Platform Identity Provider"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "identity_provider", "terraformint.accounts400.ondemand.com"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "name", "terraformint-platform"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "read_only", "false"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "type", "Platform"),
					),
				},
			},
		})
	})
	t.Run("happy path - custom idp - not existing", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_trust_configuration.custom_idp_not_existing")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountTrustConfiguration("uut", "fuh"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "id", ""),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "description", ""),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "identity_provider", ""),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "name", ""),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "protocol", ""),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "read_only", "false"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "status", ""),
						resource.TestCheckResourceAttr("data.btp_globalaccount_trust_configuration.uut", "type", ""),
					),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountTrustConfiguration(resourceName string, origin string) string {
	template := `data "btp_globalaccount_trust_configuration" "%s" { origin = "%s" }`

	return fmt.Sprintf(template, resourceName, origin)
}
