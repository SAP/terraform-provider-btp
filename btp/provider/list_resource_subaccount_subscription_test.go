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

func TestSubaccountSubscriptionListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "ba268910-81e6-4ac1-9016-cae7ed196889"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_subscription")
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
					Config: hclProviderFor(user) + listSubaccountSubscriptionQueryConfig(
						"subscription_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_subscription.subscription_list", 5),

						querycheck.ExpectIdentity(
							"btp_subaccount_subscription.subscription_list",
							map[string]knownvalue.Check{
								"app_name":      knownvalue.StringExact("feature-flags-dashboard"),
								"subaccount_id": knownvalue.StringExact(subaccountID),
								"plan_name":     knownvalue.StringExact("dashboard"),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountSubscriptionQueryConfigWithIncludeResource(
						"subscription_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_subscription.subscription_list", 5),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_subscription.subscription_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"app_name":      knownvalue.StringExact("feature-flags-dashboard"),
								"subaccount_id": knownvalue.StringExact(subaccountID),
								"plan_name":     knownvalue.StringExact("dashboard"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("app_id"),
									KnownValue: knownvalue.StringExact("feature-flags!b18"),
								},
								{
									Path:       tfjsonpath.New("app_name"),
									KnownValue: knownvalue.StringExact("feature-flags-dashboard"),
								},
								{
									Path:       tfjsonpath.New("commercial_app_name"),
									KnownValue: knownvalue.StringExact("feature-flags-dashboard"),
								},
								{
									Path:       tfjsonpath.New("category"),
									KnownValue: knownvalue.StringExact("Foundation / Cross Services"),
								},
								{
									Path:       tfjsonpath.New("description"),
									KnownValue: knownvalue.NotNull(),
								},
								{
									Path:       tfjsonpath.New("state"),
									KnownValue: knownvalue.StringExact("NOT_SUBSCRIBED"),
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
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_subscription_bad_request")
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
					Config: hclProviderFor(user) + listSubaccountSubscriptionQueryConfig("bad_request", "btp", badRequestSubaccountID),

					ExpectError: regexp.MustCompile(`API Error Reading Resource Subscription \(SubAccount\).*Detail:Failed to list subscriptions: received response with unexpected status: 400`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountSubscriptionListResource().(list.ListResourceWithConfigure)
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

func listSubaccountSubscriptionQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_subscription" "%s" {
  provider = "%s"
  config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountSubscriptionQueryConfigWithIncludeResource(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_subscription" "%s" {
  provider = "%s"
  include_resource = true
   config {
   subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}
