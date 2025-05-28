package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceSubAccountRole(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountRole(
						"uut",
						"integration-test-acc-static",
						"Subaccount Viewer Test",
						"Subaccount_Viewer",
						"cis-local!b2",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "name", "Subaccount Viewer Test"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "role_template_name", "Subaccount_Viewer"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "app_id", "cis-local!b2"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "read_only", "false"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "scopes.#", "0"),
					),
				},
				{
					ResourceName:      "btp_subaccount_role.uut",
					ImportStateIdFunc: getIdForSubaccountRoleImportId("btp_subaccount_role.uut", "Subaccount Viewer Test", "Subaccount_Viewer", "cis-local!b2"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - with attributes", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_with_attributes")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountRoleWithAttributes(
						"uut",
						"integration-test-acc-static",
						"Custom Role with Attributes",
						"Application_Frontend_Developer",
						"eu12-appfront!b390135",
						"namespace",
						"saml",
						"custom",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "name", "Custom Role with Attributes"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "role_template_name", "Application_Frontend_Developer"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "app_id", "eu12-appfront!b390135"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "read_only", "false"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "attribute_list.#", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "attribute_list.0.attribute_name", "namespace"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "attribute_list.0.attribute_value_origin", "saml"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "attribute_list.0.attribute_values.#", "1"),
						resource.TestCheckResourceAttr("btp_subaccount_role.uut", "attribute_list.0.attribute_values.0", "custom"),
					),
				},
				{
					ResourceName:      "btp_subaccount_role.uut",
					ImportStateIdFunc: getIdForSubaccountRoleImportId("btp_subaccount_role.uut", "Custom Role with Attributes", "Application_Frontend_Developer", "eu12-appfront!b390135"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("error path - subaccount_id, name, role_template_name and app_id are mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_subaccount_role" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "(subaccount_id|name|role_template_name|app_id)" is required, but no definition was found.`),
				},
			},
		})
	})
	t.Run("error path - subaccount_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountRoleBySubaccountId("uut", "this-is-not-a-uuid", "a", "b", "c"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
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
					Config:      hclResourceSubaccountRole("uut", "integration-test-acc-static", "", "b", "c"),
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
					Config:      hclDatasourceSubaccountRole("uut", "integration-test-acc-static", "a", "", "c"),
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
					Config:      hclDatasourceSubaccountRole("uut", "integration-test-acc-static", "a", "b", ""),
					ExpectError: regexp.MustCompile(`Attribute app_id string length must be at least 1, got: 0`),
				},
			},
		})
	})

	t.Run("error path - import fails", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role.error_import")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountRole("uut", "integration-test-acc-static", "Subaccount Viewer Test", "Subaccount_Viewer", "cis-local!b2"),
				},
				{
					ResourceName:      "btp_subaccount_role.uut",
					ImportStateIdFunc: getSubaccountRoleImportIdNoAppIdNoRoleTemplateName("btp_subaccount_role.uut"),
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Unexpected Import Identifier`),
				},
			},
		})
	})

}

func hclResourceSubaccountRole(resourceName string, subaccountName string, name string, roleTemplateName string, appId string) string {
	return fmt.Sprintf(`
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_role" "%s" {
        subaccount_id       = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
		name      			= "%s"
		role_template_name  = "%s"
		app_id              = "%s"
    }`, resourceName, subaccountName, name, roleTemplateName, appId)
}

func hclResourceSubaccountRoleBySubaccountId(resourceName string, subaccountId string, name string, roleTemplateName string, appId string) string {
	template := `
resource "btp_subaccount_role" "%s" {
    subaccount_id       = "%s"
    name                = "%s"
    role_template_name  = "%s"
    app_id              = "%s"
}`

	return fmt.Sprintf(template, resourceName, subaccountId, name, roleTemplateName, appId)
}

func hclResourceSubaccountRoleWithAttributes(resourceName string, subaccountName string, name string, roleTemplateName string, appId string, attributeName string, attributeValueOrigin string, attributeValues string) string {
	template := `
data "btp_subaccounts" "all" {}
resource "btp_subaccount_role" "%s" {
    subaccount_id       = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
    name                = "%s"
    role_template_name  = "%s"
    app_id              = "%s"
		attribute_list = [
    {
      attribute_name         = "%s"
      attribute_value_origin = "%s"
      attribute_values       = ["%s"]
    },
  ]
}`

	return fmt.Sprintf(template, resourceName, subaccountName, name, roleTemplateName, appId, attributeName, attributeValueOrigin, attributeValues)
}

func getSubaccountRoleImportIdNoAppIdNoRoleTemplateName(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.Attributes["subaccount_id"], nil
	}
}

func getIdForSubaccountRoleImportId(resourceName string, name string, role_template_name string, app_id string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s,%s,%s,%s", rs.Primary.Attributes["subaccount_id"], name, role_template_name, app_id), nil
	}
}
