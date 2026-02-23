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

func TestSubaccountListResource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount")
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
					Config: hclProviderFor(user) + listSubaccountQueryConfig(
						"subaccount_list",
						"btp",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount.subaccount_list", 11),

						querycheck.ExpectIdentity(
							"btp_subaccount.subaccount_list",
							map[string]knownvalue.Check{
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountQueryConfigWithFilter(
						"subaccount_list",
						"btp",
						"eu12",
					),
					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount.subaccount_list", 10),

						querycheck.ExpectIdentity(
							"btp_subaccount.subaccount_list",
							map[string]knownvalue.Check{
								"subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountQueryConfigWithIncludeResource(
						"subaccount_list",
						"btp",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount.subaccount_list", 11),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount.subaccount_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"subaccount_id": knownvalue.StringExact("4e981c0f-de50-4442-a26e-54798120f141"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("integration-test-acc-entitlements-stacked"),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("parent_id"),
									KnownValue: knownvalue.StringRegexp(regexpValidUUID),
								},
								{
									Path:       tfjsonpath.New("subdomain"),
									KnownValue: knownvalue.StringExact("integration-test-acc-entitlements-stacked-gddtpz5i"),
								},
								{
									Path:       tfjsonpath.New("region"),
									KnownValue: knownvalue.StringExact("eu12"),
								},
								{
									Path:       tfjsonpath.New("usage"),
									KnownValue: knownvalue.StringExact("NOT_USED_FOR_PRODUCTION"),
								},
								{
									Path:       tfjsonpath.New("state"),
									KnownValue: knownvalue.StringExact("OK"),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - bad request", func(t *testing.T) {
		badRequestLablesFilter := "a"
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_bad_request")
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
					Config: hclProviderFor(user) + listSubaccountQueryConfigWithLableFilter("bad_request", "btp", badRequestLablesFilter),

					ExpectError: regexp.MustCompile(`API Error Reading Resource Subaccount.*Detail:Failed to list subaccounts:.*\[Error: .*/400\]`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountListResource().(list.ListResourceWithConfigure)
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

func listSubaccountQueryConfig(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_subaccount" "%s" {
  provider = "%s"
}`, label, providerName)
}

func listSubaccountQueryConfigWithFilter(label, providerName, region string) string {
	return fmt.Sprintf(`
list "btp_subaccount" "%s" {
  provider = "%s"
  config {
   region = "%s"
  }
}`, label, providerName, region)
}

func listSubaccountQueryConfigWithLableFilter(label, providerName, labelsFilter string) string {
	return fmt.Sprintf(`
list "btp_subaccount" "%s" {
  provider = "%s"
  config {
   labels_filter = "%s"
  }
}`, label, providerName, labelsFilter)
}

func listSubaccountQueryConfigWithIncludeResource(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_subaccount" "%s" {
  provider = "%s"
  include_resource = true
}`, label, providerName)
}
