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

func TestDataSourceSubaccountRoleCollectionRoles(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_role_collection_roles")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountRoleCollectionRoles("uut", "integration-test-acc-static", "Subaccount Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_role_collection_roles.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_role_collection_roles.uut", "name", "Subaccount Viewer"),
						// Check that the set contains the expected number of roles (adjust based on your environment)
						resource.TestCheckResourceAttr("data.btp_subaccount_role_collection_roles.uut", "values.#", "8"),
						// Verify that a specific role exists within the set
						resource.TestCheckTypeSetElemNestedAttrs("data.btp_subaccount_role_collection_roles.uut", "values.*", map[string]string{
							"role_name":          "Subaccount Viewer",
							"role_template_name": "Subaccount_Viewer",
						}),
					),
				},
			},
		})
	})

	t.Run("error path - subaccount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_subaccount_role_collection_roles" "uut" { name = "Subaccount Viewer" }`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required`),
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
					Config:      hclDatasourceSubaccountRoleCollectionRolesBySubaccountId("uut", "00000000-0000-0000-0000-000000000000", ""),
					ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceSubaccountRoleCollectionRolesBySubaccountId("uut", "00000000-0000-0000-0000-000000000000", "Subaccount Viewer"),
					ExpectError: regexp.MustCompile(`received response with unexpected status: 404`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountRoleCollectionRoles(resourceName string, subaccountName string, collectionName string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}

data "btp_subaccount_role_collection_roles" "%s" {
    subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
    name          = "%s"
}`, resourceName, subaccountName, collectionName)
}

func hclDatasourceSubaccountRoleCollectionRolesBySubaccountId(resourceName string, subaccountId string, collectionName string) string {
	return fmt.Sprintf(`
data "btp_subaccount_role_collection_roles" "%s" {
    subaccount_id = "%s"
    name          = "%s"
}`, resourceName, subaccountId, collectionName)
}
