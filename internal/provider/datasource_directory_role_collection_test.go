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

func TestDataSourceDirectoryRoleCollection(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_role_collection")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDirectoryRoleCollection("uut", "integration-test-dir-se-static", "Directory Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_directory_role_collection.uut", "name", "Directory Viewer"),
						resource.TestCheckResourceAttr("data.btp_directory_role_collection.uut", "description", "Read-only access to the directory"),
						resource.TestCheckResourceAttr("data.btp_directory_role_collection.uut", "read_only", "true"),
						resource.TestCheckResourceAttr("data.btp_directory_role_collection.uut", "roles.#", "3"),
					),
				},
			},
		})
	})
	t.Run("error path - directory not security enabled", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_role_collection.not_security_enabled")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceDirectoryRoleCollection("uut", "integration-test-dir-static", "Directory Viewer"),
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
					Config:      `data "btp_directory_role_collections" "uut" {}`,
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
					Config:      hclDatasourceDirectoryRoleCollectionByDirectoryId("uut", "this-is-not-a-uuid", "Directory Viewer"),
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
					Config:      hclDatasourceDirectoryRoleCollectionByDirectoryId("uut", "5357bda0-8651-4eab-a69d-12d282bc3247", ""),
					ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1, got: 0`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceDirectoryRoleCollectionByDirectoryId("uut", "5357bda0-8651-4eab-a69d-12d282bc3247", "Directory Viewer"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceDirectoryRoleCollectionByDirectoryId(resourceName string, directoryId string, name string) string {
	template := `
data "btp_directory_role_collection" "%s" {
    directory_id = "%s"
    name         = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId, name)
}

func hclDatasourceDirectoryRoleCollection(resourceName string, directoryName string, name string) string {
	template := `
data "btp_directories" "all" {}
data "btp_directory_role_collection" "%s" {
    directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
    name         = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryName, name)
}
