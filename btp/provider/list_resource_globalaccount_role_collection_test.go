package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
	res "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/querycheck/queryfilter"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestGlobalaccountRoleCollectionListResource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_globalaccount_role_collection")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			Steps: []resource.TestStep{
				{
					Query:  true,
					Config: hclProviderFor(user) + listGlobalAccountRoleCollectionQueryConfig("role_collection_list", "btp"),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_globalaccount_role_collection.role_collection_list", 2),

						querycheck.ExpectIdentity(
							"btp_globalaccount_role_collection.role_collection_list",
							map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Global Account Viewer"),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query:  true,
					Config: hclProviderFor(user) + listGlobalAccountRoleCollectionQueryConfigWithIncludeResource("role_collection_list", "btp"),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_globalaccount_role_collection.role_collection_list", 2),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_globalaccount_role_collection.role_collection_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("Global Account Viewer"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("description"),
									KnownValue: knownvalue.StringExact("Read-only access to the global account"),
								},
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("Global Account Viewer"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringExact("Global Account Viewer"),
								},
								{
									Path:       tfjsonpath.New("roles"),
									KnownValue: knownvalue.ListSizeExact(6),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewGlobalaccountRoleCollectionListResource().(list.ListResourceWithConfigure)
		resp := &res.ConfigureResponse{}
		req := res.ConfigureRequest{
			ProviderData: struct{}{}, // Wrong type
		}

		r.Configure(context.Background(), req, resp)

		if !resp.Diagnostics.HasError() {
			t.Error("Expected error for invalid provider data type")
		}
	})

}

func listGlobalAccountRoleCollectionQueryConfig(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_globalaccount_role_collection" "%s" {
  provider = "%s"
}`, label, providerName)
}

func listGlobalAccountRoleCollectionQueryConfigWithIncludeResource(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_globalaccount_role_collection" "%s" {
  provider = "%s"
  include_resource = true
}`, label, providerName)
}
