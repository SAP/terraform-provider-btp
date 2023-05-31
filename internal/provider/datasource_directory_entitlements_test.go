package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceDirectoryEntitlements(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_directory_entitlements")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceDirectoryEntitlements("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_directory_entitlements.uut", "id", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
						resource.TestCheckResourceAttr("data.btp_directory_entitlements.uut", "directory_id", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
						resource.TestCheckResourceAttr("data.btp_directory_entitlements.uut", "values.%", "2"),
						/* TBD: Add other entitlement details into resource attributes
						resource.TestCheckResourceAttr("data.btp_directory_entitlements.uut", "plan_description", "Directory Viewer"),
						resource.TestCheckResourceAttr("data.btp_directory_entitlements.uut", "plan_display_name", "Read-only access to the directory"),
						resource.TestCheckResourceAttr("data.btp_directory_entitlements.uut", "quota_assigned", "0"),
						resource.TestCheckResourceAttr("data.btp_directory_entitlements.uut", "quota_remaining", "3"),
						resource.TestCheckResourceAttr("data.btp_directory_entitlements.uut", "service_display_name", "3"),
						resource.TestCheckResourceAttr("data.btp_directory_entitlements.uut", "service_name", "3"),
						*/
					),
				},
			},
		})
	})
	/*
		t.Run("error path - directory not security enabled", func(t *testing.T) {
			rec := setupVCR(t, "fixtures/datasource_directory_role_collection.not_security_enabled")
			defer stopQuietly(rec)

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
				Steps: []resource.TestStep{
					{
						Config:      hclProvider() + hclDatasourceDirectoryRoleCollection("uut", "5357bda0-8651-4eab-a69d-12d282bc3247", "Directory Viewer"),
						ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 400; Correlation ID:\s+[a-f0-9\-]+\]`),
					},
				},
			})
		})
		t.Run("error path - directory_id mandatory", func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: getProviders(nil),
				Steps: []resource.TestStep{
					{
						Config:      hclProvider() + `data "btp_directory_role_collections" "uut" {}`,
						ExpectError: regexp.MustCompile(`The argument "directory_id" is required, but no definition was found`),
					},
				},
			})
		})
		t.Run("error path - directory_id not a valid UUID", func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: getProviders(nil),
				Steps: []resource.TestStep{
					{
						Config:      hclProvider() + hclDatasourceDirectoryRoleCollection("uut", "this-is-not-a-uuid", "Directory Viewer"),
						ExpectError: regexp.MustCompile(`Attribute directory_id value must be a valid UUID, got: this-is-not-a-uuid`),
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
						Config:      hclProviderWithCLIServerURL(srv.URL) + hclDatasourceDirectoryRoleCollection("uut", "5357bda0-8651-4eab-a69d-12d282bc3247", "Directory Viewer"),
						ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
					},
				},
			})
		})
	*/
}

func hclDatasourceDirectoryEntitlements(resourceName string, directoryId string) string {
	template := `
data "btp_directory_entitlements" "%s" {
  directory_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId)
}
