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

func TestDataSourceDirectoryUsers(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_users.default_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDirectoryUsersDefaultIdp("uut", "integration-test-dir-se-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_directory_users.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_directory_users.uut", "values.#", "8"),
					),
				},
			},
		})
	})
	t.Run("happy path with custom idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_users.custom_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDirectoryUsersWithCustomIdp("uut", "integration-test-dir-se-static", "iasprovidertestblr-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_directory_users.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_directory_users.uut", "values.#", "2"),
					),
				},
			},
		})
	})
	t.Run("error path - non existing idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_users.non_existing_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceDirectoryUsersWithCustomIdp("uut", "integration-test-dir-se-static", "this-doesnt-exist"),
					ExpectError: regexp.MustCompile(`API Error Reading Resource Users \(Directory\)`),
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
					Config:      `data "btp_directory_users" "uut" {}`,
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
					Config:      hclDatasourceDirectoryUsersDefaultIdpByDirectoryId("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute directory_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
	t.Run("error path - origin must not be empty if given", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceDirectoryUsersWithCustomIdp("uut", "integration-test-dir-se-static", ""),
					ExpectError: regexp.MustCompile(`Attribute origin string length must be at least 1, got: 0`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceDirectoryUsersDefaultIdpByDirectoryId("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceDirectoryUsersDefaultIdpByDirectoryId(resourceName string, directoryId string) string {
	template := `data "btp_directory_users" "%s" { directory_id = "%s" }`
	return fmt.Sprintf(template, resourceName, directoryId)
}

func hclDatasourceDirectoryUsersDefaultIdp(resourceName string, directoryName string) string {
	template := `
data "btp_directories" "all" {}
data "btp_directory_users" "%s" {
	directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, directoryName)
}

func hclDatasourceDirectoryUsersWithCustomIdp(resourceName string, directoryName string, origin string) string {
	template := `
data "btp_directories" "all" {}
data "btp_directory_users" "%s" {
    directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
    origin       = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryName, origin)
}
