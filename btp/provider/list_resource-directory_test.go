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

func TestDirectoryListResource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_directory")
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
					Config: hclProviderFor(user) + listDirectoryQueryConfig(
						"directory_list",
						"btp",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_directory.directory_list", 11),

						querycheck.ExpectIdentity(
							"btp_directory.directory_list",
							map[string]knownvalue.Check{
								"directory_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},

				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listDirectoryQueryConfigWithIncludeResource(
						"directory_list",
						"btp",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_directory.directory_list", 11),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_directory.directory_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"directory_id": knownvalue.StringExact("14870944-4832-4e76-83f7-d2913661cf6d"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("integration-test-dir-se-static"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("parent_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("subdomain"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("features"),
									KnownValue: knownvalue.NotNull(),
								},
								{
									Path:       tfjsonpath.New("labels"),
									KnownValue: knownvalue.NotNull(),
								},
								{
									Path:       tfjsonpath.New("description"),
									KnownValue: knownvalue.StringExact("Please don't modify. This is used for integration tests."),
								},
								{
									Path:       tfjsonpath.New("state"),
									KnownValue: knownvalue.StringExact("OK"),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewDirectoryListResource().(list.ListResourceWithConfigure)
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

func listDirectoryQueryConfig(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_subaccount" "%s" {
  provider = "%s"
}`, label, providerName)
}

func listDirectoryQueryConfigWithIncludeResource(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_subaccount" "%s" {
  provider = "%s"
  include_resource = true
}`, label, providerName)
}
