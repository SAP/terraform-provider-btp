package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountIdentityProvider(t *testing.T) {
	t.Parallel()
	t.Run("happy path - get global idp by host", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_identity_provider.by_host")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountIdentityProviderByHost("uut", "a2dynbhnd.accounts400.ondemand.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_identity_provider.uut", "host", "a2dynbhnd.accounts400.ondemand.com"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_identity_provider.uut", "tenant_id", "a2dynbhnd"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_identity_provider.uut", "status", "ACTIVE"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_identity_provider.uut", "tenant_type", "trial"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_identity_provider.uut", "common_host", "a2dynbhnd.accounts400.cloud.sap"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_identity_provider.uut", "data_center_id", "EU2_QA"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_identity_provider.uut", "cost_center_id", "101000297"),
					),
				},
			},
		})
	})

	t.Run("error path - host mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: `data "btp_globalaccount_identity_provider" "uut" {

}`,
					ExpectError: regexp.MustCompile(`The argument "host" is required`),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountIdentityProviderByHost(resourceName string, host string) string {
	return fmt.Sprintf(`
data "btp_globalaccount_identity_provider" "%s" {
    host = "%s"
}`, resourceName, host)
}
