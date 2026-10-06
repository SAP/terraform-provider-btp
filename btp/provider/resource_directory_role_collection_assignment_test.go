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

func TestResourceDirectoryRoleCollectionAssignment(t *testing.T) {
	t.Parallel()
	t.Run("happy path - simple role collection assignment", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "jenny.doe@test.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection_assignment.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "role_collection_name", "Directory Viewer"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "user_name", "jenny.doe@test.com"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "origin", "ldap"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection assignment with origin", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.with_origin")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentWithOriginByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "john.doe@test.com", "iasprovidertestblr-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection_assignment.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "role_collection_name", "Directory Viewer"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "user_name", "john.doe@test.com"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "origin", "iasprovidertestblr-platform"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection assignment with origin and group", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.with_origin_and_group")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentWithOriginAndGroupByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "tf-test-group", "iasprovidertestblr-platform"),
					// We do not get back any information about the group, so if the call succeeds we assume that the asssignment/unassignment worked
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection_assignment.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "role_collection_name", "Directory Viewer"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "group_name", "tf-test-group"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "origin", "iasprovidertestblr-platform"),
					),
				},
			},
		})
	})

	t.Run("happy path - role collection assignment with origin and attribute", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.with_origin_and_attribute")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentWithOriginAndAttributeByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "tf_attr_name_test", "tf_attr_val_test", "iasprovidertestblr-platform"),
					// We do not get back any information about the group, so if the call succeeds we assume that the asssignment/unassignment worked
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory_role_collection_assignment.uut", "directory_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "role_collection_name", "Directory Viewer"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "attribute_name", "tf_attr_name_test"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "attribute_value", "tf_attr_val_test"),
						resource.TestCheckResourceAttr("btp_directory_role_collection_assignment.uut", "origin", "iasprovidertestblr-platform"),
					),
				},
			},
		})
	})

	t.Run("happy path - import role collection assignment by user", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.import")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "jenny.doe@test.com"),
				},
				{
					ResourceName:            "btp_directory_role_collection_assignment.uut",
					ImportState:             true,
					ImportStateIdFunc:       getImportStateIdForDirectoryRoleCollectionAssignmentUser("btp_directory_role_collection_assignment.uut", "Directory Viewer", "jenny.doe@test.com"),
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"id"},
				},
			},
		})
	})

	t.Run("happy path - import role collection assignment with resource identity", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.import_identity")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_12_0),
			},
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "jenny.doe@test.com"),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_directory_role_collection_assignment.uut", map[string]knownvalue.Check{
							"directory_id":         knownvalue.NotNull(),
							"role_collection_name": knownvalue.StringExact("Directory Viewer"),
							"user_name":            knownvalue.StringExact("jenny.doe@test.com"),
							"group_name":           knownvalue.Null(),
							"attribute_name":       knownvalue.Null(),
							"attribute_value":      knownvalue.Null(),
							"origin":               knownvalue.StringExact("ldap"),
						}),
					},
				},
				{
					ResourceName:    "btp_directory_role_collection_assignment.uut",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})

	t.Run("error path - import with invalid identifier", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/resource_directory_role_collection_assignment.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryRoleCollectionAssignmentByDirectory("uut", "integration-test-dir-se-static", "Directory Viewer", "jenny.doe@test.com"),
				},
				{
					ResourceName:  "btp_directory_role_collection_assignment.uut",
					ImportState:   true,
					ImportStateId: "invalid-id",
					ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
				},
			},
		})
	})

	t.Run("error path - directory_id mandatory", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_directory_role_collection_assignment" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "directory_id" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - role_collection_name mandatory", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_directory_role_collection_assignment" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "role_collection_name" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclResourceDirectoryRoleCollectionAssignmentByDirectory(resourceName string, directoryName string, roleCollectionName string, userName string) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_role_collection_assignment" "%s"{
    directory_id        = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
	role_collection_name = "%s"
	user_name            = "%s"
}`, resourceName, directoryName, roleCollectionName, userName)
}

func hclResourceDirectoryRoleCollectionAssignmentWithOriginByDirectory(resourceName string, directoryName string, roleCollectionName string, userName string, origin string) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_role_collection_assignment" "%s"{
    directory_id        = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
	role_collection_name = "%s"
	user_name            = "%s"
	origin               = "%s"
}`, resourceName, directoryName, roleCollectionName, userName, origin)
}

func hclResourceDirectoryRoleCollectionAssignmentWithOriginAndGroupByDirectory(resourceName string, directoryName string, roleCollectionName string, groupName string, origin string) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_role_collection_assignment" "%s"{
    directory_id        = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
	role_collection_name = "%s"
	origin               = "%s"
	group_name           = "%s"
}`, resourceName, directoryName, roleCollectionName, origin, groupName)
}

func hclResourceDirectoryRoleCollectionAssignmentWithOriginAndAttributeByDirectory(resourceName string, directoryName string, roleCollectionName string, attributeName string, attributeValue string, origin string) string {
	return fmt.Sprintf(`
data "btp_directories" "all" {}
resource "btp_directory_role_collection_assignment" "%s"{
    directory_id        = [for dir in data.btp_directories.all.values : dir.id if dir.name == "%s"][0]
	role_collection_name = "%s"
	origin               = "%s"
	attribute_name       = "%s"
	attribute_value      = "%s"
}`, resourceName, directoryName, roleCollectionName, origin, attributeName, attributeValue)
}

func getImportStateIdForDirectoryRoleCollectionAssignmentUser(resourceName, roleCollectionName, userName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s,%s,%s", rs.Primary.Attributes["directory_id"], roleCollectionName, userName, rs.Primary.Attributes["origin"]), nil
	}
}
