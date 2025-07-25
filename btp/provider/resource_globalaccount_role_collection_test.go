package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Needed for JSON mapping - fails with data types of globalaccountRoleCollectionRoleRef struc
type globalaccountRoleCollectionRoleRefTestType struct {
	Name              string `json:"name"`
	RoleTemplateAppId string `json:"role_template_app_id"`
	RoleTemplateName  string `json:"role_template_name"`
}

func TestResourceGlobalaccountRoleCollection(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role_collection")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollection(
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

	t.Run("happy path - import with resource identity", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role_collection.import_by_resource_identity")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_12_0),
			},
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollection(
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
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_globalaccount_role_collection.uut", map[string]knownvalue.Check{
							"name": knownvalue.StringExact("My new role collection"),
						}),
					},
				},
				{
					ResourceName:    "btp_globalaccount_role_collection.uut",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})

	t.Run("happy path - update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role_collection.update")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollection(
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
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollection(
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

	t.Run("happy path - update removing description", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role_collection.update_rm_desc")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollection(
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
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollection(
						"uut",
						"My new role collection",
						"",
						globalaccountRoleCollectionRoleRefTestType{
							Name:              "System Landscape Viewer",
							RoleTemplateAppId: "cmp!b17875",
							RoleTemplateName:  "GlobalAccount_System_Landscape_Viewer",
						}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "name", "My new role collection"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "description", ""),
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

	t.Run("happy path - update omitting description", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role_collection.update_wo_desc")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollection(
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
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollectionWoDescription(
						"uut",
						"My new role collection",
						globalaccountRoleCollectionRoleRefTestType{
							Name:              "System Landscape Viewer",
							RoleTemplateAppId: "cmp!b17875",
							RoleTemplateName:  "GlobalAccount_System_Landscape_Viewer",
						}),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "name", "My new role collection"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection.uut", "description", "Description of my new role collection"),
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
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role_collection.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollection("uut", "My new role collection", "Description of my new role collection"),
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

func hclResourceGlobalaccountRoleCollection(resourceName string, displayName string, description string, roles ...globalaccountRoleCollectionRoleRefTestType) string {
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

func hclResourceGlobalaccountRoleCollectionWoDescription(resourceName string, displayName string, roles ...globalaccountRoleCollectionRoleRefTestType) string {
	if roles == nil {
		roles = []globalaccountRoleCollectionRoleRefTestType{}
	}
	rolesJson, _ := json.Marshal(roles)

	return fmt.Sprintf(`resource "btp_globalaccount_role_collection" "%s" {
        name         = "%s"
        roles        = %v
    }`, resourceName, displayName, string(rolesJson))
}
