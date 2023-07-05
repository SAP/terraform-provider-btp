package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Needed for JSON mapping - fails with data types of globalaccountRoleCollectionRoleRef struc
type globalaccountRoleCollectionRoleRefTestType struct {
	Name              string `json:"name"`
	RoleTemplateAppId string `json:"role_template_app_id"`
	RoleTemplateName  string `json:"role_template_name"`
}

func TestResourceGlobalAccountRoleCollection(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_globalaccount_role_collection")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceGlobalAccountRoleCollection("uut", "My new role collection", "Description of my new role collection", "Global Account Viewer", "cis-central!b13", "GlobalAccount_Viewer"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "name", "My new role collection"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "description", "Description of my new role collection"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "roles.#", "1"),
					),
				},
				{
					ResourceName:      "btp_globalaccount_role_collection.uut",
					ImportStateId:     "My new role collection",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("error path - import fails", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_globalaccount_role_collection_import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceGlobalAccountRoleCollection("uut", "My new role collection", "Description of my new role collection", "Global Account Viewer", "cis-central!b13", "GlobalAccount_Viewer"),
				},
				{
					ResourceName:      "btp_globalaccount_role_collection.uut",
					ImportStateId:     "ef23ace8-6ade-4d78-9c1f-8df729548bbf,My new role collection",
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Expected import identifier with format: name. Got:`),
				},
			},
		})
	})
}

func hclResourceGlobalAccountRoleCollection(resourceName string, displayName string, description string, roleName string, RoleTemplateAppId string, RoleTemplateName string) string {

	roles := []globalaccountRoleCollectionRoleRefTestType{}

	roles = append(roles, globalaccountRoleCollectionRoleRefTestType{
		Name:              roleName,
		RoleTemplateAppId: RoleTemplateAppId,
		RoleTemplateName:  RoleTemplateName,
	})
	rolesJson, _ := json.Marshal(roles)

	return fmt.Sprintf(`resource "btp_globalaccount_role_collection" "%s" {
        name         = "%s"
        description  = "%s"
        roles        = %v
    }`, resourceName, displayName, description, string(rolesJson))
}
