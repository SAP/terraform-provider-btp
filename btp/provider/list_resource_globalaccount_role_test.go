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

func TestGlobalaccountRoleListResource(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/list_resource_globalaccount_role")
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
					Config: hclProviderFor(user) + listGlobalAccountRoleQueryConfig("roles_list", "btp"),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_globalaccount_role.roles_list", 14),

						querycheck.ExpectIdentity(
							"btp_globalaccount_role.roles_list",
							map[string]knownvalue.Check{
								"name":               knownvalue.StringExact("Auditlog_Auditor"),
								"role_template_name": knownvalue.StringExact("Auditlog_Auditor"),
								"app_id":             knownvalue.StringExact("auditlog-management!b3034"),
							},
						),
					},
				},
				{
					// List Query with include_resource = true
					Query:  true,
					Config: hclProviderFor(user) + listGlobalAccountRoleQueryConfigWithIncludeResource("roles_list", "btp"),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_globalaccount_role.roles_list", 14),

						// Verify full resource data is populated
						querycheck.ExpectResourceKnownValues(
							"btp_globalaccount_role.roles_list",
							queryfilter.ByResourceIdentity(map[string]knownvalue.Check{
								"name":               knownvalue.StringExact("Auditlog_Auditor"),
								"role_template_name": knownvalue.StringExact("Auditlog_Auditor"),
								"app_id":             knownvalue.StringExact("auditlog-management!b3034"),
							}),
							[]querycheck.KnownValueCheck{
								{
									Path:       tfjsonpath.New("description"),
									KnownValue: knownvalue.StringExact("Read access to audit logs"),
								},
								{
									Path:       tfjsonpath.New("role_template_name"),
									KnownValue: knownvalue.StringExact("Auditlog_Auditor"),
								},
								{
									Path:       tfjsonpath.New("app_id"),
									KnownValue: knownvalue.StringExact("auditlog-management!b3034"),
								},
								{
									Path:       tfjsonpath.New("name"),
									KnownValue: knownvalue.StringExact("Auditlog_Auditor"),
								},
								{
									Path:       tfjsonpath.New("read_only"),
									KnownValue: knownvalue.Bool(true),
								},
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		r := NewGlobalaccountRoleListResource().(list.ListResourceWithConfigure)
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

func listGlobalAccountRoleQueryConfig(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_globalaccount_role" "%s" {
  provider = "%s"
}`, label, providerName)
}

func listGlobalAccountRoleQueryConfigWithIncludeResource(label, providerName string) string {
	return fmt.Sprintf(`
list "btp_globalaccount_role" "%s" {
  provider = "%s"
  include_resource = true
}`, label, providerName)
}
