package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
}

func hclResourceGlobalAccountRole(resourceName string, name string, roleTemplateName string, appId string) string {
	return fmt.Sprintf(`
	resource "btp_globalaccount_role" "%s" {
		name      			= "%s"
		role_template_name  = "%s"
		app_id              = "%s"
    }`, resourceName, name, roleTemplateName, appId)
}
