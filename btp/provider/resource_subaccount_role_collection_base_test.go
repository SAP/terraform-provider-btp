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

func TestResourceSubaccountRoleCollectionBase(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_base")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubAccountRoleCollectionBaseBySubaccount(
						"uut",
						"integration-test-acc-static",
						"My new role collection base",
						"Description of my new role collection base"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "name", "My new role collection base"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "description", "Description of my new role collection base"),
					),
				},
				{
					ResourceName:      "btp_subaccount_role_collection_base.uut",
					ImportStateIdFunc: getImportIdForRoleCollection("btp_subaccount_role_collection_base.uut", "My new role collection base"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - import with resource identity ", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_base.import_resource_identity")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_12_0), // ImportBlockWithResourceIdentity requires Terraform 1.12.0 or later
			},
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubAccountRoleCollectionBaseBySubaccount(
						"uut",
						"integration-test-acc-static",
						"My new role collection base",
						"Description of my new role collection base"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "name", "My new role collection base"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "description", "Description of my new role collection base"),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_subaccount_role_collection_base.uut", map[string]knownvalue.Check{
							"subaccount_id": knownvalue.NotNull(),
							"name":          knownvalue.StringExact("My new role collection base"),
						}),
					},
				},
				{
					ResourceName:    "btp_subaccount_role_collection_base.uut",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})

	t.Run("happy path - update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_base.update")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubAccountRoleCollectionBaseBySubaccount(
						"uut",
						"integration-test-acc-static",
						"My new role collection base",
						"Updated description of my new role collection base"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "name", "My new role collection base"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "description", "Updated description of my new role collection base"),
					),
				},
				{
					ResourceName:      "btp_subaccount_role_collection_base.uut",
					ImportStateIdFunc: getImportIdForRoleCollectionBase("btp_subaccount_role_collection_base.uut", "My new role collection base"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - update removing description", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_base.update_rm_desc")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubAccountRoleCollectionBaseBySubaccount(
						"uut",
						"integration-test-acc-static",
						"My new role collection base",
						"Description of my new role collection base"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "name", "My new role collection base"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "description", "Description of my new role collection base"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubAccountRoleCollectionBaseBySubaccount(
						"uut",
						"integration-test-acc-static",
						"My new role collection base",
						""),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "name", "My new role collection base"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "description", ""),
					),
				},
				{
					ResourceName:      "btp_subaccount_role_collection_base.uut",
					ImportStateIdFunc: getImportIdForRoleCollectionBase("btp_subaccount_role_collection_base.uut", "My new role collection base"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - update omitting description", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_base.update_wo_desc")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubAccountRoleCollectionBaseBySubaccount(
						"uut",
						"integration-test-acc-static",
						"My new role collection base",
						"Description of my new role collection base"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "name", "My new role collection base"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "description", "Description of my new role collection base"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubAccountRoleCollectionBaseWoDescriptionBySubaccount(
						"uut",
						"integration-test-acc-static",
						"My new role collection base"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "name", "My new role collection base"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "description", "Description of my new role collection base"),
					),
				},
				{
					ResourceName:      "btp_subaccount_role_collection_base.uut",
					ImportStateIdFunc: getImportIdForRoleCollectionBase("btp_subaccount_role_collection_base.uut", "My new role collection base"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("error path - import with wrong key", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_base.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubAccountRoleCollectionBaseBySubaccount("uut", "integration-test-acc-static",
						"My new role collection base", "Description of my new role collection base"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "name", "My new role collection base"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_base.uut", "description", "Description of my new role collection base"),
					),
				},
				{
					ResourceName:      "btp_subaccount_role_collection_base.uut",
					ImportStateId:     "00000000-0000-0000-0000-000000000000",
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Expected import identifier with format: subaccount_id, name. Got:`),
				},
			},
		})
	})

	t.Run("error path - subacount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubAccountRoleCollectionBaseNoSubaccountId("uut", "My new role collection base", "Description of my new role collection base"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

}

func hclResourceSubAccountRoleCollectionBaseBySubaccount(resourceName string, subaccountName string, displayName string, description string) string {
	return fmt.Sprintf(`
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_role_collection_base" "%s" {
        subaccount_id       = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
		name      			= "%s"
        description      	= "%s"
    }`, resourceName, subaccountName, displayName, description)
}

func hclResourceSubAccountRoleCollectionBaseNoSubaccountId(resourceName string, displayName string, description string) string {
	return fmt.Sprintf(`resource "btp_subaccount_role_collection_base" "%s" {
        name      			= "%s"
        description      	= "%s"
    }`, resourceName, displayName, description)
}

func hclResourceSubAccountRoleCollectionBaseWoDescriptionBySubaccount(resourceName string, subaccountName string, displayName string) string {
	return fmt.Sprintf(`
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_role_collection_base" "%s" {
        subaccount_id       = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
		name      			= "%s"
    }`, resourceName, subaccountName, displayName)
}

func getImportIdForRoleCollectionBase(resourceName string, roleCollectionDisplayName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["subaccount_id"], roleCollectionDisplayName), nil
	}
}
