package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceDirectoryUser(t *testing.T) {
	t.Parallel()
	t.Run("happy path - default idp", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_directory_user.default_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceDirectoryUserDefaultIdp("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "jenny.doe@test.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "directory_id", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "origin", "ldap"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "active", "false"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "family_name", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "given_name", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "id", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "role_collections.#", "0"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "verified", "false"),
					),
				},
			},
		})
	})
	t.Run("happy path - custom idp", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_directory_user.custom_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceDirectoryUserCustomIdp("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "jenny.doe@test.com", "terraformint"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "directory_id", "05368777-4934-41e8-9f3c-6ec5f4d564b9"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "origin", "terraformint"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "active", "false"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "family_name", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "given_name", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "id", ""),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "role_collections.#", "0"),
						resource.TestCheckResourceAttr("data.btp_directory_user.uut", "verified", "false"),
					),
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
					Config:      hclProvider() + hclDatasourceDirectoryUserDefaultIdp("uut", "this-is-not-a-uuid", "jenny.doe@test.com"),
					ExpectError: regexp.MustCompile(`Attribute directory_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
	t.Run("error path - directory_id, user_name and origin are mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + `data "btp_directory_user" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "(directory_id|user_name)" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclDatasourceDirectoryUserCustomIdp(resourceName string, directoryId string, userName string, origin string) string {
	template := `data "btp_directory_user" "%s" {
	directory_id = "%s"
	user_name 	 = "%s"
  	origin    	 = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId, userName, origin)
}

func hclDatasourceDirectoryUserDefaultIdp(resourceName string, directoryId string, userName string) string {
	template := `
data "btp_directory_user" "%s" {
	directory_id = "%s"
	user_name 	 = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId, userName)
}
