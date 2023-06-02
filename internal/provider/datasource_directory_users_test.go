package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceDirectoryUsers(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_directory_users.default_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceDirectoryUsersDefaultIdp("defaultidp", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_directory_users.defaultidp", "directory_id", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
						resource.TestCheckResourceAttr("data.btp_directory_users.defaultidp", "values.#", "2"),
					),
				},
			},
		})
	})
	/* TODO: find root cause for error message "cannot unmarshal object into Go value of type []string" (see https://github.com/SAP/terraform-provider-btp/issues/71)
	t.Run("happy path with custom idp", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_directory_users.custom_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceDirectoryUsersCustomIdp("mycustomidp", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "terraformint"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_directory_users.mycustomidp", "directory_id", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
						resource.TestCheckResourceAttr("data.btp_directory_users.mycustomidp", "values.#", "2"),
					),
				},
			},
		})
	})
	*/
	t.Run("error path - directory_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + `data "btp_directory_users" "uut" {}`,
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
					Config:      hclProvider() + hclDatasourceDirectoryUsersDefaultIdp("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute directory_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})

}

func hclDatasourceDirectoryUsersDefaultIdp(resourceName string, directoryId string) string {
	template := `data "btp_directory_users" "%s" {
  directory_id    = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId)
}

/* TODO: find root cause for error message "cannot unmarshal object into Go value of type []string" (see https://github.com/SAP/terraform-provider-btp/issues/71)
func hclDatasourceDirectoryUsersCustomIdp(resourceName string, directoryId string, origin string) string {
	template := `data "btp_directory_users" "%s" {
  directory_id    = "%s"
  origin    = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId, origin)
}
*/
