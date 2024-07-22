package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceDirectoryRole(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRole("uut", "integration-test-dir-roles", "Directory Viewer Test", "Directory_Viewer", "cis-central!b13"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role.uut", "name", "Directory Viewer Test"),
						resource.TestCheckResourceAttr("btp_directory_role.uut", "role_template_name", "Directory_Viewer"),
						resource.TestCheckResourceAttr("btp_directory_role.uut", "app_id", "cis-central!b13"),
						resource.TestCheckResourceAttr("btp_directory_role.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_directory_role.uut", "read_only", "false"),
						resource.TestCheckResourceAttr("btp_directory_role.uut", "scopes.#", "7"),
					),
				},
				{
					ResourceName:      "btp_directory_role.uut",
					ImportStateIdFunc: getIdForDirectoryRoleImportId("btp_directory_role.uut", "Directory Viewer Test", "Directory_Viewer", "cis-central!b13"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("error path - directory not security enabled", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role.not_security_enabled")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclResourceDirectoryRole("uut", "integration-test-dir-static", "Directory Viewer", "Directory_Viewer", "cis-central!b13"),
					ExpectError: regexp.MustCompile(`Access forbidden due to insufficient authorization.*`), //error message has a line break, we only check the first part
				},
			},
		})
	})

	t.Run("error path - directory_id, name, role_template_name and app_id are mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_directory_role" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "(directory_id|name|role_template_name|app_id)" is required, but no definition was found.`),
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
					Config:      hclResourceDirectoryRoleByDirectoryId("uut", "this-is-not-a-uuid", "a", "b", "c"),
					ExpectError: regexp.MustCompile(`Attribute directory_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})

	t.Run("error path - name must not be empty", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceDirectoryRole("uut", "integration-test-dir-roles", "", "b", "c"),
					ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1, got: 0`),
				},
			},
		})
	})

	t.Run("error path - role_template_name must not be empty", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceDirectoryRole("uut", "integration-test-dir-roles", "a", "", "c"),
					ExpectError: regexp.MustCompile(`Attribute role_template_name string length must be at least 1, got: 0`),
				},
			},
		})
	})

	t.Run("error path - app_id must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceDirectoryRole("uut", "integration-test-dir-roles", "a", "b", ""),
					ExpectError: regexp.MustCompile(`Attribute app_id string length must be at least 1, got: 0`),
				},
			},
		})
	})

	t.Run("error path - update role name", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role.update")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRole(
						"uut",
						"integration-test-dir-roles",
						"Directory Viewer Test",
						"Directory_Viewer",
						"cis-central!b13",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_directory_role.uut", "name", "Directory Viewer Test"),
						resource.TestCheckResourceAttr("btp_directory_role.uut", "role_template_name", "Directory_Viewer"),
						resource.TestCheckResourceAttr("btp_directory_role.uut", "app_id", "cis-central!b13"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRole(
						"uut",
						"integration-test-dir-roles",
						"Directory Viewer Test Updated",
						"Directory_Viewer",
						"cis-central!b13",
					),
					ExpectError: regexp.MustCompile(`This resource is not supposed to be updated`),
				},
			},
		})
	})

	t.Run("error path - import fails", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role.error_import")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRole("uut", "integration-test-dir-roles", "Directory Viewer Test", "Directory_Viewer", "cis-central!b13"),
				},
				{
					ResourceName:      "btp_directory_role.uut",
					ImportStateIdFunc: getDirectoryRoleImportIdNoAppIdNoRoleTemplateName("btp_directory_role.uut"),
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Unexpected Import Identifier`),
				},
			},
		})
	})
}

func hclResourceDirectoryRole(resourceName string, directoryName string, name string, roleTemplateName string, appId string) string {
	template := `
data "btp_directories" "all" {}
resource "btp_directory_role" "%s" {
    directory_id       = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
    name               = "%s"
    role_template_name = "%s"
    app_id             = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryName, name, roleTemplateName, appId)
}

func hclResourceDirectoryRoleByDirectoryId(resourceName string, directoryId string, name string, roleTemplateName string, appId string) string {
	template := `
resource "btp_directory_role" "%s" {
    directory_id       = "%s"
    name               = "%s"
    role_template_name = "%s"
    app_id             = "%s"
}`
	return fmt.Sprintf(template, resourceName, directoryId, name, roleTemplateName, appId)
}

func getDirectoryRoleImportIdNoAppIdNoRoleTemplateName(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.Attributes["directory_id"], nil
	}
}

func getIdForDirectoryRoleImportId(resourceName string, name string, role_template_name string, app_id string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s,%s,%s,%s", rs.Primary.Attributes["directory_id"], name, role_template_name, app_id), nil
	}
}
