package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGlobalaccountRoleCollectionAssignment(t *testing.T) {
	t.Parallel()
	t.Run("happy path - simple role collection assignment", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role_collection_assignment")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollectionAssignment("uut", "Global Account Viewer", "jenny.doe@test.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection_assignment.uut", "role_collection_name", "Global Account Viewer"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection_assignment.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection_assignment.uut", "origin", "ldap"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection assignment with origin", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role_collection_assignment.with_origin")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollectionAssignmentWithOrigin("uut", "Global Account Viewer", "john.doe@test.com", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection_assignment.uut", "role_collection_name", "Global Account Viewer"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection_assignment.uut", "user_name", "john.doe@test.com"),
						resource.TestCheckResourceAttr("btp_globalaccount_role_collection_assignment.uut", "origin", "terraformint-platform"),
					),
				},
			},
		})
	})

	t.Run("error path - role collection import fails", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_role_collection_assignment.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountRoleCollectionAssignment("uut", "Global Account Viewer", "jenny.doe@test.com"),
				},
				{
					ResourceName:      "btp_globalaccount_role_collection_assignment.uut",
					ImportStateId:     "anyID",
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Import Not Supported`),
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
					Config:      `resource "btp_globalaccount_role_collection_assignment" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "role_collection_name" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclResourceGlobalaccountRoleCollectionAssignment(resourceName string, roleCollectionName string, userName string) string {

	return fmt.Sprintf(`
resource "btp_globalaccount_role_collection_assignment" "%s"{
	role_collection_name = "%s"
	user_name            = "%s"
}`, resourceName, roleCollectionName, userName)
}

func hclResourceGlobalaccountRoleCollectionAssignmentWithOrigin(resourceName string, roleCollectionName string, userName string, origin string) string {

	return fmt.Sprintf(`
resource "btp_globalaccount_role_collection_assignment" "%s"{
	role_collection_name = "%s"
	user_name            = "%s"
	origin               = "%s"
}`, resourceName, roleCollectionName, userName, origin)
}
