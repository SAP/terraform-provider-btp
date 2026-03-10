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

func TestSubaccountEntitlementListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "77395f6a-a601-4c9e-8cd0-c1fcefc7f60f"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_entitlement")
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
					Config: hclProviderFor(user) + listSubaccountEntitlementQueryConfig(
						"entitlement_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_entitlement.entitlement_list", 100),

						querycheck.ExpectIdentity(
							"btp_subaccount_entitlement.entitlement_list",
							map[string]knownvalue.Check{
								"plan_name":     knownvalue.StringExact("oauth2"),
								"service_name":  knownvalue.StringExact("cias"),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountEntitlementQueryConfigWithIncludeResource(
						"entitlement_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_entitlement.entitlement_list", 100),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_entitlement.entitlement_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"plan_name":     knownvalue.StringExact("oauth2"),
								"service_name":  knownvalue.StringExact("cias"),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringExact("cias-oauth2"),
								},
								{
									Path:       tfjsonpath.New("category"),
									KnownValue: knownvalue.StringExact("ELASTIC_SERVICE"),
								},
								{
									Path:       tfjsonpath.New("plan_id"),
									KnownValue: knownvalue.StringExact("cias-oauth2"),
								},
								{
									Path:       tfjsonpath.New("plan_name"),
									KnownValue: knownvalue.StringExact("oauth2"),
								},
								{
									Path:       tfjsonpath.New("plan_unique_identifier"),
									KnownValue: knownvalue.StringExact("cias-oauth2"),
								},
								{
									Path:       tfjsonpath.New("service_name"),
									KnownValue: knownvalue.StringExact("cias"),
								},
								{
									Path:       tfjsonpath.New("subaccount_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountEntitlementListResource().(list.ListResourceWithConfigure)
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

func listSubaccountEntitlementQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_entitlement" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountEntitlementQueryConfigWithIncludeResource(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_entitlement" "%s" {
  provider = "%s"
  include_resource = true
   config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}
