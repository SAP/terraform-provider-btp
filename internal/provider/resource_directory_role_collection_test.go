package provider

import (
	"encoding/json"
	"fmt"

	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Needed for JSON mapping - fails with data types of directoryRoleCollectionRoleRefType struct
type directoryRoleCollectionRoleRefTestType struct {
	Name              string `json:"name"`
	RoleTemplateAppId string `json:"role_template_app_id"`
	RoleTemplateName  string `json:"role_template_name"`
}

func TestResourceDirectoryRoleCollection(t *testing.T) {

	t.Run("happy path - no description", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection.no_description")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionNoDescriptionByDirectory("uut", "integration-test-dir-se-static", "My own role collection", "Directory Viewer", "cis-central!b13", "Directory_Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My own role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.#", "1"),
					),
				},
				{
					ResourceName:      "btp_directory_role_collection.uut",
					ImportStateIdFunc: getIdForDirectoryRoleCollectionImportId("btp_directory_role_collection.uut", "My own role collection"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - import with resource identity", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection.import_by_resource_identity")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_12_0),
			},
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionNoDescriptionByDirectory("uut", "integration-test-dir-se-static", "My own role collection", "Directory Viewer", "cis-central!b13", "Directory_Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My own role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.#", "1"),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_directory_role_collection.uut", map[string]knownvalue.Check{
							"directory_id": knownvalue.NotNull(),
							"name":         knownvalue.StringExact("My own role collection"),
						}),
					},
				},
				{
					ResourceName:    "btp_directory_role_collection.uut",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})

	t.Run("happy path - with description", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection.with_description")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionWithDescriptionByDirectory("uut", "integration-test-dir-se-static", "My own role collection", "This is my new role collection", "Directory Viewer", "cis-central!b13", "Directory_Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My own role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "description", "This is my new role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.#", "1"),
					),
				},
				{
					ResourceName:      "btp_directory_role_collection.uut",
					ImportStateIdFunc: getIdForDirectoryRoleCollectionImportId("btp_directory_role_collection.uut", "My own role collection"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - multiple roles", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection.multiple_roles")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionByDirectory(
						"uut",
						"integration-test-dir-se-static",
						"My role collection",
						"This is my new role collection",
						directoryRoleCollectionRoleRefTestType{
							Name:              "Directory Viewer",
							RoleTemplateAppId: "cis-central!b13",
							RoleTemplateName:  "Directory_Viewer",
						},
						directoryRoleCollectionRoleRefTestType{
							Name:              "Directory Usage Reporting Viewer",
							RoleTemplateAppId: "uas!b10418",
							RoleTemplateName:  "Directory_Usage_Reporting_Viewer",
						}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "description", "This is my new role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.#", "2"),
					),
				},
				{
					ResourceName:      "btp_directory_role_collection.uut",
					ImportStateIdFunc: getIdForDirectoryRoleCollectionImportId("btp_directory_role_collection.uut", "My role collection"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection.update")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionWithDescriptionByDirectory(
						"uut",
						"integration-test-dir-se-static",
						"My own role collection",
						"This is my new role collection",
						"Directory Viewer",
						"cis-central!b13",
						"Directory_Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My own role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "description", "This is my new role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.#", "1"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.0.name", "Directory Viewer"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionByDirectory(
						"uut",
						"integration-test-dir-se-static",
						"My own role collection",
						"This is my updated role collection",
						directoryRoleCollectionRoleRefTestType{
							Name:              "Directory Viewer",
							RoleTemplateAppId: "cis-central!b13",
							RoleTemplateName:  "Directory_Viewer",
						},
						directoryRoleCollectionRoleRefTestType{
							Name:              "Directory Usage Reporting Viewer",
							RoleTemplateAppId: "uas!b10418",
							RoleTemplateName:  "Directory_Usage_Reporting_Viewer",
						}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My own role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "description", "This is my updated role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.#", "2"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionNoDescriptionByDirectory(
						"uut",
						"integration-test-dir-se-static",
						"My own role collection",
						"Directory Viewer",
						"cis-central!b13",
						"Directory_Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My own role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "description", "This is my updated role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.#", "1"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.0.name", "Directory Viewer"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionWithDescriptionByDirectory(
						"uut",
						"integration-test-dir-se-static",
						"My own role collection",
						"",
						"Directory Viewer",
						"cis-central!b13",
						"Directory_Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My own role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.#", "1"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.0.name", "Directory Viewer"),
					),
				},
				{
					ResourceName:      "btp_directory_role_collection.uut",
					ImportStateIdFunc: getIdForDirectoryRoleCollectionImportId("btp_directory_role_collection.uut", "My own role collection"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("error path - import fails", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection.error_import")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionNoDescriptionByDirectory("uut", "integration-test-dir-se-static", "My special role collection", "Directory Viewer", "cis-central!b13", "Directory_Viewer"),
				},
				{
					ResourceName:      "btp_directory_role_collection.uut",
					ImportStateIdFunc: getIdForDirectoryRoleCollectionImportId("btp_directory_role_collection.uut", ""),
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Expected import identifier with format: directory_id, name. Got:`),
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
					Config:      `resource "btp_directory_role_collection" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "directory_id" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - name mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_directory_role_collection" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - roles mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_directory_role_collection" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "roles" is required, but no definition was found.`),
				},
			},
		})
	})
}

func createRoles(roleName string, roleTemplateAppId string, roleTemplateName string) string {
	roles := []directoryRoleCollectionRoleRefTestType{}
	roles = append(roles, directoryRoleCollectionRoleRefTestType{
		Name:              roleName,
		RoleTemplateAppId: roleTemplateAppId,
		RoleTemplateName:  roleTemplateName,
	})
	rolesJson, _ := json.Marshal(roles)
	return string(rolesJson)
}

func createEmptyRolesIfNil(roles []directoryRoleCollectionRoleRefTestType) string {
	if roles == nil {
		roles = []directoryRoleCollectionRoleRefTestType{}
	}
	rolesJson, _ := json.Marshal(roles)
	return string(rolesJson)
}

func hclResourceDirectoryRoleCollectionNoDescriptionByDirectory(resourceName string, directoryName string, roleCollectionName string, roleName string, roleTemplateAppId string, roleTemplateName string) string {
	return fmt.Sprintf(`
	data "btp_directories" "all" {}
	resource "btp_directory_role_collection" "%s" {
        directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
        name         = "%s"
        roles  		 = %v
    }`, resourceName, directoryName, roleCollectionName, createRoles(roleName, roleTemplateAppId, roleTemplateName))
}

func hclResourceDirectoryRoleCollectionWithDescriptionByDirectory(resourceName string, directoryName string, roleCollectionName string, roleCollectionDescription string, roleName string, roleTemplateAppId string, roleTemplateName string) string {
	return fmt.Sprintf(`
	data "btp_directories" "all" {}
	resource "btp_directory_role_collection" "%s" {
        directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
        name         = "%s"
		description  = "%s"
        roles  		 = %v
    }`, resourceName, directoryName, roleCollectionName, roleCollectionDescription, createRoles(roleName, roleTemplateAppId, roleTemplateName))
}

func hclResourceDirectoryRoleCollectionByDirectory(resourceName string, directoryName string, roleCollectionName string, roleCollectionDescription string, roles ...directoryRoleCollectionRoleRefTestType) string {
	return fmt.Sprintf(`
	data "btp_directories" "all" {}
	resource "btp_directory_role_collection" "%s" {
        directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
        name         = "%s"
		description  = "%s"
        roles  		 = %v
    }`, resourceName, directoryName, roleCollectionName, roleCollectionDescription, createEmptyRolesIfNil(roles))
}

func getIdForDirectoryRoleCollectionImportId(resourceName string, name string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		if name != "" {
			return fmt.Sprintf("%s,%s", rs.Primary.Attributes["directory_id"], name), nil
		}
		return rs.Primary.Attributes["directory_id"], nil
	}
}
