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

func TestDirectoryRoleListResource(t *testing.T) {
	t.Parallel()

	directoryID := "14870944-4832-4e76-83f7-d2913661cf6d"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_directory_role")
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
					Config: hclProviderFor(user) + listDirectoryRoleQueryConfig(
						"directory_role_list",
						"btp",
						directoryID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_directory_role.directory_role_list", 8),

						querycheck.ExpectIdentity(
							"btp_directory_role.directory_role_list",
							map[string]knownvalue.Check{
								"directory_id":       knownvalue.StringExact(directoryID),
								"name":               knownvalue.StringExact("Global Account Admin"),
								"role_template_name": knownvalue.StringExact("GlobalAccount_Admin"),
								"app_id":             knownvalue.StringExact("cis-central!b13"),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listDirectoryRoleQueryConfigWithIncludeResource(
						"directory_role_list",
						"btp",
						directoryID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_directory_role.directory_role_list", 8),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_directory_role.directory_role_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"directory_id":       knownvalue.StringExact(directoryID),
								"name":               knownvalue.StringExact("Global Account Admin"),
								"role_template_name": knownvalue.StringExact("GlobalAccount_Admin"),
								"app_id":             knownvalue.StringExact("cis-central!b13"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("description"),
									KnownValue: knownvalue.NotNull(),
								},
								{
									Path:       tfjsonpath.New("role_template_name"),
									KnownValue: knownvalue.StringExact("GlobalAccount_Admin"),
								},
								{
									Path:       tfjsonpath.New("app_id"),
									KnownValue: knownvalue.StringExact("cis-central!b13"),
								},
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("Global Account Admin"),
								},
								{
									Path:       tfjsonpath.New("read_only"),
									KnownValue: knownvalue.Bool(true),
								},
								{
									Path:       tfjsonpath.New("scopes"),
									KnownValue: knownvalue.NotNull(),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - bad request", func(t *testing.T) {
		badRequestDirectoryID := ""
		rec, user := setupVCR(t, "fixtures/list_resource_directory_role_bad_request")
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
					Config: hclProviderFor(user) + listDirectoryRoleQueryConfig("bad_request", "btp", badRequestDirectoryID),

					ExpectError: regexp.MustCompile(`API Error Reading Resource Directory Roles .*Detail:Failed to list directory roles: received response with unexpected status: 400`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewDirectoryRoleListResource().(list.ListResourceWithConfigure)
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

func listDirectoryRoleQueryConfig(label, providerName, directoryID string) string {
	return fmt.Sprintf(`
list "btp_directory_role" "%s" {
  provider = "%s"
  config {
   directory_id = "%s"
  }
}`, label, providerName, directoryID)
}

func listDirectoryRoleQueryConfigWithIncludeResource(label, providerName, directoryID string) string {
	return fmt.Sprintf(`
list "btp_directory_role" "%s" {
  provider = "%s"
  include_resource = true
   config {
   directory_id = "%s"
  }
}`, label, providerName, directoryID)
}
