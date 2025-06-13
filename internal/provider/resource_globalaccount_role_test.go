package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestResourceGlobalAccountRole(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalAccountRole(
						"uut",
						"GlobalAccount Viewer Test",
						"GlobalAccount_Viewer",
						"cis-central!b13",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "name", "GlobalAccount Viewer Test"),
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "role_template_name", "GlobalAccount_Viewer"),
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "app_id", "cis-central!b13"),
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "read_only", "false"),
					),
				},
				{
					ResourceName:      "btp_globalaccount_role.uut",
					ImportStateIdFunc: getIdForGlobalAccountRoleImportId("btp_globalaccount_role.uut", "GlobalAccount Viewer Test", "GlobalAccount_Viewer", "cis-central!b13"),
					ImportState:       true,
				},
			},
		})
	})

	t.Run("happy path - import with resource identity", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role.import_by_resource_identity")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_12_0),
			},
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalAccountRole(
						"uut",
						"GlobalAccount Viewer Test",
						"GlobalAccount_Viewer",
						"cis-central!b13",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "name", "GlobalAccount Viewer Test"),
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "role_template_name", "GlobalAccount_Viewer"),
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "app_id", "cis-central!b13"),
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "read_only", "false"),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_globalaccount_role.uut", map[string]knownvalue.Check{
							"name":               knownvalue.StringExact("GlobalAccount Viewer Test"),
							"role_template_name": knownvalue.StringExact("GlobalAccount_Viewer"),
							"app_id":             knownvalue.StringExact("cis-central!b13"),
						}),
					},
				},
				{
					ResourceName:    "btp_globalaccount_role.uut",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})

	t.Run("error path - name, role_template_name and app_id are mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_globalaccount_role" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "(name|role_template_name|app_id)" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - name must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceGlobalAccountRole("uut", "", "b", "c"),
					ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1, got: 0`),
				},
			},
		})
	})
	t.Run("error path - role_template_name must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceGlobalAccountRole("uut", "a", "", "c"),
					ExpectError: regexp.MustCompile(`Attribute role_template_name string length must be at least 1, got: 0`),
				},
			},
		})
	})
	t.Run("error path - app_id must not be empty", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceGlobalAccountRole("uut", "a", "b", ""),
					ExpectError: regexp.MustCompile(`Attribute app_id string length must be at least 1, got: 0`),
				},
			},
		})
	})

	t.Run("error path - update role name", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role.update")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalAccountRole(
						"uut",
						"GlobalAccount Viewer Test",
						"GlobalAccount_Viewer",
						"cis-central!b13",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "name", "GlobalAccount Viewer Test"),
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "role_template_name", "GlobalAccount_Viewer"),
						resource.TestCheckResourceAttr("btp_globalaccount_role.uut", "app_id", "cis-central!b13"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceGlobalAccountRole(
						"uut",
						"GlobalAccount Viewer Test updated",
						"GlobalAccount_Viewer",
						"cis-central!b13",
					),
					ExpectError: regexp.MustCompile(`this resource is not supposed to be updated`),
				},
			},
		})
	})

}

func hclResourceGlobalAccountRole(resourceName string, name string, roleTemplateName string, appId string) string {
	return fmt.Sprintf(`
	resource "btp_globalaccount_role" "%s" {
		name      			= "%s"
		role_template_name  = "%s"
		app_id              = "%s"
    }`, resourceName, name, roleTemplateName, appId)
}

func getIdForGlobalAccountRoleImportId(resourceName string, name string, role_template_name string, app_id string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s,%s,%s", rs.Primary.Attributes["name"], role_template_name, app_id), nil
	}
}
