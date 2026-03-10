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

func TestSubaccountRoleListResource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_role")
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
					Config: hclProviderFor(user) + listSubaccountRoleQueryConfig("roles_list", "btp"),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role.roles_list", 44),

						querycheck.ExpectIdentity(
							"btp_subaccount_role.roles_list",
							map[string]knownvalue.Check{
								"name":               knownvalue.StringExact("Destination Administrator"),
								"role_template_name": knownvalue.StringExact("Destination_Administrator"),
								"app_id":             knownvalue.StringExact("destination-xsappname!b9"),
								"subaccount_id":      knownvalue.StringExact("ba268910-81e6-4ac1-9016-cae7ed196889"),
							},
						),
					},
				},
				{
					// List query with include_resource = true
					Query:  true,
					Config: hclProviderFor(user) + listSubaccountRoleQueryConfigWithIncludeResource("roles_list", "btp"),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role.roles_list", 44),

						querycheck.ExpectResourceKnownValues(
							"btp_subaccount_role.roles_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"name":               knownvalue.StringExact("Destination Administrator"),
								"role_template_name": knownvalue.StringExact("Destination_Administrator"),
								"app_id":             knownvalue.StringExact("destination-xsappname!b9"),
								"subaccount_id":      knownvalue.StringExact("ba268910-81e6-4ac1-9016-cae7ed196889"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("Destination Administrator"),
								},
								{
									Path:       tfjsonpath.New("role_template_name"),
									KnownValue: knownvalue.StringExact("Destination_Administrator"),
								},
								{
									Path:       tfjsonpath.New("app_id"),
									KnownValue: knownvalue.StringExact("destination-xsappname!b9"),
								},
								{
									Path:       tfjsonpath.New("subaccount_id"),
									KnownValue: knownvalue.StringExact("ba268910-81e6-4ac1-9016-cae7ed196889"),
								},
								{
									Path: tfjsonpath.New("id"),
									KnownValue: knownvalue.StringExact(
										"ba268910-81e6-4ac1-9016-cae7ed196889,Destination Administrator,Destination_Administrator,destination-xsappname!b9",
									),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewSubaccountRoleListResource().(list.ListResourceWithConfigure)

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

func listSubaccountRoleQueryConfig(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_subaccount_role" "%s" {
  provider = "%s"
  config {
    subaccount_id = "ba268910-81e6-4ac1-9016-cae7ed196889"
  }
}`, label, providerName)
}

func listSubaccountRoleQueryConfigWithIncludeResource(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_subaccount_role" "%s" {
  provider = "%s"
  include_resource = true
  config {
    subaccount_id = "ba268910-81e6-4ac1-9016-cae7ed196889"
  }
}`, label, providerName)
}
