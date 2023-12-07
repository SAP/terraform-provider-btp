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

func TestDataSourceDirectory(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{ // normal directory
					Config: hclProviderFor(user) + hclDatasourceDirectory("uut", "integration-test-dir-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_directory.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttrSet("data.btp_directory.uut", "created_by"),
						//resource.TestCheckResourceAttr("data.btp_directory.uut", "created_date", "2023-05-16T08:39:33Z"),
						resource.TestCheckResourceAttr("data.btp_directory.uut", "description", "Please don't modify. This is used for integration tests."),
						resource.TestCheckResourceAttr("data.btp_directory.uut", "labels.#", "0"),
						resource.TestMatchResourceAttr("data.btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_directory.uut", "name", "integration-test-dir-static"),
						//resource.TestCheckResourceAttr("data.btp_directory.uut", "parent_id", "03760ecf-9d89-4189-a92a-1c7efed09298"),
						resource.TestCheckResourceAttr("data.btp_directory.uut", "state", "OK"),
						resource.TestCheckResourceAttr("data.btp_directory.uut", "subdomain", ""),
					),
				},
				{ // security enabled directory
					Config: hclProviderFor(user) + hclDatasourceDirectory("uut_security_enabled", "integration-test-dir-se-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_directory.uut_security_enabled", "id", regexpValidUUID),
						resource.TestCheckResourceAttrSet("data.btp_directory.uut_security_enabled", "created_by"),
						//resource.TestCheckResourceAttr("data.btp_directory.uut_security_enabled", "created_date", "2023-05-16T08:46:24Z"),
						resource.TestCheckResourceAttr("data.btp_directory.uut_security_enabled", "description", "Please don't modify. This is used for integration tests."),
						resource.TestCheckResourceAttr("data.btp_directory.uut_security_enabled", "labels.#", "0"),
						resource.TestMatchResourceAttr("data.btp_directory.uut_security_enabled", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_directory.uut_security_enabled", "name", "integration-test-dir-se-static"),
						//resource.TestCheckResourceAttr("data.btp_directory.uut_security_enabled", "parent_id", "03760ecf-9d89-4189-a92a-1c7efed09298"),
						resource.TestCheckResourceAttr("data.btp_directory.uut_security_enabled", "state", "OK"),
						//resource.TestCheckResourceAttr("data.btp_directory.uut_security_enabled", "subdomain", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
					),
				},
			},
		})
	})
	t.Run("error path - id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_directory" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found`),
				},
			},
		})
	})
	t.Run("error path - id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceDirectoryById("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute id value must be a valid UUID, got: this-is-not-a-uuid`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceDirectoryById("uut", "00000000-0000-0000-0000-000000000000"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceDirectoryById(resourceName string, id string) string {
	return fmt.Sprintf(`data "btp_directory" "%s" { id = "%s" }`, resourceName, id)
}

func hclDatasourceDirectory(resourceName string, directoryName string) string {
	template := `
data "btp_directories" "all" {}
data "btp_directory" "%s" {
    id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, directoryName)
}
