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

func TestDirectoryEntitlementListResource(t *testing.T) {
	t.Parallel()

	directoryID := "0f7a9b71-0b19-4b6c-b20b-ab2e5445bdc2"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_directory_entitlement")
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
					Config: hclProviderFor(user) + listDirectoryEntitlementQueryConfig("entitlement_list", "btp", directoryID),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_directory_entitlement.entitlement_list", 34),

						querycheck.ExpectIdentity(
							"btp_directory_entitlement.entitlement_list",
							map[string]knownvalue.Check{
								"directory_id": knownvalue.StringRegexp(regexpValidUUID),
								"service_name": knownvalue.StringExact("auditlog-management"),
								"plan_name":    knownvalue.StringExact("default"),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query:  true,
					Config: hclProviderFor(user) + listDirectoryEntitlementQueryConfigWithIncludeResource("entitlement_list", "btp", directoryID),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_directory_entitlement.entitlement_list", 34),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_directory_entitlement.entitlement_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"directory_id": knownvalue.StringRegexp(regexpValidUUID),
								"service_name": knownvalue.StringExact("auditlog-management"),
								"plan_name":    knownvalue.StringExact("default"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("directory_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("service_name"),
									KnownValue: knownvalue.StringExact("auditlog-management"),
								},
								{
									Path:       tfjsonpath.New("plan_name"),
									KnownValue: knownvalue.StringExact("default"),
								},
								{
									Path:       tfjsonpath.New("plan_unique_identifier"),
									KnownValue: knownvalue.StringExact("auditlog-management-default"),
								},
								{
									Path:       tfjsonpath.New("auto_assign"),
									KnownValue: knownvalue.Bool(false),
								},
								{
									Path:       tfjsonpath.New("auto_distribute_amount"),
									KnownValue: knownvalue.Int64Exact(0),
								},
								{
									Path:       tfjsonpath.New("category"),
									KnownValue: knownvalue.StringExact("ELASTIC_SERVICE"),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - directory not found", func(t *testing.T) {
		invalidDirectoryID := "00000000-0000-0000-0000-000000000000"
		rec, user := setupVCR(t, "fixtures/list_resource_directory_entitlement_not_found")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Query:       true,
					Config:      hclProviderFor(user) + listDirectoryEntitlementQueryConfig("not_found", "btp", invalidDirectoryID),
					ExpectError: regexp.MustCompile(`Service plan assignment not allowed.*not child of parent`),
				},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewDirectoryEntitlementListResource().(list.ListResourceWithConfigure)
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

func listDirectoryEntitlementQueryConfig(label, providerName, directoryID string) string {
	return fmt.Sprintf(`
list "btp_directory_entitlement" "%s" {
  provider = "%s"
  config {
   directory_id = "%s"
  }
}`, label, providerName, directoryID)
}

func listDirectoryEntitlementQueryConfigWithIncludeResource(label, providerName, directoryID string) string {
	return fmt.Sprintf(`
list "btp_directory_entitlement" "%s" {
	provider = "%s"
	include_resource = true
	config {
	  directory_id = "%s"
	}
}`, label, providerName, directoryID)
}
