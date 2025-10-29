package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountRole(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_role")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountRole("uut", "integration-test-acc-static", "Subaccount Viewer", "Subaccount_Viewer", "cis-local!b2"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_role.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "name", "Subaccount Viewer"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "role_template_name", "Subaccount_Viewer"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "app_id", "cis-local!b2"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "description", "Role for subaccount members with read-only authorizations for core commercialization operations, such as viewing subaccount entitlements, details of environment instances, and job results."),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "read_only", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "scopes.#", "6"),
					),
				},
			},
		})
	})

	t.Run("happy path - with attributes", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_role_with_attributes")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountRole("uut", "integration-test-acc-static", "Application_Frontend_Developer", "Application_Frontend_Developer", "eu12-appfront!b390135"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_role.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "name", "Application_Frontend_Developer"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "role_template_name", "Application_Frontend_Developer"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "app_id", "eu12-appfront!b390135"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "description", "Developer access to Application Frontend"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "read_only", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "attribute_list.#", "1"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "attribute_list.0.attribute_name", "namespace"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "attribute_list.0.attribute_value_origin", "static"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "attribute_list.0.attribute_values.#", "1"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "attribute_list.0.attribute_values.0", ""),
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "attribute_list.0.value_required", "false"),
					),
				},
			},
		})
	})
	t.Run("error path - subaccount_id, name, role_template_name and app_id are mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_subaccount_role" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "(subaccount_id|name|role_template_name|app_id)" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - name must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountRoleBySubaccountId("uut", "00000000-0000-0000-0000-000000000000", "", "b", "c"),
					ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1, got: 0`),
				},
			},
		})
	})
	t.Run("error path - role_template_name must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountRoleBySubaccountId("uut", "00000000-0000-0000-0000-000000000000", "a", "", "c"),
					ExpectError: regexp.MustCompile(`Attribute role_template_name string length must be at least 1, got: 0`),
				},
			},
		})
	})
	t.Run("error path - app_id must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountRoleBySubaccountId("uut", "00000000-0000-0000-0000-000000000000", "a", "b", ""),
					ExpectError: regexp.MustCompile(`Attribute app_id string length must be at least 1, got: 0`),
				},
			},
		})
	})
	t.Run("error path - cli server returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/login/") {
				_, _ = fmt.Fprintf(w, "{}")
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(srv.Client()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceSubaccountRoleBySubaccountId("uut", "00000000-0000-0000-0000-000000000000", "Subaccount Viewer", "Subaccount_Viewer", "cis-local!b2"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountRoleBySubaccountId(resourceName string, subaccountId string, name string, roleTemplateName string, appId string) string {
	template := `
data "btp_subaccount_role" "%s" {
    subaccount_id       = "%s"
    name                = "%s"
    role_template_name  = "%s"
    app_id              = "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId, name, roleTemplateName, appId)
}

func hclDatasourceSubaccountRole(resourceName string, subaccountName string, name string, roleTemplateName string, appId string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_role" "%s" {
    subaccount_id       = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
    name                = "%s"
    role_template_name  = "%s"
    app_id              = "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountName, name, roleTemplateName, appId)
}
