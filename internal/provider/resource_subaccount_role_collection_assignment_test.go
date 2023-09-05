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
					Config: hclProviderFor(user) + hclResourceRoleCollectionAssignment("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "Destination Administrator", "jenny.doe@test.com"),
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
					Config: hclProviderFor(user) + hclResourceRoleCollectionAssignmentWithOrigin("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "Destination Administrator", "john.doe@test.com", "terraformint-platform"),
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

	t.Run("error path - role collection import fails", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_role_collection_assignment.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceRoleCollectionAssignment("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "Destination Administrator", "jenny.doe@test.com"),
				},
				{
					ResourceName:      "btp_subaccount_role_collection_assignment.uut",
					ImportStateId:     "ef23ace8-6ade-4d78-9c1f-8df729548bbf",
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

func hclResourceRoleCollectionAssignment(resourceName string, subaccountId string, roleCollectionName string, userName string) string {

	return fmt.Sprintf(`
resource "btp_subaccount_role_collection_assignment" "%s"{
    subaccount_id        = "%s"
	role_collection_name = "%s"
	user_name            = "%s"
}`, resourceName, subaccountId, roleCollectionName, userName)
}

func hclResourceRoleCollectionAssignmentWithOrigin(resourceName string, subaccountId string, roleCollectionName string, userName string, origin string) string {

	return fmt.Sprintf(`
resource "btp_subaccount_role_collection_assignment" "%s"{
    subaccount_id        = "%s"
	role_collection_name = "%s"
	user_name            = "%s"
	origin               = "%s"
}`, resourceName, subaccountId, roleCollectionName, userName, origin)
}
