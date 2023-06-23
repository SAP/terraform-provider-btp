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
		rec := setupVCR(t, "fixtures/datasource_directory_user.default_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceDirectoryUserDefaultIdp("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "jenny.doe@test.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "directory_id", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "origin", "ldap"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "active", "true"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "family_name", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "given_name", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "id", "40c72ef9-b901-4b89-91fb-3d283231f7b4"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "role_collections.#", "0"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "verified", "false"),
					),
				},
			},
		})
	})
	t.Run("happy path - custom idp", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_directory_user.custom_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceDirectoryUserCustomIdp("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "jenny.doe@test.com", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "directory_id", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "origin", "terraformint-platform"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "active", "true"),
                        resource.TestCheckResourceAttr("data.btp_directory_user.uut", "family_name", "unknown"), //FIXME should be empty, see NGPBUG-357810
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "given_name", "unknown"), //FIXME should be empty, see NGPBUG-357810
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "id", "2b5382f4-1922-4803-8dcb-5babe097b12b"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "role_collections.#", "0"),
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
					Config:      hclProvider() + hclDatasourceDirectoryUserDefaultIdp("uut", "this-is-not-a-uuid", "jenny.doe@test.com"),
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
					Config:      hclProvider() + `data "btp_directory_user" "uut" {}`,
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
					Config:      hclProvider() + hclDatasourceDirectoryUserCustomIdp("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "", "terraformint"),
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
					Config:      hclProvider() + hclDatasourceDirectoryUserCustomIdp("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "jenny.doe@test.com", ""),
					ExpectError: regexp.MustCompile(`Attribute origin string length must be at least 1, got: 0`),
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
					Config:      hclProviderWithCLIServerURL(srv.URL) + hclDatasourceDirectoryUserDefaultIdp("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "jenny.doe@test.com"),
					ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceDirectoryUserCustomIdp(resourceName string, directoryId string, userName string, origin string) string {
	template := `data "btp_directory_user" "%s" {
	directory_id = "%s"
	user_name 	 = "%s"
  	origin    	 = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId, userName, origin)
}

func hclDatasourceDirectoryUserDefaultIdp(resourceName string, directoryId string, userName string) string {
	template := `
data "btp_directory_user" "%s" {
	directory_id = "%s"
	user_name 	 = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId, userName)
}
