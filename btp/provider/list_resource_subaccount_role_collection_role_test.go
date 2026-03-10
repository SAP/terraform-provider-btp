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

func TestSubaccountRoleCollectionRoleListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "ba268910-81e6-4ac1-9016-cae7ed196889"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_role_collection_role")
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
					Config: hclProviderFor(user) + listSubaccountRoleCollectionRoleQueryConfig(
						"role_collection_role_list",
						"btp",
						subaccountID,
						"Subaccount Viewer",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role_collection_role.role_collection_role_list", 8),

						querycheck.ExpectIdentity(
							"btp_subaccount_role_collection_role.role_collection_role_list",
							map[string]knownvalue.Check{
								"name":                 knownvalue.StringExact("Subaccount Viewer"),
								"subaccount_id":        knownvalue.StringExact(subaccountID),
								"role_name":            knownvalue.StringExact("Subaccount Viewer"),
								"role_template_name":   knownvalue.StringExact("Subaccount_Viewer"),
								"role_template_app_id": knownvalue.StringExact("cis-local!b2"),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountRoleCollectionRoleQueryConfigWithIncludeResource(
						"role_collection_role_list",
						"btp",
						subaccountID,
						"Subaccount Viewer",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role_collection_role.role_collection_role_list", 8),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_role_collection_role.role_collection_role_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"name":                 knownvalue.StringExact("Subaccount Viewer"),
								"subaccount_id":        knownvalue.StringExact(subaccountID),
								"role_name":            knownvalue.StringExact("Subaccount Viewer"),
								"role_template_name":   knownvalue.StringExact("Subaccount_Viewer"),
								"role_template_app_id": knownvalue.StringExact("cis-local!b2"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringExact("ba268910-81e6-4ac1-9016-cae7ed196889,Subaccount Viewer,Subaccount Viewer,cis-local!b2,Subaccount_Viewer"),
								},
								{
									Path:       tfjsonpath.New("subaccount_id"),
									KnownValue: knownvalue.StringExact(subaccountID),
								},
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("Subaccount Viewer"),
								},
								{
									Path:       tfjsonpath.New("role_name"),
									KnownValue: knownvalue.StringExact("Subaccount Viewer"),
								},
								{
									Path:       tfjsonpath.New("role_template_app_id"),
									KnownValue: knownvalue.StringExact("cis-local!b2"),
								},
								{
									Path:       tfjsonpath.New("role_template_name"),
									KnownValue: knownvalue.StringExact("Subaccount_Viewer"),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - bad request", func(t *testing.T) {
		badRequestSubaccountID := ""
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_role_collection_role_bad_request")
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
					Config: hclProviderFor(user) + listSubaccountRoleCollectionRoleQueryConfig("bad_request", "btp", badRequestSubaccountID, ""),

					ExpectError: regexp.MustCompile(`(?i)API Error Reading Resource Role Collection Role \(Subaccount\).*Detail:Failed to list role collection roles: received response with unexpected status: 400`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountRoleCollectionRoleListResource().(list.ListResourceWithConfigure)
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

func listSubaccountRoleCollectionRoleQueryConfig(label, providerName, subaccountID, name string) string {
	return fmt.Sprintf(`
list "btp_subaccount_role_collection_role" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
   name          = "%s"
  }
}`, label, providerName, subaccountID, name)
}

func listSubaccountRoleCollectionRoleQueryConfigWithIncludeResource(label, providerName, subaccountID, name string) string {
	return fmt.Sprintf(`
list "btp_subaccount_role_collection_role" "%s" {
  provider = "%s"
  include_resource = true
   config {
   subaccount_id = "%s"
   name          = "%s"
  }
}`, label, providerName, subaccountID, name)
}
