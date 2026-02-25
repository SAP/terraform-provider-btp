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

func TestGlobalaccountTrustConfigurationListResource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_globalaccount_trust_configuration")
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
					Config: hclProviderFor(user) + listGlobalaccountTrustConfigurationQueryConfig(
						"trust_configuration_list",
						"btp",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_globalaccount_trust_configuration.trust_configuration_list", 4),

						querycheck.ExpectIdentity(
							"btp_globalaccount_trust_configuration.trust_configuration_list",
							map[string]knownvalue.Check{
								"origin": knownvalue.StringExact("terraformint-platform"),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listGlobalaccountTrustConfigurationQueryConfigWithIncludeResource(
						"trust_configuration_list",
						"btp",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_globalaccount_trust_configuration.trust_configuration_list", 4),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_globalaccount_trust_configuration.trust_configuration_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"origin": knownvalue.StringExact("terraformint-platform"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("status"),
									KnownValue: knownvalue.StringExact("active"),
								},
								{
									Path:       tfjsonpath.New("description"),
									KnownValue: knownvalue.StringExact("Custom Platform Identity Provider"),
								},
								{
									Path:       tfjsonpath.New("domain"),
									KnownValue: knownvalue.StringExact("terraformint.accounts400.ondemand.com"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringExact("terraformint-platform"),
								},
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("terraformint-platform"),
								},
								{
									Path:       tfjsonpath.New("origin"),
									KnownValue: knownvalue.StringExact("terraformint-platform"),
								},
								{
									Path:       tfjsonpath.New("protocol"),
									KnownValue: knownvalue.StringExact("OpenID Connect"),
								},
								{
									Path:       tfjsonpath.New("type"),
									KnownValue: knownvalue.StringExact("Platform"),
								},
								{
									Path:       tfjsonpath.New("identity_provider"),
									KnownValue: knownvalue.StringExact("terraformint.accounts400.ondemand.com"),
								},

								{
									Path:       tfjsonpath.New("read_only"),
									KnownValue: knownvalue.Bool(false),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewGlobalaccountTrustConfigurationListResource().(list.ListResourceWithConfigure)
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

func listGlobalaccountTrustConfigurationQueryConfig(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_globalaccount_trust_configuration" "%s" {
  provider = "%s"
}`, label, providerName)
}

func listGlobalaccountTrustConfigurationQueryConfigWithIncludeResource(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_globalaccount_trust_configuration" "%s" {
  provider = "%s"
  include_resource = true
}`, label, providerName)
}
