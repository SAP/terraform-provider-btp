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
					Config: hclProviderFor(user) + hclDatasourceSubaccountRole("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "Subaccount Viewer", "Subaccount_Viewer", "cis-local!b2"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_role.uut", "subaccount_id", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
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
	t.Run("error path - subaccount_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountRole("uut", "this-is-not-a-uuid", "a", "b", "c"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
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
					Config:      hclDatasourceSubaccountRole("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "", "b", "c"),
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
					Config:      hclDatasourceSubaccountRole("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "a", "", "c"),
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
					Config:      hclDatasourceSubaccountRole("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "a", "b", ""),
					ExpectError: regexp.MustCompile(`Attribute app_id string length must be at least 1, got: 0`),
				},
			},
		})
	})
	t.Run("error path - cli server returns error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/login/") {
				fmt.Fprintf(w, "{}")
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceSubaccountRole("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "Subaccount Viewer", "Subaccount_Viewer", "cis-local!b2"),
					ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountRole(resourceName string, subaccountId string, name string, roleTemplateName string, appId string) string {
	template := `
data "btp_subaccount_role" "%s" {
    subaccount_id       = "%s"
    name                = "%s"
    role_template_name  = "%s"
    app_id              = "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId, name, roleTemplateName, appId)
}
