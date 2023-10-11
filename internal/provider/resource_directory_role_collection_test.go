package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionNoDescription("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "My own role collection", "Directory Viewer", "cis-central!b13", "Directory_Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My own role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.#", "1"),
					),
				},
				{
					ResourceName:      "btp_directory_role_collection.uut",
					ImportStateId:     "05368777-4934-41e8-9f3c-6ec5f4d564b9,My own role collection",
					ImportState:       true,
					ImportStateVerify: true,
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
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionWithDescription("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "My own role collection", "This is my new role collection", "Directory Viewer", "cis-central!b13", "Directory_Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My own role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "description", "This is my new role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "roles.#", "1"),
					),
				},
				{
					ResourceName:      "btp_directory_role_collection.uut",
					ImportStateId:     "05368777-4934-41e8-9f3c-6ec5f4d564b9,My own role collection",
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
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollection(
						"uut",
						"05368777-4934-41e8-9f3c-6ec5f4d564b9",
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
					ImportStateId:     "05368777-4934-41e8-9f3c-6ec5f4d564b9,My role collection",
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
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionWithDescription(
						"uut",
						"05368777-4934-41e8-9f3c-6ec5f4d564b9",
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
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollection(
						"uut",
						"05368777-4934-41e8-9f3c-6ec5f4d564b9",
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
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionNoDescription(
						"uut",
						"05368777-4934-41e8-9f3c-6ec5f4d564b9",
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
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionWithDescription(
						"uut",
						"05368777-4934-41e8-9f3c-6ec5f4d564b9",
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
					ImportStateId:     "05368777-4934-41e8-9f3c-6ec5f4d564b9,My own role collection",
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
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionNoDescription("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "My special role collection", "Directory Viewer", "cis-central!b13", "Directory_Viewer"),
				},
				{
					ResourceName:      "btp_directory_role_collection.uut",
					ImportStateId:     "05368777-4934-41e8-9f3c-6ec5f4d564b9",
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

func hclResourceDirectoryRoleCollectionNoDescription(resourceName string, directoryId string, roleCollectionName string, roleName string, RoleTemplateAppId string, RoleTemplateName string) string {
	roles := []directoryRoleCollectionRoleRefTestType{}

	roles = append(roles, directoryRoleCollectionRoleRefTestType{
		Name:              roleName,
		RoleTemplateAppId: RoleTemplateAppId,
		RoleTemplateName:  RoleTemplateName,
	})
	rolesJson, _ := json.Marshal(roles)

	return fmt.Sprintf(`resource "btp_directory_role_collection" "%s" {
        directory_id = "%s"
        name         = "%s"
        roles  		 = %v
    }`, resourceName, directoryId, roleCollectionName, string(rolesJson))
}

func hclResourceDirectoryRoleCollectionWithDescription(resourceName string, directoryId string, roleCollectionName string, roleCollectionDescription string, roleName string, RoleTemplateAppId string, RoleTemplateName string) string {
	roles := []directoryRoleCollectionRoleRefTestType{}

	roles = append(roles, directoryRoleCollectionRoleRefTestType{
		Name:              roleName,
		RoleTemplateAppId: RoleTemplateAppId,
		RoleTemplateName:  RoleTemplateName,
	})
	rolesJson, _ := json.Marshal(roles)

	return fmt.Sprintf(`resource "btp_directory_role_collection" "%s" {
        directory_id = "%s"
        name         = "%s"
		description  = "%s"
        roles  		 = %v
    }`, resourceName, directoryId, roleCollectionName, roleCollectionDescription, string(rolesJson))
}

func hclResourceDirectoryRoleCollection(resourceName string, directoryId string, roleCollectionName string, roleCollectionDescription string, roles ...directoryRoleCollectionRoleRefTestType) string {
	if roles == nil {
		roles = []directoryRoleCollectionRoleRefTestType{}
	}
	rolesJson, _ := json.Marshal(roles)

	return fmt.Sprintf(`resource "btp_directory_role_collection" "%s" {
        directory_id = "%s"
        name         = "%s"
		description  = "%s"
        roles  		 = %v
    }`, resourceName, directoryId, roleCollectionName, roleCollectionDescription, string(rolesJson))
}
