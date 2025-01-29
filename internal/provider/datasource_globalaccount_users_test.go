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

func TestDataSourceGlobalaccountUsers(t *testing.T) {
	t.Parallel()
	t.Run("happy path - default idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_users.default_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountUsers("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_users.uut", "values.#", "35"),
					),
				},
			},
		})
	})
	t.Run("happy path - with custom idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_users.custom_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountUsersWithCustomIdp("uut", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_users.uut", "values.#", "8"),
					),
				},
			},
		})
	})
	// TODO: error path with non existing idp
	t.Run("error path - origin must not be empty if given", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceGlobalaccountUsersWithCustomIdp("uut", ""),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceGlobalaccountUsersWithCustomIdp("uut", "terraformint-platform"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountUsers(resourceName string) string {
	template := `data "btp_globalaccount_users" "%s" {}`
	return fmt.Sprintf(template, resourceName)
}

func hclDatasourceGlobalaccountUsersWithCustomIdp(resourceName string, origin string) string {
	template := `data "btp_globalaccount_users" "%s" { origin = "%s" }`

	return fmt.Sprintf(template, resourceName, origin)
}
