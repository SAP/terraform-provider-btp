package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
	res "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestSubaccountRoleCollectionAssignmentListResource(t *testing.T) {
	t.Parallel()

	subaccountID := "fc26cc61-ac5e-4c7d-9747-725f32a8994e"

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_role_collection_assignment")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			Steps: []resource.TestStep{
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountRoleCollectionAssignmentQueryConfig(
						"subaccount_role_collection_assignment_list",
						"btp",
						subaccountID,
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role_collection_assignment.subaccount_role_collection_assignment_list", 26),
					},
				},
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountRoleCollectionAssignmentQueryConfigWithRoleCollectionName(
						"subaccount_role_collection_assignment_list2",
						"btp",
						subaccountID,
						"Subaccount Administrator",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role_collection_assignment.subaccount_role_collection_assignment_list2", 8),
					},
				},
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountRoleCollectionAssignmentQueryConfigWithOrigin(
						"subaccount_role_collection_assignment_list3",
						"btp",
						subaccountID,
						"iasprovidertestblr-platform",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role_collection_assignment.subaccount_role_collection_assignment_list3", 14),
					},
				},
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountRoleCollectionAssignmentQueryConfigWithRoleCollectionNameAndOrigin(
						"subaccount_role_collection_assignment_list4",
						"btp",
						subaccountID,
						"test-role-collection-attribute",
						"iasprovidertestblr-platform",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role_collection_assignment.subaccount_role_collection_assignment_list4", 1),

						querycheck.ExpectIdentity("btp_subaccount_role_collection_assignment.subaccount_role_collection_assignment_list4",
							map[string]knownvalue.Check{
								"role_collection_name": knownvalue.StringExact("test-role-collection-attribute"),
								"subaccount_id":        knownvalue.StringExact("fc26cc61-ac5e-4c7d-9747-725f32a8994e"),
								"user_name":            knownvalue.Null(),
								"group_name":           knownvalue.Null(),
								"attribute_name":       knownvalue.StringExact("test1"),
								"attribute_value":      knownvalue.StringExact("value1"),
								"origin":               knownvalue.StringExact("iasprovidertestblr-platform"),
							},
						),
					},
				},
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountRoleCollectionAssignmentQueryConfigWithRoleCollectionNameAndOrigin(
						"subaccount_role_collection_assignment_list5",
						"btp",
						subaccountID,
						"test-role-collection-group",
						"iasproviderdevblr",
					),

					QueryResultChecks: []querycheck.QueryResultCheck{
						querycheck.ExpectLength("btp_subaccount_role_collection_assignment.subaccount_role_collection_assignment_list5", 1),

						querycheck.ExpectIdentity("btp_subaccount_role_collection_assignment.subaccount_role_collection_assignment_list5",
							map[string]knownvalue.Check{
								"role_collection_name": knownvalue.StringExact("test-role-collection-group"),
								"subaccount_id":        knownvalue.StringExact("fc26cc61-ac5e-4c7d-9747-725f32a8994e"),
								"user_name":            knownvalue.Null(),
								"group_name":           knownvalue.StringExact("test-group"),
								"attribute_name":       knownvalue.Null(),
								"attribute_value":      knownvalue.Null(),
								"origin":               knownvalue.StringExact("iasproviderdevblr"),
							},
						),
					},
				},
			},
		})
	})

	t.Run("error path - Access forbidden", func(t *testing.T) {
		t.Parallel()
		notFoundSubaccountID := "00000000-0000-0000-0000-000000000000"

		rec, user := setupVCR(t, "fixtures/list_resource_subaccount_role_collection_assignment_access_forbidden")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			Steps: []resource.TestStep{
				{
					Query: true,
					Config: hclProviderFor(user) + listSubaccountRoleCollectionAssignmentQueryConfig(
						"not_found_test",
						"btp",
						notFoundSubaccountID,
					),

					ExpectError: regexp.MustCompile(`Access forbidden due to insufficient authorization|403`),
				},
			},
		})
	})

	t.Run("error path - configure", func(t *testing.T) {
		t.Parallel()
		r := NewSubaccountRoleCollectionAssignmentListResource().(list.ListResourceWithConfigure)
		resp := &res.ConfigureResponse{}
		req := res.ConfigureRequest{
			ProviderData: struct{}{}, // Wrong type
		}

		r.Configure(context.Background(), req, resp)

		if !resp.Diagnostics.HasError() {
			t.Error("Expected error for invalid provider data type")
		}
	})
}

func listSubaccountRoleCollectionAssignmentQueryConfig(label, providerName, subaccountID string) string {
	return fmt.Sprintf(`
list "btp_subaccount_role_collection_assignment" "%s" {
  provider = "%s"
  config {
    subaccount_id = "%s"
  }
}`, label, providerName, subaccountID)
}

func listSubaccountRoleCollectionAssignmentQueryConfigWithRoleCollectionName(label, providerName, subaccountID, roleCollectionName string) string {
	return fmt.Sprintf(`
list "btp_subaccount_role_collection_assignment" "%s" {
  provider = "%s"
  config {
    subaccount_id = "%s"
    role_collection_name = "%s"
  }
}`, label, providerName, subaccountID, roleCollectionName)
}

func listSubaccountRoleCollectionAssignmentQueryConfigWithOrigin(label, providerName, subaccountID, origin string) string {
	return fmt.Sprintf(`
list "btp_subaccount_role_collection_assignment" "%s" {
  provider = "%s"
  config {
    subaccount_id = "%s"
    origin = "%s"
  }
}`, label, providerName, subaccountID, origin)
}

func listSubaccountRoleCollectionAssignmentQueryConfigWithRoleCollectionNameAndOrigin(label, providerName, subaccountID, roleCollectionName, origin string) string {
	return fmt.Sprintf(`
list "btp_subaccount_role_collection_assignment" "%s" {
  provider = "%s"
  config {
    subaccount_id = "%s"
    role_collection_name = "%s"
    origin = "%s"
  }
}`, label, providerName, subaccountID, roleCollectionName, origin)
}

func subaccountRoleCollectionAssignmentQueryImportConfig(subaccountID string) string {
	return fmt.Sprintf(`
resource "btp_subaccount_role_collection_assignment" "attribute_assignment" {
  subaccount_id       = "%s"
  role_collection_name = "test-role-collection-attribute"
  attribute_name      = "test1"
  attribute_value     = "value1"
  origin              = "iasprovidertestblr-platform"
}

import {
  to = btp_subaccount_role_collection_assignment.attribute_assignment
  identity = {
    subaccount_id        = "%s"
    role_collection_name = "test-role-collection-attribute"
    attribute_name       = "test1"
    attribute_value      = "value1"
    origin               = "iasprovidertestblr-platform"
  }
}

resource "btp_subaccount_role_collection_assignment" "group_assignment" {
  subaccount_id       = "%s"
  role_collection_name = "test-role-collection-group"
  group_name          = "test-group"
  origin              = "iasproviderdevblr"
}

import {
  to = btp_subaccount_role_collection_assignment.group_assignment
  identity = {
    subaccount_id        = "%s"
    role_collection_name = "test-role-collection-group"
    group_name           = "test-group"
    origin               = "iasproviderdevblr"
  }
}`, subaccountID, subaccountID, subaccountID, subaccountID)
}
