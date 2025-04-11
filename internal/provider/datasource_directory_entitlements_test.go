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

func TestDataSourceDirectoryEntitlements(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_entitlements")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDirectoryEntitlements("uut", "integration-test-dir-se-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_directory_entitlements.uut", "values.%", "2"),
					),
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
					Config:      hclDatasourceDirectoryEntitlementsByDirectoryId("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute directory_id value must be a valid UUID, got: this-is-not-a-uuid`),
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
					Config:      `data "btp_directory_entitlements" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "directory_id" is required, but no definition was found`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceDirectoryEntitlements("uut", "integration-test-dir-se-static"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceDirectoryEntitlementsByDirectoryId(resourceName string, directoryId string) string {
	template := `
data "btp_directory_entitlements" "%s" {
    directory_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId)
}

func hclDatasourceDirectoryEntitlements(resourceName string, directoryName string) string {
	template := `
data "btp_directories" "all" {}
data "btp_directory_entitlements" "%s" {
    directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, directoryName)
}
