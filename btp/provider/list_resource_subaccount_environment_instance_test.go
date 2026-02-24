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

func TestSubaccountEnvironmentInstanceListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "59cd458e-e66e-4b60-b6d8-8f219379f9a5"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_environment_instance")
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
					Config: hclProviderFor(user) + listSubaccountEnvironmentInstanceQueryConfig(
						"environment_instances_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_environment_instance.environment_instances_list", 1),

						querycheck.ExpectIdentity(
							"btp_subaccount_environment_instance.environment_instances_list",
							map[string]knownvalue.Check{
								"id":            knownvalue.StringExact("28EE0D05-966B-4218-8286-14D15B71B610"),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountEnvironmentInstanceQueryConfigWithIncludeResource(
						"environment_instances_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_environment_instance.environment_instances_list", 1),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_environment_instance.environment_instances_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"id":            knownvalue.StringExact("28EE0D05-966B-4218-8286-14D15B71B610"),
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("integration-test-services-4ie3yr1a_integration-test-services-static"),
								},
								{
									Path:       tfjsonpath.New("service_name"),
									KnownValue: knownvalue.StringExact("cloudfoundry"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("tenant_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("platform_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("service_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("subaccount_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("state"),
									KnownValue: knownvalue.StringExact("OK"),
								},
								{
									Path:       tfjsonpath.New("type"),
									KnownValue: knownvalue.StringExact("Provision"),
								},
								{
									Path:       tfjsonpath.New("operation"),
									KnownValue: knownvalue.StringExact("provision"),
								},
								{
									Path:       tfjsonpath.New("parameters"),
									KnownValue: knownvalue.NotNull(),
								},
								{
									Path:       tfjsonpath.New("labels"),
									KnownValue: knownvalue.NotNull(),
								},
								{
									Path:       tfjsonpath.New("environment_type"),
									KnownValue: knownvalue.StringExact("cloudfoundry"),
								},
								{
									Path:       tfjsonpath.New("landscape_label"),
									KnownValue: knownvalue.StringExact("cf-eu12"),
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
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_environment_instance_bad_request")
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
					Config: hclProviderFor(user) + listSubaccountEnvironmentInstanceQueryConfig("bad_request", "btp", badRequestSubaccountID),

					ExpectError: regexp.MustCompile(`API Error Reading Resource Environment Instance \(Subaccount\).*Detail:Failed to list environment instances: received response with unexpected status: 400`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountEnvironmentInstanceListResource().(list.ListResourceWithConfigure)
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

func listSubaccountEnvironmentInstanceQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_environment_instance" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountEnvironmentInstanceQueryConfigWithIncludeResource(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_environment_instance" "%s" {
  provider = "%s"
  include_resource = true
   config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}
