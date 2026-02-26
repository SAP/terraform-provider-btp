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

func TestSubaccountSecuritySettingsListResource(t *testing.T) {
	t.Parallel()

	subAccountID := "59cd458e-e66e-4b60-b6d8-8f219379f9a5"

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_security_settings")
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
					Config: hclProviderFor(user) + listSubaccountSecuritySettingsQueryConfig(
						"security_settings_list",
						"btp",
						subAccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_security_settings.security_settings_list", 1),

						querycheck.ExpectIdentity(
							"btp_subaccount_security_settings.security_settings_list",
							map[string]knownvalue.Check{
								"subaccount_id": knownvalue.StringExact(subAccountID),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query: true,
					Config: hclProviderFor(user) + listSubaccountSecuritySettingsQueryConfigWithIncludeResource(
						"security_settings_list",
						"btp",
						subAccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_security_settings.security_settings_list", 1),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_security_settings.security_settings_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"subaccount_id": knownvalue.StringExact(subAccountID),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("access_token_validity"),
									KnownValue: knownvalue.Int64Exact(-1),
								},
								{
									Path:       tfjsonpath.New("refresh_token_validity"),
									KnownValue: knownvalue.Int64Exact(-1),
								},
								{
									Path:       tfjsonpath.New("id"),
									KnownValue: knownvalue.StringExact(subAccountID),
								},
								{
									Path:       tfjsonpath.New("default_identity_provider"),
									KnownValue: knownvalue.StringExact("sap.default"),
								},
								{
									Path:       tfjsonpath.New("treat_users_with_same_email_as_same_user"),
									KnownValue: knownvalue.Bool(false),
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
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_security_settings_bad_request")
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
					Config: hclProviderFor(user) + listSubaccountSecuritySettingsQueryConfig("bad_request", "btp", badRequestSubaccountID),

					ExpectError: regexp.MustCompile(`API Error Reading Resource Security Settings \(Subaccount\).*Detail:Failed to list security settings: received response with unexpected status: 400`)},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountSecuritySettingsListResource().(list.ListResourceWithConfigure)
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

func listSubaccountSecuritySettingsQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_security_settings" "%s" {
  provider = "%s"
  config {
  subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountSecuritySettingsQueryConfigWithIncludeResource(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_security_settings" "%s" {
  provider = "%s"
  include_resource = true
  config {
  subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}
