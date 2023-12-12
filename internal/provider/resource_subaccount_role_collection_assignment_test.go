package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceRolCollectionAssignment(t *testing.T) {
	t.Parallel()
	t.Run("happy path - simple role collection assignment", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_assignment")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceRoleCollectionAssignmentBySubaccount("uut", "integration-test-acc-static", "Destination Administrator", "jenny.doe@test.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_role_collection_assignment.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "role_collection_name", "Destination Administrator"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "origin", "ldap"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection assignment with origin", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_assignment.with_origin")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceRoleCollectionAssignmentWithOriginBySubaccount("uut", "integration-test-acc-static", "Destination Administrator", "john.doe@test.com", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_role_collection_assignment.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "role_collection_name", "Destination Administrator"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "user_name", "john.doe@test.com"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "origin", "terraformint-platform"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection assignment with origin and group", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_assignment.with_origin_and_group")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceRoleCollectionAssignmentWithOriginAndGroupBySubaccount("uut", "integration-test-acc-static", "Destination Administrator", "tf-test-group", "terraformint-platform"),
					// We do not get back any information about the group, so if the call succeeds we assume that the asssignment/unassignment worked
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_role_collection_assignment.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "role_collection_name", "Destination Administrator"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "group_name", "tf-test-group"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "origin", "terraformint-platform"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection assignment with origin and attribute", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_assignment.with_origin_and_attribute")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceRoleCollectionAssignmentWithOriginAndAttributeBySubaccount("uut", "integration-test-acc-static", "Destination Administrator", "tf_attr_name_test", "tf_attr_val_test", "terraformint-platform"),
					// We do not get back any information about the group, so if the call succeeds we assume that the asssignment/unassignment worked
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_role_collection_assignment.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "role_collection_name", "Destination Administrator"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "attribute_name", "tf_attr_name_test"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "attribute_value", "tf_attr_val_test"),
						resource.TestCheckResourceAttr("btp_subaccount_role_collection_assignment.uut", "origin", "terraformint-platform"),
					),
				},
			},
		})
	})

	t.Run("error path - role collection import fails", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_assignment.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceRoleCollectionAssignmentBySubaccount("uut", "integration-test-acc-static", "Destination Administrator", "jenny.doe@test.com"),
				},
				{
					ResourceName:      "btp_subaccount_role_collection_assignment.uut",
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Import Not Supported`),
				},
			},
		})
	})

	t.Run("error path - subaccount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_subaccount_role_collection_assignment" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - role_collection_name mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_subaccount_role_collection_assignment" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "role_collection_name" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclResourceRoleCollectionAssignmentBySubaccount(resourceName string, subaccountName string, roleCollectionName string, userName string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_role_collection_assignment" "%s"{
    subaccount_id        = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	role_collection_name = "%s"
	user_name            = "%s"
}`, resourceName, subaccountName, roleCollectionName, userName)
}

func hclResourceRoleCollectionAssignmentWithOriginBySubaccount(resourceName string, subaccountName string, roleCollectionName string, userName string, origin string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_role_collection_assignment" "%s"{
    subaccount_id        = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	role_collection_name = "%s"
	user_name            = "%s"
	origin               = "%s"
}`, resourceName, subaccountName, roleCollectionName, userName, origin)
}

func hclResourceRoleCollectionAssignmentWithOriginAndGroupBySubaccount(resourceName string, subaccountName string, roleCollectionName string, groupName string, origin string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_role_collection_assignment" "%s"{
    subaccount_id        = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	role_collection_name = "%s"
	origin               = "%s"
	group_name           = "%s"
}`, resourceName, subaccountName, roleCollectionName, origin, groupName)
}

func hclResourceRoleCollectionAssignmentWithOriginAndAttributeBySubaccount(resourceName string, subaccountName string, roleCollectionName string, attributeName string, attributeValue string, origin string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}
resource "btp_subaccount_role_collection_assignment" "%s"{
    subaccount_id        = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	role_collection_name = "%s"
	origin               = "%s"
	attribute_name       = "%s"
	attribute_value      = "%s"
}`, resourceName, subaccountName, roleCollectionName, origin, attributeName, attributeValue)
}
