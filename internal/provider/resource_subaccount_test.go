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

func TestResourceSubaccount(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceSubaccount("uut", "integration-test-acc-dyn", "eu12", "integration-test-acc-dyn"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "parent_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", "john.doe@int.test"),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "UNSET"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "beta_enabled", "false"),
					),
				},
				{
					Config: hclProvider() + hclResourceSubaccount("uut", "Integration Test Acc Dyn", "eu12", "integration-test-acc-dyn"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "Integration Test Acc Dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "parent_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", "john.doe@int.test"),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "UNSET"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "beta_enabled", "false"),
					),
				},
				{
					ResourceName:      "btp_subaccount.uut",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - usage", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_subaccount.usage_set")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceSubaccountForProdUsage("uut", "integration-test-acc-dyn-2", "eu12", "integration-test-acc-dyn-2"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "integration-test-acc-dyn-2"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "parent_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn-2"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", "john.doe@int.test"),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "USED_FOR_PRODUCTION"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "beta_enabled", "false"),
					),
				},
			},
		})
	})

	t.Run("error path - parent_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + hclResourceSubaccountWithParent("uut", "this-is-not-a-uuid", "a-subaccount", "eu12", "a-subaccount"),
					ExpectError: regexp.MustCompile(`Attribute parent_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
	t.Run("error path - name must not contain slashes", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + hclResourceSubaccount("uut", "a/subaccount", "eu12", "a-subaccount"),
					ExpectError: regexp.MustCompile(`Attribute name must not contain '/', not be empty and not exceed 255`),
				},
			},
		})
	})
	t.Run("error path - subdomain must be valid", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + hclResourceSubaccount("uut", "a.subaccount", "eu12", "a.subaccount"),
					ExpectError: regexp.MustCompile(`Attribute subdomain must only contain letters \(a-z\), digits \(0-9\)`),
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
					Config:      hclProviderWithCLIServerURL(srv.URL) + hclResourceSubaccount("uut", "a-subaccount", "eu12", "a-subaccount"),
					ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclResourceSubaccount(resourceName string, displayName string, region string, subdomain string) string {
	template := `
resource "btp_subaccount" "%s" {
    name      = "%s"
    region    = "%s"
    subdomain = "%s"
}`

	return fmt.Sprintf(template, resourceName, displayName, region, subdomain)
}

func hclResourceSubaccountForProdUsage(resourceName string, displayName string, region string, subdomain string) string {
	template := `
resource "btp_subaccount" "%s" {
    name      = "%s"
    region    = "%s"
    subdomain = "%s"
	used_for_production = "true"
}`

	return fmt.Sprintf(template, resourceName, displayName, region, subdomain)
}

func hclResourceSubaccountWithParent(resourceName string, parentId string, displayName string, region string, subdomain string) string {
	template := `
resource "btp_subaccount" "%s" {
    parent_id = "%s"
    name      = "%s"
    region    = "%s"
    subdomain = "%s"
}`

	return fmt.Sprintf(template, resourceName, parentId, displayName, region, subdomain)
}
