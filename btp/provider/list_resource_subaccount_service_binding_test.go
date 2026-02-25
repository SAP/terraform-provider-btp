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

func TestSubaccountServiceBindingListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "59cd458e-e66e-4b60-b6d8-8f219379f9a5"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_service_binding")
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
					Config: hclProviderFor(user) + listSubaccountServiceBindingQueryConfig(
						"service_binding_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_service_binding.service_binding_list", 3),

						querycheck.ExpectIdentity(
							"btp_subaccount_service_binding.service_binding_list",
							map[string]knownvalue.Check{
								"id":            knownvalue.StringExact("7b6fd614-167f-49ce-bb23-3a19a63a5359"),
								"subaccount_id": knownvalue.StringExact(subaccountID),
							},
						),
					},
				},
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountServiceBindingQueryConfigWithFilter(
						"service_binding_list",
						"btp",
						subaccountID,
						"ready eq 'true'",
					),
					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_service_binding.service_binding_list", 3),

						querycheck.ExpectIdentity(
							"btp_subaccount_service_binding.service_binding_list",
							map[string]knownvalue.Check{
								"id":            knownvalue.StringExact("7b6fd614-167f-49ce-bb23-3a19a63a5359"),
								"subaccount_id": knownvalue.StringExact(subaccountID),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountServiceBindingQueryConfigWithIncludeResource(
						"service_binding_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_service_binding.service_binding_list", 3),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_service_binding.service_binding_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"id":            knownvalue.StringExact("7b6fd614-167f-49ce-bb23-3a19a63a5359"),
								"subaccount_id": knownvalue.StringExact(subaccountID),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("service_instance_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("test-service-binding-malware"),
								},
								{
									Path:       tfjsonpath.New("credentials"),
									KnownValue: knownvalue.Null(),
								},

								{
									Path:       tfjsonpath.New("ready"),
									KnownValue: knownvalue.Bool(true),
								},
								{
									Path:       tfjsonpath.New("context"),
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
		badRequestSubaccountID := ""
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_service_binding_bad_request")
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
					Config: hclProviderFor(user) + listSubaccountServiceBindingQueryConfig("bad_request", "btp", badRequestSubaccountID),

					ExpectError: regexp.MustCompile(`API Error Reading Resource Service Binding \(Subaccount\).*Detail:Failed to list service bindings: received response with unexpected status: 400`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountServiceBindingListResource().(list.ListResourceWithConfigure)
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

func listSubaccountServiceBindingQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_service_binding" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountServiceBindingQueryConfigWithFilter(label, providerName, subaccountID, fieldFilter string) string {
	return fmt.Sprintf(`
list "btp_subaccount_service_binding" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
   fields_filter = "%s"
  }
}`, label, providerName, subaccountID, fieldFilter)
}

func listSubaccountServiceBindingQueryConfigWithIncludeResource(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_service_binding" "%s" {
  provider = "%s"
  include_resource = true
   config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}
