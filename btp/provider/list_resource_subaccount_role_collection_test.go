package provider

import (
	"context"
	"fmt"
	"regexp"
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

func TestSubaccountRoleCollectionListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "ba268910-81e6-4ac1-9016-cae7ed196889"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_role_collection")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			Steps: []resource.TestStep{
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountRoleCollectionQueryConfig(
						"subaccount_role_collection_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role_collection.subaccount_role_collection_list", 7),

						querycheck.ExpectIdentity(
							"btp_subaccount_role_collection.subaccount_role_collection_list",
							map[string]knownvalue.Check{
								"name":          knownvalue.StringExact("Subaccount Administrator"),
								"subaccount_id": knownvalue.StringExact("ba268910-81e6-4ac1-9016-cae7ed196889"),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountRoleCollectionQueryConfigWithIncludeResource(
						"subaccount_role_collection_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role_collection.subaccount_role_collection_list", 7),

						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_role_collection.subaccount_role_collection_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"name":          knownvalue.StringExact("Subaccount Administrator"),
								"subaccount_id": knownvalue.StringExact("ba268910-81e6-4ac1-9016-cae7ed196889"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("description"),
									KnownValue: knownvalue.StringExact("Administrative access to the subaccount"),
								},
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("Subaccount Administrator"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringRegexp(regexp.MustCompile(`.+,Subaccount Administrator`)),
								},
								{
									Path:       tfjsonpath.New("subaccount_id"),
									KnownValue: knownvalue.StringExact("ba268910-81e6-4ac1-9016-cae7ed196889"),
								},
								{
									Path:       tfjsonpath.New("roles"),
									KnownValue: knownvalue.ListSizeExact(9),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - Access forbidden", func(t *testing.T) {
		notFoundSubaccountID := "00000000-0000-0000-0000-000000000000"

		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_role_collection_access_forbidden")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			Steps: []resource.TestStep{
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountRoleCollectionQueryConfig(
						"not_found_test",
						"btp",
						notFoundSubaccountID,
					),

					ExpectError: regexp.MustCompile(`Access forbidden due to insufficient authorization|403`),
				},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountRoleCollectionListResource().(list.ListResourceWithConfigure)
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

func listSubaccountRoleCollectionQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_role_collection" "%s" {
  provider = "%s"
  config {
    subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountRoleCollectionQueryConfigWithIncludeResource(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_role_collection" "%s" {
  provider = "%s"
  include_resource = true
  config {
    subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}
