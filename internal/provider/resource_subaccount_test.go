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
		rec, user := setupVCR(t, "fixtures/resource_subaccount")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccount("uut", "integration-test-acc-dyn", "eu12", "integration-test-acc-dyn"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", user.Username),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "UNSET"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "beta_enabled", "false"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccount("uut", "Integration Test Acc Dyn", "eu12", "integration-test-acc-dyn"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "Integration Test Acc Dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", user.Username),
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
	t.Run("happy path used for prod", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount.used_for_production")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountUsedForProd("uut", "integration-test-acc-dyn", "eu12", "integration-test-acc-dyn"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", user.Username),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "USED_FOR_PRODUCTION"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "beta_enabled", "false"),
					),
				},
				{
					// Update name wo change of usage but provide usage explicitly again
					Config: hclProviderFor(user) + hclResourceSubaccountUsedForProd("uut", "Integration Test Acc Dyn", "eu12", "integration-test-acc-dyn"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "Integration Test Acc Dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", user.Username),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "USED_FOR_PRODUCTION"),
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

	t.Run("happy path change to used for prod", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount.change_to_used_for_production")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccount("uut", "integration-test-acc-dyn", "eu12", "integration-test-acc-dyn"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", user.Username),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "UNSET"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "beta_enabled", "false"),
					),
				},
				{
					//Update name wo change of usage but provide usage explicitly again
					Config: hclProviderFor(user) + hclResourceSubaccountUsedForProd("uut", "Integration Test Acc Dyn", "eu12", "integration-test-acc-dyn"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "Integration Test Acc Dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", user.Username),
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

	t.Run("happy path full config with update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount.full_config")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountAll("uut", "integration-test-acc-dyn", "eu12", "integration-test-acc-dyn", "My subaccount description", "NOT_USED_FOR_PRODUCTION", true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", "My subaccount description"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", user.Username),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "NOT_USED_FOR_PRODUCTION"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "beta_enabled", "true"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "labels.foo.0", "bar"),
					),
				},
				{
					// Update name wo change of usage but omit optional parameters
					Config: hclProviderFor(user) + hclResourceSubaccount("uut", "Integration Test Acc Dyn", "eu12", "integration-test-acc-dyn"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount.uut", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "name", "Integration Test Acc Dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "description", "My subaccount description"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "subdomain", "integration-test-acc-dyn"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "created_by", user.Username),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "state", "OK"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "usage", "NOT_USED_FOR_PRODUCTION"),
						resource.TestCheckResourceAttr("btp_subaccount.uut", "beta_enabled", "true"),
						resource.TestCheckNoResourceAttr("btp_subaccount.uut", "labels"),
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

	t.Run("error path - name must not contain slashes", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccount("uut", "a/subaccount", "eu12", "a-subaccount"),
					ExpectError: regexp.MustCompile(`Attribute name must not contain '/', not be empty and not exceed 255`),
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
					Config:      hclProviderForCLIServerAt(srv.URL) + hclResourceSubaccount("uut", "a-subaccount", "eu12", "a-subaccount"),
					ExpectError: regexp.MustCompile(`received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
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

func hclResourceSubaccountAll(resourceName string, displayName string, region string, subdomain string, description string, usage string, betaEnabled bool) string {
	template := `
resource "btp_subaccount" "%s" {
    name         = "%s"
    region       = "%s"
    subdomain    = "%s"
	description  = "%s"
	usage        = "%s"
	beta_enabled = %t
	labels	     = {"foo" = ["bar"]}
}`

	result := fmt.Sprintf(template, resourceName, displayName, region, subdomain, description, usage, betaEnabled)
	return result
}

func hclResourceSubaccountUsedForProd(resourceName string, displayName string, region string, subdomain string) string {
	template := `
resource "btp_subaccount" "%s" {
    name      = "%s"
    region    = "%s"
    subdomain = "%s"
	usage     = "USED_FOR_PRODUCTION"
}`

	return fmt.Sprintf(template, resourceName, displayName, region, subdomain)
}
