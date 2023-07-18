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
					Config: hclProvider() + hclResourceGlobalAccountRoleCollection(
						"uut",
						"My new role collection",
						"Description of my new role collection",
						globalaccountRoleCollectionRoleRefTestType{
							Name:              "Global Account Viewer",
							RoleTemplateAppId: "cis-central!b13",
							RoleTemplateName:  "GlobalAccount_Viewer",
						}),
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

	t.Run("happy path - update", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_globalaccount_role_collection.update")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceGlobalAccountRoleCollection(
						"uut",
						"My new role collection",
						"Description of my new role collection",
						globalaccountRoleCollectionRoleRefTestType{
							Name:              "Global Account Viewer",
							RoleTemplateAppId: "cis-central!b13",
							RoleTemplateName:  "GlobalAccount_Viewer",
						}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "name", "My new role collection"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "description", "Description of my new role collection"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "roles.#", "1"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "roles.0.name", "Global Account Viewer"),
					),
				},
				{
					Config: hclProvider() + hclResourceGlobalAccountRoleCollection(
						"uut",
						"My new role collection",
						"Description of my updated role collection",
						globalaccountRoleCollectionRoleRefTestType{
							Name:              "System Landscape Viewer",
							RoleTemplateAppId: "cmp!b17875",
							RoleTemplateName:  "GlobalAccount_System_Landscape_Viewer",
						}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "name", "My new role collection"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "description", "Description of my updated role collection"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "roles.#", "1"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "roles.0.name", "System Landscape Viewer"),
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
		rec := setupVCR(t, "fixtures/resource_globalaccount_role_collection.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceGlobalAccountRoleCollection("uut", "My new role collection", "Description of my new role collection"),
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

func hclResourceGlobalAccountRoleCollection(resourceName string, displayName string, description string, roles ...globalaccountRoleCollectionRoleRefTestType) string {
	if roles == nil {
		roles = []globalaccountRoleCollectionRoleRefTestType{}
	}
	rolesJson, _ := json.Marshal(roles)

	return fmt.Sprintf(`resource "btp_globalaccount_role_collection" "%s" {
        name         = "%s"
        description  = "%s"
        roles        = %v
    }`, resourceName, displayName, description, string(rolesJson))
}
