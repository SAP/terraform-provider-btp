package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceWhoami(t *testing.T) {
	t.Parallel()
	t.Run("happy path with default idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_whoami")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceWhoami("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_whoami.uut", "id", user.Username),
						resource.TestCheckResourceAttr("data.btp_whoami.uut", "email", user.Username),
						resource.TestCheckResourceAttr("data.btp_whoami.uut", "issuer", user.Issuer),
					),
				},
			},
		})
	})
}

func hclDatasourceWhoami(resourceName string) string {
	return fmt.Sprintf(`data "btp_whoami" "%s" {}`, resourceName)
}
