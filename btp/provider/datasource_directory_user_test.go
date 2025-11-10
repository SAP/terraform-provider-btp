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

func TestDataSourceDirectoryUser(t *testing.T) {
	t.Parallel()
	t.Run("happy path - default idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_user.default_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDirectoryUserDefaultIdp("uut", "integration-test-dir-se-static", "jenny.doe@test.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_directory_user.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "origin", "ldap"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "active", "true"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "family_name", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "given_name", ""),
						resource.TestMatchResourceAttr("data.btp_directory_user.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "role_collections.#", "1"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "verified", "false"),
					),
				},
			},
		})
	})
	t.Run("happy path - custom idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_directory_user.custom_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDirectoryUserCustomIdp("uut", "integration-test-dir-se-static", "jenny.doe@test.com", "iasprovidertestblr-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_directory_user.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "origin", "iasprovidertestblr-platform"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "active", "true"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "family_name", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "given_name", ""),
						resource.TestMatchResourceAttr("data.btp_directory_user.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "role_collections.#", "2"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "verified", "false"),
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
					Config:      hclDatasourceDirectoryUserDefaultIdpByDirectoryId("uut", "this-is-not-a-uuid", "jenny.doe@test.com"),
					ExpectError: regexp.MustCompile(`Attribute directory_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
	t.Run("error path - directory_id, user_name and origin are mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_directory_user" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "(directory_id|user_name)" is required, but no definition was found.`),
				},
			},
		})
	})
	t.Run("error path - user_name must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceDirectoryUserCustomIdp("uut", "integration-test-dir-se-static", "", "terraformint"),
					ExpectError: regexp.MustCompile(`Attribute user_name string length must be between 1 and 256, got: 0`),
				},
			},
		})
	})
	t.Run("error path - user_name must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceDirectoryUserCustomIdp("uut", "integration-test-dir-se-static", "jenny.doe@test.com", ""),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceDirectoryUserDefaultIdpByDirectoryId("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "jenny.doe@test.com"),
					ExpectError: regexp.MustCompile(`received response with unexpected status: 404 \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceDirectoryUserCustomIdp(resourceName string, directoryName string, userName string, origin string) string {
	template := `
data "btp_directories" "all" {}
data "btp_directory_user" "%s" {
	directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
	user_name 	 = "%s"
  	origin    	 = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryName, userName, origin)
}

func hclDatasourceDirectoryUserDefaultIdpByDirectoryId(resourceName string, directoryId string, userName string) string {
	template := `
data "btp_directory_user" "%s" {
	directory_id = "%s"
	user_name 	 = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId, userName)
}

func hclDatasourceDirectoryUserDefaultIdp(resourceName string, directoryName string, userName string) string {
	template := `
data "btp_directories" "all" {}
data "btp_directory_user" "%s" {
	directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
	user_name 	 = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryName, userName)
}
