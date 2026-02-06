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

func TestDataSourceGlobalaccountEntitlementsWithDataCenter(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_entitlements_with_data_center")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountEntitlementsWithDataCenter("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlements_with_data_centers.uut", "values.%", "181"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlements_with_data_centers.uut", "values.SAPLaunchpad:foundation.datacenter_information.%", "5"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlements_with_data_centers.uut", "values.SAPLaunchpad:foundation.datacenter_information.eu12.dc_region", "eu12"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlements_with_data_centers.uut", "values.SAPLaunchpad:foundation.datacenter_information.eu12.dc_name", "cf-eu12"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlements_with_data_centers.uut", "values.SAPLaunchpad:foundation.datacenter_information.eu12.dc_display_name", "cf-eu12"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlements_with_data_centers.uut", "values.SAPLaunchpad:foundation.datacenter_information.eu12.dc_iaas_provider", "AWS"),
					),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclDatasourceGlobalaccountEntitlementsWithDataCenter("uut"),
					ExpectError: regexp.MustCompile(`received response with unexpected status: 404 \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountEntitlementsWithDataCenter(resourceName string) string {
	template := `data "btp_globalaccount_entitlements_with_data_centers" "%s" {}`
	return fmt.Sprintf(template, resourceName)
}
