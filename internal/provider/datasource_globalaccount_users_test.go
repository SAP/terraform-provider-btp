package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountUsers(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_users")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountUsers("defaultidp", "sap.default"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_users.defaultidp", "values.#", "6"),
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

func hclDatasourceGlobalaccountUsers(resourceName string, origin string) string {
	template := `data "btp_globalaccount_users" "%s" {
  origin    = "%s"
}`
	return fmt.Sprintf(template, resourceName, origin)
}
