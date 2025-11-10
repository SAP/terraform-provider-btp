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

	t.Run("happy path - subaccount by subdomain", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_by_subdomain")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountBySubdomain("test", "eu12", "integration-test-acc-static-b8xxozer"),
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

	t.Run("error path - id or subdomain with region is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_subaccount" "test" {}`,
					ExpectError: regexp.MustCompile(`At least one attribute out of \[id,subdomain\] must be specified`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceSubaccountById("test", "00000000-0000-0000-0000-000000000000"),
					ExpectError: regexp.MustCompile(`received response with unexpected status: 404 \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})

	t.Run("error path - either id or subdomain with region should be provided, not both", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: `
					data "btp_subaccount" "test" {
						id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
						subdomain = "integration-test-acc-static-b8xxozer"
						region = "eu12"
					}`,
					ExpectError: regexp.MustCompile(`Attribute "subdomain" cannot be specified when "id" is specified`),
				},
			},
		})
	})

	t.Run("error path - region not provided with subdomain", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_subaccount" "test" {subdomain = "integration-test-acc-static-b8xxozer"}`,
					ExpectError: regexp.MustCompile(`Attribute "region" must be specified when "subdomain" is specified`),
				},
			},
		})
	})

	t.Run("error path - subdomain is invalid", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount.err_invalid_subdomain")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceSubaccountBySubdomain("test", "eu12", "invalid-subdomain"),
					ExpectError: regexp.MustCompile(`Subaccount not found with subdomain`),
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

func hclDatasourceSubaccountBySubdomain(resourceName string, region string, subDomain string) string {
	template := `
	data "btp_subaccount" "%s" {
		region = "%s"
		subdomain = "%s"
	}`
	return fmt.Sprintf(template, resourceName, region, subDomain)
}
