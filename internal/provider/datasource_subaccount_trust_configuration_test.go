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

func TestDataSourceSubaccountTrustConfiguration(t *testing.T) {
	t.Parallel()
	t.Run("happy path - with default idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_trust_configuration.default")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountTrustConfiguration("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "sap.default"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "id", "sap.default"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "description", ""),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "identity_provider", ""),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "name", "sap.default"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "read_only", "false"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "type", "Application"),
					),
				},
			},
		})
	})
	t.Run("happy path - with custom idp", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_trust_configuration.custom_idp_exists")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountTrustConfiguration("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "id", "terraformint-platform"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "description", "Custom Platform Identity Provider"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "identity_provider", "terraformint.accounts400.ondemand.com"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "name", "terraformint-platform"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "protocol", "OpenID Connect"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "read_only", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "status", "active"),
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configuration.uut", "type", "Platform"),
					),
				},
			},
		})
	})
	t.Run("error path - custom idp not existing", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_trust_configuration.custom_idp_not_existing")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceSubaccountTrustConfiguration("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "fuh"),
					ExpectError: regexp.MustCompile(`API Error Reading Resource Trust Configuration \(Subaccount\)`),
				},
			},
		})
	})
	t.Run("error path - origin must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountTrustConfiguration("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", ""),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceSubaccountTrustConfiguration("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "sap.default"),
					ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountTrustConfiguration(resourceName string, subaccountId string, origin string) string {
	template := `
data "btp_subaccount_trust_configuration" "%s" {
    subaccount_id = "%s"
	origin = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, origin)
}
