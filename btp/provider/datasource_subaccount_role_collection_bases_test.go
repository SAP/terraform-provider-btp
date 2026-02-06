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

func TestDataSourceSubaccountRoleCollectionBases(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_role_collection_bases")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountRoleCollectionBases("uut", "integration-test-acc-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_role_collection_bases.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_role_collection_bases.uut", "values.#", "7"),
						resource.TestCheckTypeSetElemNestedAttrs("data.btp_subaccount_role_collection_bases.uut", "values.*", map[string]string{
							"name": "Subaccount Viewer",
						}),
					),
				},
			},
		})
	})

	t.Run("error path - cli server returns 404", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/login/") || strings.Contains(r.URL.Path, "accounts/global") {
				w.Header().Set("Content-Type", "application/json")
				_, _ = fmt.Fprintf(w, `{"proxyType": "none", "isLicensed": true}`)
				return
			}
			w.WriteHeader(http.StatusNotFound)
			_, _ = fmt.Fprintf(w, `{"error": "No role collections found for subaccount"}`)
		}))
		defer srv.Close()

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(srv.Client()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderForCLIServerAt(srv.URL) + `
	data "btp_subaccount_role_collection_bases" "uut" {
		subaccount_id = "00000000-0000-0000-0000-000000000000"
	}`,
					ExpectError: regexp.MustCompile(`received response with unexpected status: 404`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountRoleCollectionBases(resourceName string, subaccountName string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
data "btp_subaccount_role_collection_bases" "%s" {
    subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
}`, resourceName, subaccountName)
}
