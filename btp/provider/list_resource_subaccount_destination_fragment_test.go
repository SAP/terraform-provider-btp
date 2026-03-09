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

func TestSubaccountDestinationFragmentListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "ba268910-81e6-4ac1-9016-cae7ed196889"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_destination_fragment")
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
					Config: hclProviderFor(user) + listSubaccountDestinationFragmentQueryConfig(
						"destination_fragment_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_destination_fragment.destination_fragment_list", 1),

						querycheck.ExpectIdentity(
							"btp_subaccount_destination_fragment.destination_fragment_list",
							map[string]knownvalue.Check{
								"name":                knownvalue.StringExact("test_destination_fragment"),
								"subaccount_id":       knownvalue.StringExact(subaccountID),
								"service_instance_id": knownvalue.Null(),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountDestinationFragmentQueryConfigWithIncludeResource(
						"destination_fragment_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_destination_fragment.destination_fragment_list", 1),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_destination_fragment.destination_fragment_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"name":                knownvalue.StringExact("test_destination_fragment"),
								"subaccount_id":       knownvalue.StringExact(subaccountID),
								"service_instance_id": knownvalue.Null(),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("fragment_content"),
									KnownValue: knownvalue.NotNull(),
								},
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("test_destination_fragment"),
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
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_destination_fragment_bad_request")
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
					Config: hclProviderFor(user) + listSubaccountDestinationFragmentQueryConfig("bad_request", "btp", badRequestSubaccountID),

					ExpectError: regexp.MustCompile(`(?i)API Error Reading Resource Destination Fragment \(Subaccount\).*Detail:Failed to list destination fragments: received response with unexpected status: 400`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountDestinationFragmentListResource().(list.ListResourceWithConfigure)
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

func listSubaccountDestinationFragmentQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_destination_fragment" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountDestinationFragmentQueryConfigWithIncludeResource(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_destination_fragment" "%s" {
  provider = "%s"
  include_resource = true
   config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}
