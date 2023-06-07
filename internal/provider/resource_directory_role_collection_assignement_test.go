package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceDirectoryRoleCollectionAssignment(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.user")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceDirectoryRoleCollectionAssignmentUser("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "Global Account Viewer", "jenny.doe@test.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection_assignment.uut", "directory_id", regexpValidUUID),

						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "role_collection_name", "My role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "user_name", "Description of my role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "group_name", "Description of my role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "origin", "Description of my role collection"),
					),
				}, /*
					{
						ResourceName:      "btp_directory_role_collection_assignment.uut",
						ImportState:       true,
						ImportStateVerify: true,
					}*/
			},
		})
	})
}

func hclResourceDirectoryRoleCollectionAssignmentUser(resourceName string, directoryId string, roleCollectionName string, userName string) string {
	return fmt.Sprintf(`resource "btp_directory_role_collection_assignment" "%s" {
        directory_id 			= "%s"
        role_collection_name    = "%s"
        user_name      			= "%s"
    }`, resourceName, directoryId, roleCollectionName, userName)
}
