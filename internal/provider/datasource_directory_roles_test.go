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

func TestDataSourceDirectoryRoles(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_roles")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDirectoryRoles("uut", "integration-test-dir-se-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_directory_roles.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_directory_roles.uut", "values.#", "8"),
						resource.TestCheckResourceAttrSet("data.btp_directory_roles.uut", "values.0.app_name"),
					),
				},
			},
		})
	})
	t.Run("error path - directory not security enabled", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_roles.not_security_enabled")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceDirectoryRoles("uut", "integration-test-dir-static"),
					ExpectError: regexp.MustCompile(`Access forbidden due to insufficient authorization.*`), //error message has a line break, we only check the first part
				},
			},
		})
	})
	t.Run("error path - directory_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_directory_roles" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "directory_id" is required, but no definition was found`),
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
					Config:      hclDatasourceDirectoryRolesByDirectoryId("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute directory_id value must be a valid UUID, got: this-is-not-a-uuid`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceDirectoryRolesByDirectoryId("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceDirectoryRolesByDirectoryId(resourceName string, directoryId string) string {
	return fmt.Sprintf(`data "btp_directory_roles" "%s" { directory_id = "%s" }`, resourceName, directoryId)
}

func hclDatasourceDirectoryRoles(resourceName string, directoryName string) string {
	template := `
data "btp_directories" "all" {}
data "btp_directory_roles" "%s" {
	directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, directoryName)
}
