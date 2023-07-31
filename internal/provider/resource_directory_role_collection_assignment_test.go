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
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignment("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "Directory Viewer", "jenny.doe@test.com"),
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
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment_with_origin")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentWithOrigin("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "Directory Viewer", "john.doe@test.com", "terraformint-platform"),
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

	t.Run("error path - role collection import fails", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment_import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignment("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "Directory Viewer", "jenny.doe@test.com"),
				},
				{
					ResourceName:      "btp_directory_role_collection_assignment.uut",
					ImportStateId:     "05368777-4934-41e8-9f3c-6ec5f4d564b9",
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

func hclResourceDirectoryRoleCollectionAssignment(resourceName string, directoryId string, roleCollectionName string, userName string) string {

	return fmt.Sprintf(`
resource "btp_directory_role_collection_assignment" "%s"{
    directory_id        = "%s"
	role_collection_name = "%s"
	user_name            = "%s"
}`, resourceName, directoryId, roleCollectionName, userName)
}

func hclResourceDirectoryRoleCollectionAssignmentWithOrigin(resourceName string, directoryId string, roleCollectionName string, userName string, origin string) string {

	return fmt.Sprintf(`
resource "btp_directory_role_collection_assignment" "%s"{
    directory_id        = "%s"
	role_collection_name = "%s"
	user_name            = "%s"
	origin               = "%s"
}`, resourceName, directoryId, roleCollectionName, userName, origin)
}
