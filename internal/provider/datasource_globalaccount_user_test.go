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

func TestDataSourceGlobalaccountUser(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_user")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountUserWithCustomIdp("uut", "jenny.doe@test.com", "sap.default"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "origin", "sap.default"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "active", "true"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "family_name", "unknown"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "given_name", "unknown"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "id", "86535387-54aa-4282-af13-67dd50cdd13c"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "role_collections.#", "1"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_user.uut", "verified", "false"),
					),
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
					Config:      hclDatasourceGlobalaccountUser("uut", ""),
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
					Config:      hclDatasourceGlobalaccountUserWithCustomIdp("uut", "jenny.doe@test.com", ""),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceGlobalaccountUser("uut", "jenny.doe@test.com"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountUser(resourceName string, userName string) string {
	return fmt.Sprintf(`data "btp_globalaccount_user" "%s" { user_name = "%s" }`, resourceName, userName)
}

func hclDatasourceGlobalaccountUserWithCustomIdp(resourceName string, userName string, origin string) string {
	template := `
data "btp_globalaccount_user" "%s" {
    user_name = "%s"
    origin    = "%s"
}`

	return fmt.Sprintf(template, resourceName, userName, origin)
}
