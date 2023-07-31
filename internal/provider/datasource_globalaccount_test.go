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

func TestDataSourceGlobalAccount(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalAccount("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "id", "03760ecf-9d89-4189-a92a-1c7efed09298"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "name", "terraform-integration-canary"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "contract_status", "ACTIVE"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "geo_access", "STANDARD"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "license_type", "SAPDEV"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "costobject_type", "COST_CENTER"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "usage", "Testing"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "commercial_model", "Subscription"),
						resource.TestCheckResourceAttr("data.btp_globalaccount.uut", "consumption_based", "true"),
					),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceGlobalAccount("uut"),
					ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceGlobalAccount(resourceName string) string {
	return fmt.Sprintf(`data "btp_globalaccount" "%s" {}`, resourceName)
}
