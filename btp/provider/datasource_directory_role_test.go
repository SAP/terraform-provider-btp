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

func TestDataSourceDirectoryRole(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_role")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDirectoryRole("uut", "integration-test-dir-se-static", "Directory Viewer", "Directory_Viewer", "cis-central!b13"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_directory_role.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_directory_role.uut", "name", "Directory Viewer"),
						resource.TestCheckResourceAttr("data.btp_directory_role.uut", "role_template_name", "Directory_Viewer"),
						resource.TestCheckResourceAttr("data.btp_directory_role.uut", "app_id", "cis-central!b13"),
						resource.TestCheckResourceAttr("data.btp_directory_role.uut", "description", "Role for directory members with read-only authorizations for core commercialization operations, such as viewing directories, subaccounts, entitlements, and regions."),
						resource.TestCheckResourceAttr("data.btp_directory_role.uut", "read_only", "true"),
						resource.TestCheckResourceAttr("data.btp_directory_role.uut", "scopes.#", "7"),
					),
				},
			},
		})
	})

	t.Run("error path - directory not security enabled", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_role.not_security_enabled")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceDirectoryRole("uut", "integration-test-dir-static", "Directory Viewer", "Directory_Viewer", "cis-central!b13"),
					ExpectError: regexp.MustCompile(`Access forbidden due to insufficient authorization.*`), //error message has a line break, we only check the first part
				},
			},
		})
	})

	t.Run("error path - directory_id, name, role_template_name and app_id are mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_directory_role" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "(directory_id|name|role_template_name|app_id)" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - directory_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceDirectoryRoleByDirectoryId("uut", "this-is-not-a-uuid", "a", "b", "c"),
					ExpectError: regexp.MustCompile(`Attribute directory_id value must be a valid UUID, got: this-is-not-a-uuid`),
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
					Config:      hclDatasourceDirectoryRole("uut", "integration-test-dir-se-static", "", "b", "c"),
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
					Config:      hclDatasourceDirectoryRoleByDirectoryId("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "a", "", "c"),
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
					Config:      hclDatasourceDirectoryRole("uut", "integration-test-dir-se-static", "a", "b", ""),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceDirectoryRoleByDirectoryId("uut", "00000000-0000-0000-0000-000000000000", "Directory Viewer", "Directory_Viewer", "cis-central!b13"),
					ExpectError: regexp.MustCompile(`received response with unexpected status: 404 \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceDirectoryRoleByDirectoryId(resourceName string, directoryId string, name string, roleTemplateName string, appId string) string {
	template := `
data "btp_directory_role" "%s" {
    directory_id       = "%s"
    name               = "%s"
    role_template_name = "%s"
    app_id             = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId, name, roleTemplateName, appId)
}

func hclDatasourceDirectoryRole(resourceName string, directoryName string, name string, roleTemplateName string, appId string) string {
	template := `
data "btp_directories" "all" {}
data "btp_directory_role" "%s" {
    directory_id       = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
    name               = "%s"
    role_template_name = "%s"
    app_id             = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryName, name, roleTemplateName, appId)
}
