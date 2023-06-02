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
					),
				},
			},
		})
	})

	/*
		t.Run("error path - directory_id not a valid UUID", func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: getProviders(nil),
				Steps: []resource.TestStep{
					{
						Config:      hclProvider() + hclDatasourceDirectoryEntitlements("uut", "this-is-not-a-uuid"),
						ExpectError: regexp.MustCompile(`Attribute directory_id value must be a valid UUID, got: this-is-not-a-uuid`),
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
						Config:      hclProvider() + `data "btp_directory_entitlements" "uut" {}`,
						ExpectError: regexp.MustCompile(`The argument "directory_id" is required, but no definition was found`),
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
