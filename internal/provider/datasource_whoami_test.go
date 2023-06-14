package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceWhoami(t *testing.T) {
	t.Parallel()
	t.Run("happy path with default idp", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_whoami")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceWhoami("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_whoami.uut", "id", "john.doe@int.test"),
						resource.TestCheckResourceAttr("data.btp_whoami.uut", "email", "john.doe@int.test"),
						resource.TestCheckResourceAttr("data.btp_whoami.uut", "issuer", "accounts.sap.com"),
					),
				},
			},
		})
	})
}

func hclDatasourceWhoami(resourceName string) string {
	return fmt.Sprintf(`data "btp_whoami" "%s" {}`, resourceName)
}
