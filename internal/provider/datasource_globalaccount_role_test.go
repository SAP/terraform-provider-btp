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

func TestDataSourceGlobalaccountRole(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_role")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountRole("uut", "Global Account Viewer", "GlobalAccount_Viewer", "cis-central!b13"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_role.uut", "name", "Global Account Viewer"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role.uut", "role_template_name", "GlobalAccount_Viewer"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role.uut", "app_id", "cis-central!b13"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role.uut", "description", "Role for global account members with read-only authorizations for core commercialization operations, such as viewing global accounts, subaccounts, entitlements, and regions."),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role.uut", "read_only", "true"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role.uut", "scopes.#", "9"),
					),
				},
			},
		})
	})
	t.Run("error path - name, role_template_name and app_id are mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_globalaccount_role" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "(name|role_template_name|app_id)" is required, but no definition was found.`),
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
					Config:      hclDatasourceGlobalaccountRole("uut", "", "b", "c"),
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
					Config:      hclDatasourceGlobalaccountRole("uut", "a", "", "c"),
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
					Config:      hclDatasourceGlobalaccountRole("uut", "a", "b", ""),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceGlobalaccountRole("uut", "Global Account Viewer", "GlobalAccount_Viewer", "cis-local!b13"),
					ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountRole(resourceName string, name string, roleTemplateName string, appId string) string {
	template := `
data "btp_globalaccount_role" "%s" {
    name               = "%s"
    role_template_name = "%s"
    app_id             = "%s"
}`

	return fmt.Sprintf(template, resourceName, name, roleTemplateName, appId)
}
