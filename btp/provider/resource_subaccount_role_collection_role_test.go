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

func TestResourceSubaccountRoleCollectionRole(t *testing.T) {
	// 1. HAPPY PATH: Create Base + Role, then Import via 5-part String
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_role")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					// We must create the Base collection first in the same config
					Config: hclProviderFor(user) +
						hclResourceSubAccountRoleCollectionBaseBySubaccount(
							"base",
							"integration-test-acc-static",
							"MyTestCollection",
							"Desc") +
						hclResourceSubAccountRoleCollectionRoleBySubaccount(
							"uut",
							"integration-test-acc-static",
							"MyTestCollection",
							"Subaccount Admin",
							"Subaccount_Admin",
							"cis-local!b2"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_role.uut", "role_name", "Subaccount Admin"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_role.uut", "role_template_name", "Subaccount_Admin"),
					),
				},
				{
					// Test the 5-part String Import: sub_id,name,role_name,app_id,temp_name
					ResourceName:      "btp_subaccount_role_collection_role.uut",
					ImportStateIdFunc: getImportIdForRoleCollectionRole("btp_subaccount_role_collection_role.uut", "MyTestCollection", "Subaccount Admin", "Subaccount_Admin", "cis-local!b2"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	// 2. IDENTITY IMPORT: Test modern Terraform 1.12+ Identity Import
	t.Run("happy path - import with resource identity", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_role.import_identity")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_12_0),
			},
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) +
						hclResourceSubAccountRoleCollectionBaseBySubaccount("base", "integration-test-acc-static", "MyTestCollection", "Desc") +
						hclResourceSubAccountRoleCollectionRoleBySubaccount("uut", "integration-test-acc-static", "MyTestCollection", "Subaccount Admin", "Subaccount_Admin", "cis-local!b2"),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_subaccount_role_collection_role.uut", map[string]knownvalue.Check{
							"subaccount_id":        knownvalue.NotNull(),
							"name":                 knownvalue.StringExact("MyTestCollection"),
							"role_name":            knownvalue.StringExact("Subaccount Admin"),
							"role_template_name":   knownvalue.StringExact("Subaccount_Admin"),
							"role_template_app_id": knownvalue.StringExact("cis-local!b2"),
						}),
					},
				},
				{
					ResourceName:    "btp_subaccount_role_collection_role.uut",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})

	// 3. ERROR PATH: Test invalid ID format (too few parts)
	t.Run("error path - import with wrong key", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_role.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) +
						hclResourceSubAccountRoleCollectionBaseBySubaccount("base", "integration-test-acc-static", "MyTestCollection", "Desc") +
						hclResourceSubAccountRoleCollectionRoleBySubaccount("uut", "integration-test-acc-static", "MyTestCollection", "Subaccount Admin", "Subaccount_Admin", "cis-local!b2"),
				},
				{
					ResourceName:      "btp_subaccount_role_collection_role.uut",
					ImportStateId:     "too,short", // Only 2 parts instead of 5
					ImportState:       true,
					ImportStateVerify: false,
					ExpectError:       regexp.MustCompile(`Expected: subaccount_id,collection_name,role_name,app_id,template_name`),
				},
			},
		})
	})

	// 4. ERROR PATH: Missing required subaccount_id
	t.Run("error path - subaccount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubAccountRoleCollectionRoleNoSubaccountId("uut", "CollName", "RoleName", "TempName", "AppId"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required`),
				},
			},
		})
	})

}

func hclResourceSubAccountRoleCollectionRoleBySubaccount(resourceName, subaccountName, displayName, roleName, roleTemplateName, roleTemplateAppId string) string {
	return fmt.Sprintf(`
	resource "btp_subaccount_role_collection_role" "%s" {
        subaccount_id       = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
		name      			= "%s"
        role_name            = "%s"
        role_template_name   = "%s"
        role_template_app_id = "%s"
		depends_on = [btp_subaccount_role_collection_base.base]
    }`, resourceName, subaccountName, displayName, roleName, roleTemplateName, roleTemplateAppId)
}

func hclResourceSubAccountRoleCollectionRoleNoSubaccountId(resourceName, displayName, roleName, roleTemplateName, roleTemplateAppId string) string {
	return fmt.Sprintf(`resource "btp_subaccount_role_collection_role" "%s" {
        name      			= "%s"
		role_name            = "%s"
        role_template_name   = "%s"
        role_template_app_id = "%s"
		depends_on = [btp_subaccount_role_collection_base.base]
    }`, resourceName, displayName, roleName, roleTemplateName, roleTemplateAppId)
}

func getImportIdForRoleCollectionRole(resourceName, roleCollectionDisplayName, roleName, roleTemplateName, roleTemplateAppId string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s,%s,%s,%s", rs.Primary.Attributes["subaccount_id"], roleCollectionDisplayName, roleName, roleTemplateAppId, roleTemplateName), nil
	}
}
