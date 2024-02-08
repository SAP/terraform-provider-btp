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

func TestDataSourceSubaccount(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountByName("test", "integration-test-acc-static"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount.test", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount.test", "beta_enabled", "false"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount.test", "created_by"),
						resource.TestMatchResourceAttr("data.btp_subaccount.test", "created_date", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_subaccount.test", "description", "Please don't modify. This is used for integration tests."),
						resource.TestCheckResourceAttr("data.btp_subaccount.test", "labels.#", "0"),
						resource.TestMatchResourceAttr("data.btp_subaccount.test", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_subaccount.test", "name", "integration-test-acc-static"),
						resource.TestMatchResourceAttr("data.btp_subaccount.test", "parent_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount.test", "region", "eu12"),
						resource.TestCheckResourceAttr("data.btp_subaccount.test", "state", "OK"),
						resource.TestCheckResourceAttr("data.btp_subaccount.test", "usage", "NOT_USED_FOR_PRODUCTION"),
					),
				},
			},
		})
	})

	t.Run("error path - subaccount doesn't exist", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount.err_subaccount_doesnt_exist")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceSubaccountById("test", "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"),
					ExpectError: regexp.MustCompile(`404 Not Found: \[no body\] \[Error: 404\]`), // TODO improve error text
				},
			},
		})
	})

	t.Run("error path - id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_subaccount" "test" {}`,
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountById("test", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute id value must be a valid UUID, got: this-is-not-a-uuid`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceSubaccountById("test", "00000000-0000-0000-0000-000000000000"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountById(resourceName string, id string) string {
	return fmt.Sprintf(`data "btp_subaccount" "%s" { id = "%s" }`, resourceName, id)
}

func hclDatasourceSubaccountByName(resourceName string, Name string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount" "%s" {
	id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, Name)
}
