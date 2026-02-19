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

func TestSubaccountServiceInstanceListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "59cd458e-e66e-4b60-b6d8-8f219379f9a5"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_service_instance")
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
					Config: hclProviderFor(user) + listSubaccountServiceInstanceQueryConfig(
						"service_instances_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_service_instance.service_instances_list", 3),

						querycheck.ExpectIdentity(
							"btp_subaccount_service_instance.service_instances_list",
							map[string]knownvalue.Check{
								"id":            knownvalue.StringExact("df532d07-57a7-415e-a261-23a398ef068a"),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountServiceInstanceQueryConfigWithLableFilter(
						"service_instances_list",
						"btp",
						subaccountID,
						"org eq 'testvalue'",
						"ready eq 'true'",
					),
					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_service_instance.service_instances_list", 1),

						querycheck.ExpectIdentity(
							"btp_subaccount_service_instance.service_instances_list",
							map[string]knownvalue.Check{
								"id":            knownvalue.StringRegexp(regexpValidUUID),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountServiceInstanceQueryConfigWithIncludeResource(
						"service_instances_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_service_instance.service_instances_list", 3),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_service_instance.service_instances_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"id":            knownvalue.StringExact("df532d07-57a7-415e-a261-23a398ef068a"),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("tf-testacc-alertnotification-instance"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("platform_id"),
									KnownValue: knownvalue.StringExact("service-manager"),
								},
								{
									Path:       tfjsonpath.New("serviceplan_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("subaccount_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("usable"),
									KnownValue: knownvalue.Bool(true),
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
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_service_instance_bad_request")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Query:  true,
					Config: hclProviderFor(user) + listSubaccountServiceInstanceQueryConfig("bad_request", "btp", badRequestSubaccountID),

					ExpectError: regexp.MustCompile(`API Error Reading Resource Service Instance \(Subaccount\).*Detail:Failed to list service instances: received response with unexpected status: 400`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountServiceInstanceListResource().(list.ListResourceWithConfigure)
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

func listSubaccountServiceInstanceQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_service_instance" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountServiceInstanceQueryConfigWithLableFilter(label, providerName, subaccountID, lableFilter, fieldFilter string) string {
	return fmt.Sprintf(`
list "btp_subaccount_service_instance" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
   fields_filter = "%s"
   labels_filter = "%s"

  }
}`, label, providerName, subaccountID, fieldFilter, lableFilter)
}

func listSubaccountServiceInstanceQueryConfigWithIncludeResource(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_service_instance" "%s" {
  provider = "%s"
  include_resource = true
   config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}
