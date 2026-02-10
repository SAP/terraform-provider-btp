package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountIdentityProviders(t *testing.T) {
	t.Parallel()
	t.Run("happy path - list all global idps", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_identity_providers.list_all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + `data "btp_globalaccount_identity_providers" "uut" {}`,
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_identity_providers.uut", "values.#", "12"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_identity_providers.uut", "values.0.status", "ACTIVE"),
					),
				},
			},
		})
	})
}
