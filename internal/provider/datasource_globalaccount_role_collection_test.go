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

func TestDataSourceGlobalaccountRoleCollection(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_role_collection.role_collection_exists")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountRoleCollection("uut", "Global Account Administrator"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "description", "Administrative access to the global account"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "read_only", "true"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "roles.#", "4"),
					),
				},
			},
		})
	})
	// FIXME https://github.com/SAP/terraform-provider-btp/issues/160
	/*t.Run("happy path - role collection not available", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_globalaccount_role_collection.role_collection_not_available")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceGlobalaccountRoleCollection("uut", "fuh"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "description", ""),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "read_only", "false"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_role_collection.uut", "roles.#", "0"),
					),
				},
			},
		})
	})*/
	t.Run("error path - name must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + hclDatasourceGlobalaccountRoleCollection("uut", ""),
					ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1, got: 0`),
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
					Config:      hclProviderWithCLIServerURL(srv.URL) + hclDatasourceGlobalaccountRoleCollection("uut", "Global Account Administrator"),
					ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountRoleCollection(resourceName string, name string) string {
	template := `data "btp_globalaccount_role_collection" "%s" { name = "%s" }`

	return fmt.Sprintf(template, resourceName, name)
}
