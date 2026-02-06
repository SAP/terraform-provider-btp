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

func TestDataSourceSubaccountRoleCollectionRole(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_role_collection_role")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					// Use the new helper that includes role_name
					Config: hclProviderFor(user) + hclDatasourceSubaccountRoleCollectionRole("uut", "integration-test-acc-static", "Subaccount Viewer", "Subaccount Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_role_collection_role.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_role_collection_role.uut", "name", "Subaccount Viewer"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role_collection_role.uut", "role_name", "Subaccount Viewer"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role_collection_role.uut", "role_template_name", "Subaccount_Viewer"),
						resource.TestCheckResourceAttr("data.btp_subaccount_role_collection_role.uut", "role_template_app_id", "cis-local!b2"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_role_collection_role.uut", "description"),
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
					Config: `
data "btp_subaccount_role_collection_role" "uut" {
    name      = "coll"
    role_name = "role"
}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required`),
				},
			},
		})
	})

	t.Run("error path - role not found in collection", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_role_collection_role.not_found")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceSubaccountRoleCollectionRole("uut", "integration-test-acc-static", "Subaccount Viewer", "NonExistentRole"),
					ExpectError: regexp.MustCompile(`Role 'NonExistentRole' not found in role collection 'Subaccount Viewer'`),
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
					Config:      hclDatasourceSubaccountRoleCollectionRoleBySubaccountId("uut", "this-is-not-a-uuid", "", "NonExistentRole"),
					ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1, got: 0`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceSubaccountRoleCollectionBySubaccountId("uut", "00000000-0000-0000-0000-000000000000", "Subaccount Viewer"),
					ExpectError: regexp.MustCompile(`received response with unexpected status: 404 \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountRoleCollectionRole(resourceName string, subaccountName string, collectionName string, roleName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_role_collection_role" "%s" {
    subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
    name          = "%s"
    role_name     = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, collectionName, roleName)
}

func hclDatasourceSubaccountRoleCollectionRoleBySubaccountId(resourceName string, subaccountId string, collectionName string, roleName string) string {
	template := `
data "btp_subaccount_role_collection_role" "%s" {
    subaccount_id = "%s"
    name          = "%s"
    role_name     = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, collectionName, roleName)
}
