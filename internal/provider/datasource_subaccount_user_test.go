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

func TestDataSourceSubaccountUser(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_user")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountUserWithCustomIdp("uut", "integration-test-acc-static", "jenny.doe@test.com", "sap.default"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_user.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_user.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("data.btp_subaccount_user.uut", "origin", "sap.default"),
						resource.TestCheckResourceAttr("data.btp_subaccount_user.uut", "active", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_user.uut", "family_name", ""),
						resource.TestCheckResourceAttr("data.btp_subaccount_user.uut", "given_name", ""),
						resource.TestMatchResourceAttr("data.btp_subaccount_user.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_user.uut", "role_collections.#", "1"),
						resource.TestCheckResourceAttr("data.btp_subaccount_user.uut", "verified", "false"),
					),
				},
			},
		})
	})
	t.Run("error path - subaccount_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountUserWithCustomIdpBySubaccountId("uut", "this-is-not-a-uuid", "jenny.doe@test.com", "sap.default"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
	t.Run("error path - subaccount_id, user_name and origin are mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_subaccount_user" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "(subaccount_id|user_name)" is required, but no definition was found.`),
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
					Config:      hclDatasourceSubaccountUserWithCustomIdp("uut", "integration-test-acc-static", "", "terraformint"),
					ExpectError: regexp.MustCompile(`Attribute user_name string length must be between 1 and 256, got: 0`),
				},
			},
		})
	})
	t.Run("error path - origin must not be empty if given", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountUserWithCustomIdp("uut", "integration-test-acc-static", "jenny.doe@test.com", ""),
					ExpectError: regexp.MustCompile(`Attribute origin string length must be at least 1, got: 0`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceSubaccountUserWithCustomIdp("uut", "integration-test-acc-static", "jenny.doe@test.com", "sap.default"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountUserWithCustomIdpBySubaccountId(resourceName string, subaccountId string, userName string, origin string) string {
	template := `
data "btp_subaccount_user" "%s" {
	subaccount_id = "%s"
	user_name 	 = "%s"
  	origin    	 = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, userName, origin)
}

func hclDatasourceSubaccountUserWithCustomIdp(resourceName string, subaccountName string, userName string, origin string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_user" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	user_name 	 = "%s"
  	origin    	 = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, userName, origin)
}
