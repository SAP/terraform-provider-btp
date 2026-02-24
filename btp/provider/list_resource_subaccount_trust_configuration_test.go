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

func TestSubaccountTrustConfigurationListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "59cd458e-e66e-4b60-b6d8-8f219379f9a5"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_trust_configuration")
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
					Config: hclProviderFor(user) + listSubaccountTrustConfigurationQueryConfig(
						"trust_configuration_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_trust_configuration.trust_configuration_list", 3),

						querycheck.ExpectIdentity(
							"btp_subaccount_trust_configuration.trust_configuration_list",
							map[string]knownvalue.Check{
								"origin":        knownvalue.StringExact("sap.default"),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountTrustConfigurationQueryConfigWithIncludeResource(
						"trust_configuration_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_trust_configuration.trust_configuration_list", 3),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_trust_configuration.trust_configuration_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"origin":        knownvalue.StringExact("sap.default"),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("status"),
									KnownValue: knownvalue.StringExact("active"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringExact("sap.default"),
								},
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("sap.default"),
								},
								{
									Path:       tfjsonpath.New("origin"),
									KnownValue: knownvalue.StringExact("sap.default"),
								},
								{
									Path:       tfjsonpath.New("protocol"),
									KnownValue: knownvalue.StringExact("OpenID Connect"),
								},
								{
									Path:       tfjsonpath.New("type"),
									KnownValue: knownvalue.StringExact("Application"),
								},
								{
									Path:       tfjsonpath.New("link_text"),
									KnownValue: knownvalue.StringExact("Default Identity Provider"),
								},
								{
									Path:       tfjsonpath.New("auto_create_shadow_users"),
									KnownValue: knownvalue.Bool(true),
								},
								{
									Path:       tfjsonpath.New("read_only"),
									KnownValue: knownvalue.Bool(false),
								},
								{
									Path:       tfjsonpath.New("available_for_user_logon"),
									KnownValue: knownvalue.Bool(true),
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
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_trust_configuration_bad_request")
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
					Config: hclProviderFor(user) + listSubaccountTrustConfigurationQueryConfig("bad_request", "btp", badRequestSubaccountID),

					ExpectError: regexp.MustCompile(`API Error Reading Resource Trust Configuration \(Subaccount\).*Detail:Failed to list trust configurations: received response with unexpected status: 400`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountTrustConfigurationListResource().(list.ListResourceWithConfigure)
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

func listSubaccountTrustConfigurationQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_trust_configuration" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountTrustConfigurationQueryConfigWithIncludeResource(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_trust_configuration" "%s" {
  provider = "%s"
  include_resource = true
   config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}
