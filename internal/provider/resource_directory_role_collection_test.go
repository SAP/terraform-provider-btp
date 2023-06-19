package provider

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Needed for JSON mapping - fails with data types of globalaccountRoleCollectionRoleRef struc
type directoryRoleCollectionRoleRefTestType struct {
	Name              string `json:"name"`
	RoleTemplateAppId string `json:"role_template_app_id"`
	RoleTemplateName  string `json:"role_template_name"`
}

func TestResourceDirectoryRoleCollection(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_directory_role_collection.no_description")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceDirectoryRoleCollectionNoDescription("uut", "05368777-4934-41e8-9f3c-6ec5f4d564b9", "My own role collection", "Directory Viewer", "cis-central!b13", "Directory_Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection.uut", "directory_id", regexpValidUUID),
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
