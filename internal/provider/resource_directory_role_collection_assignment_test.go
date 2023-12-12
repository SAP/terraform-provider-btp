package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceDirectoryRoleCollectionAssignment(t *testing.T) {
	t.Parallel()
	t.Run("happy path - simple role collection assignment", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "jenny.doe@test.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection_assignment.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "role_collection_name", "Directory Viewer"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "origin", "ldap"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection assignment with origin", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.with_origin")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentWithOriginByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "john.doe@test.com", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection_assignment.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "role_collection_name", "Directory Viewer"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "user_name", "john.doe@test.com"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "origin", "terraformint-platform"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection assignment with origin and group", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.with_origin_and_group")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentWithOriginAndGroupByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "tf-test-group", "terraformint-platform"),
					// We do not get back any information about the group, so if the call succeeds we assume that the asssignment/unassignment worked
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection_assignment.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "role_collection_name", "Directory Viewer"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "group_name", "tf-test-group"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "origin", "terraformint-platform"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection assignment with origin and attribute", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.with_origin_and_attribute")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentWithOriginAndAttributeByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "tf_attr_name_test", "tf_attr_val_test", "terraformint-platform"),
					// We do not get back any information about the group, so if the call succeeds we assume that the asssignment/unassignment worked
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection_assignment.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "role_collection_name", "Directory Viewer"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "attribute_name", "tf_attr_name_test"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "attribute_value", "tf_attr_val_test"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "origin", "terraformint-platform"),
					),
				},
			},
		})
	})

	t.Run("error path - role collection import fails", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "jenny.doe@test.com"),
				},
				{
					ResourceName:      "btp_directory_role_collection_assignment.uut",
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Import Not Supported`),
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
					Config:      `resource "btp_directory_role_collection_assignment" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "directory_id" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - role_collection_name mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_directory_role_collection_assignment" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "role_collection_name" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclResourceDirectoryRoleCollectionAssignmentByDirectory(resourceName string, directoryName string, roleCollectionName string, userName string) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_role_collection_assignment" "%s"{
    directory_id        = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
	role_collection_name = "%s"
	user_name            = "%s"
}`, resourceName, directoryName, roleCollectionName, userName)
}

func hclResourceDirectoryRoleCollectionAssignmentWithOriginByDirectory(resourceName string, directoryName string, roleCollectionName string, userName string, origin string) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_role_collection_assignment" "%s"{
    directory_id        = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
	role_collection_name = "%s"
	user_name            = "%s"
	origin               = "%s"
}`, resourceName, directoryName, roleCollectionName, userName, origin)
}

func hclResourceDirectoryRoleCollectionAssignmentWithOriginAndGroupByDirectory(resourceName string, directoryName string, roleCollectionName string, groupName string, origin string) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_role_collection_assignment" "%s"{
    directory_id        = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
	role_collection_name = "%s"
	origin               = "%s"
	group_name           = "%s"
}`, resourceName, directoryName, roleCollectionName, origin, groupName)
}

func hclResourceDirectoryRoleCollectionAssignmentWithOriginAndAttributeByDirectory(resourceName string, directoryName string, roleCollectionName string, attributeName string, attributeValue string, origin string) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_role_collection_assignment" "%s"{
    directory_id        = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
	role_collection_name = "%s"
	origin               = "%s"
	attribute_name       = "%s"
	attribute_value      = "%s"
}`, resourceName, directoryName, roleCollectionName, origin, attributeName, attributeValue)
}
