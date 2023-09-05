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

func TestDataSourceSubaccountUsers(t *testing.T) {
	t.Parallel()
	t.Run("happy path with default idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_users.default_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountUsersDefaultIdp("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_users.uut", "subaccount_id", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
						resource.TestCheckResourceAttr("data.btp_subaccount_users.uut", "values.#", "6"),
					),
				},
			},
		})
	})
	t.Run("happy path with custom idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_users.custom_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountUsersWithCustomIdp("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_users.uut", "subaccount_id", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
						resource.TestCheckResourceAttr("data.btp_subaccount_users.uut", "values.#", "3"),
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
					Config:      `data "btp_subaccount_users" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
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
					Config:      hclDatasourceSubaccountUsersDefaultIdp("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
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
					Config:      hclDatasourceSubaccountUsersWithCustomIdp("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", ""),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceSubaccountUsersWithCustomIdp("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "terraformint-platform"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountUsersDefaultIdp(resourceName string, subaccountId string) string {
	template := `data "btp_subaccount_users" "%s" { subaccount_id = "%s" }`

	return fmt.Sprintf(template, resourceName, subaccountId)
}

func hclDatasourceSubaccountUsersWithCustomIdp(resourceName string, subaccountId string, origin string) string {
	template := `
data "btp_subaccount_users" "%s" {
    subaccount_id = "%s"
    origin        = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, origin)
}
