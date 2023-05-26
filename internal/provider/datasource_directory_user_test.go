package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceDirectoryUser(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_directory_user")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceDirectoryUser("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "jenny.doe@test.com", "sap.default"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "directory_id", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "origin", "sap.default"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "active", "false"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "family_name", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "given_name", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "id", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "role_collections.#", "0"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "verified", "false"),
					),
				},
			},
		})
	})
	/*
		t.Run("error path - name, role_template_name and app_id are mandatory", func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: getProviders(nil),
				Steps: []resource.TestStep{
					{
						Config:      hclProvider() + `data "btp_globalaccount_role" "uut" {}`,
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
						Config:      hclProvider() + hclDatasourceGlobalaccountRole("uut", "", "b", "c"),
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
						Config:      hclProvider() + hclDatasourceGlobalaccountRole("uut", "a", "", "c"),
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
						Config:      hclProvider() + hclDatasourceGlobalaccountRole("uut", "a", "b", ""),
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
						Config:      hclProviderWithCLIServerURL(srv.URL) + hclDatasourceGlobalaccountRole("uut", "Global Account Viewer", "GlobalAccount_Viewer", "cis-local!b13"),
						ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
					},
				},
			})
		})
	*/
}

func hclDatasourceDirectoryUser(resourceName string, directoryId string, userName string, origin string) string {
	template := `data "btp_directory_user" "%s" {
	directory_id = "%s"
	user_name 	 = "%s"
  	origin    	 = "%s"
}`

	return fmt.Sprintf(template, resourceName, directoryId, userName, origin)
}
