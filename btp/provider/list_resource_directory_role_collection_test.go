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

func TestDirectoryRoleCollectionListResource(t *testing.T) {
	t.Parallel()

	directoryID := "0f7a9b71-0b19-4b6c-b20b-ab2e5445bdc2"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_directory_role_collection")
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
					Config: hclProviderFor(user) + listDirectoryRoleCollectionQueryConfig(
						"directory_role_collection_list",
						"btp",
						directoryID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_directory_role_collection.directory_role_collection_list", 2),

						querycheck.ExpectIdentity(
							"btp_directory_role_collection.directory_role_collection_list",
							map[string]knownvalue.Check{
								"name":         knownvalue.StringExact("Directory Viewer"),
								"directory_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listDirectoryRoleCollectionQueryConfigWithIncludeResource(
						"directory_role_collection_list",
						"btp",
						directoryID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_directory_role_collection.directory_role_collection_list", 2),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_directory_role_collection.directory_role_collection_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"name":         knownvalue.StringExact("Directory Viewer"),
								"directory_id": knownvalue.StringRegexp(regexpValidUUID),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("description"),
									KnownValue: knownvalue.StringExact("Read-only access to the directory"),
								},
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("Directory Viewer"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringExact("0f7a9b71-0b19-4b6c-b20b-ab2e5445bdc2,Directory Viewer"),
								},
								{
									Path:       tfjsonpath.New("directory_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("roles"),
									KnownValue: knownvalue.ListSizeExact(3),
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

func listDirectoryRoleCollectionQueryConfig(label, providerName, directoryID string) string {
	return fmt.Sprintf(`
		list "btp_directory_role_collection" "%s" {
		provider = "%s"
		config {
		directory_id = "%s"
		}
		}`, label, providerName, directoryID)
}

func listDirectoryRoleCollectionQueryConfigWithIncludeResource(label, providerName, directoryID string) string {
	return fmt.Sprintf(`	
		list "btp_directory_role_collection" "%s" {
		provider = "%s"
		include_resource = true
		config {
		directory_id = "%s"
		}
		}`, label, providerName, directoryID)
}
