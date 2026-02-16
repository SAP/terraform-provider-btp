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

func TestGlobalaccountResourceProviderListResource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_globalaccount_resource_provider")
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
					Config: hclProviderFor(user) + listGlobalAccountResourceProviderQueryConfig("resource_provider_list", "btp"),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_globalaccount_resource_provider.resource_provider_list", 1),

						querycheck.ExpectIdentity(
							"btp_globalaccount_resource_provider.resource_provider_list",
							map[string]knownvalue.Check{
								"provider_type":  knownvalue.StringExact("AWS"),
								"technical_name": knownvalue.StringExact("tf_test_resource_provider"),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query:  true,
					Config: hclProviderFor(user) + listGlobalAccountResourceProviderQueryConfigWithIncludeResource("resource_provider_list", "btp"),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_globalaccount_resource_provider.resource_provider_list", 1),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_globalaccount_resource_provider.resource_provider_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"provider_type":  knownvalue.StringExact("AWS"),
								"technical_name": knownvalue.StringExact("tf_test_resource_provider"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("provider_type"),
									KnownValue: knownvalue.StringExact("AWS"),
								},
								{
									Path:       tfjsonpath.New("description"),
									KnownValue: knownvalue.StringExact("Description of the resource provider"),
								},
								{
									Path:       tfjsonpath.New("technical_name"),
									KnownValue: knownvalue.StringExact("tf_test_resource_provider"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringExact("tf_test_resource_provider"),
								},
								{
									Path:       tfjsonpath.New("display_name"),
									KnownValue: knownvalue.StringExact("Test AWS Resource Provider"),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewGlobalaccountResourceProviderListResource().(list.ListResourceWithConfigure)
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

func listGlobalAccountResourceProviderQueryConfig(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_globalaccount_resource_provider" "%s" {
  provider = "%s"
}`, label, providerName)
}

func listGlobalAccountResourceProviderQueryConfigWithIncludeResource(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_globalaccount_resource_provider" "%s" {
  provider = "%s"
  include_resource = true
}`, label, providerName)
}
