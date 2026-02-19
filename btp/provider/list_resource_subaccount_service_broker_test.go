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

func TestSubaccountServiceBrokerListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "59cd458e-e66e-4b60-b6d8-8f219379f9a5"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_service_broker")
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
					Config: hclProviderFor(user) + listSubaccountServiceBrokerQueryConfig(
						"service_broker_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_service_broker.service_broker_list", 1),

						querycheck.ExpectIdentity(
							"btp_subaccount_service_broker.service_broker_list",
							map[string]knownvalue.Check{
								"id":            knownvalue.StringExact("d88f8a5b-1a7f-4aca-8279-fa0b2eabec7d"),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountServiceBrokerQueryConfigWithLableFilter(
						"service_broker_list",
						"btp",
						subaccountID,
						"cred_revision eq '0'",
						"ready eq 'true'",
					),
					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_service_broker.service_broker_list", 1),

						querycheck.ExpectIdentity(
							"btp_subaccount_service_broker.service_broker_list",
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
					Config: hclProviderFor(user) + listSubaccountServiceBrokerQueryConfigWithIncludeResource(
						"service_broker_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_service_broker.service_broker_list", 1),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_service_broker.service_broker_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"id":            knownvalue.StringExact("d88f8a5b-1a7f-4aca-8279-fa0b2eabec7d"),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("integration-test-static-service-broker-59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("subaccount_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("url"),
									KnownValue: knownvalue.StringExact("https://integration-test-static-service-broker.cfapps.eu12.hana.ondemand.com"),
								},
								{
									Path:       tfjsonpath.New("created_date"),
									KnownValue: knownvalue.NotNull(),
								},
								{
									Path:       tfjsonpath.New("ready"),
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
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_service_broker_bad_request")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Query:  true,
					Config: hclProviderFor(user) + listSubaccountServiceBrokerQueryConfig("bad_request", "btp", badRequestSubaccountID),

					ExpectError: regexp.MustCompile(`API Error Reading Resource Service Broker \(Subaccount\).*Detail:Failed to list service brokers: received response with unexpected status: 400`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountServiceBrokerListResource().(list.ListResourceWithConfigure)
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

func listSubaccountServiceBrokerQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_service_broker" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountServiceBrokerQueryConfigWithLableFilter(label, providerName, subaccountID, lableFilter, fieldFilter string) string {
	return fmt.Sprintf(`
list "btp_subaccount_service_broker" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
   fields_filter = "%s"
   labels_filter = "%s"

  }
}`, label, providerName, subaccountID, fieldFilter, lableFilter)
}

func listSubaccountServiceBrokerQueryConfigWithIncludeResource(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_service_broker" "%s" {
  provider = "%s"
  include_resource = true
   config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}
