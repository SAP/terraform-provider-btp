package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceDirectoryRoleCollection(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_directory_role_collection")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceDirectoryRoleCollection("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "My role collection", "Description of my role collection"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),

						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "name", "My role collection"),
						resource.TestCheckResourceAttr("btp_directory_role_collection.uut", "description", "Description of my role collection"),
					),
				}, /*
					{
						ResourceName:      "btp_directory_role_collection.uut",
						ImportState:       true,
						ImportStateVerify: true,
					}*/
			},
		})
	})
}

func hclResourceDirectoryRoleCollection(resourceName string, directoryId string, displayName string, description string) string {
	return fmt.Sprintf(`resource "btp_directory_role_collection" "%s" {
        directory_id = "%s"
        name         = "%s"
        description  = "%s"
    }`, resourceName, directoryId, displayName, description)
}
