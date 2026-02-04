package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountRoleCollectionBase(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_role_collection_base")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountRoleCollectionBase("uut", "integration-test-acc-static", "Subaccount Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_role_collection_base.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_role_collection_base.uut", "name", "Subaccount Viewer"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_role_collection_base.uut", "description"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role_collection_base.uut", "read_only", "true"),
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
					Config: `
data "btp_subaccount_role_collection_bases" "uut" {
    # subaccount_id is missing
}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required`),
				},
			},
		})
	})

	t.Run("error path - name must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: `
data "btp_subaccount_role_collection_base" "uut" {
    subaccount_id = "00000000-0000-0000-0000-000000000000"
    name          = ""
}`,
					ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountRoleCollectionBase(resourceName string, subaccountName string, collectionName string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
data "btp_subaccount_role_collection_base" "%s" {
    subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
    name          = "%s"
}`, resourceName, subaccountName, collectionName)
}
