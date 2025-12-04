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

func TestResourceSubaccountDestinationFragment(t *testing.T) {
	t.Parallel()
	t.Run("happy path: destination fragment without service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_fragment_without_service_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationFragment(
						"test",
						"integration-test-destination",
						"integration-test-fragment"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination_fragment.test", "name", "integration-test-fragment"),
						resource.TestMatchResourceAttr("btp_subaccount_destination_fragment.test", "subaccount_id", regexpValidUUID),
						resource.TestCheckNoResourceAttr("btp_subaccount_destination_fragment.test", "service_instance_id"),
						resource.TestCheckNoResourceAttr("btp_subaccount_destination_fragment.test", "fragment_content"),
					),
				},
				{
					ResourceName:      "btp_subaccount_destination_fragment.test",
					ImportStateIdFunc: getIdForSubaccounDestinationFragmentImportId("btp_subaccount_destination_fragment.test", "integration-test-fragment"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - import with resource identity", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_fragment.import_by_resource_identity")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_12_0),
			},
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationFragment(
						"test",
						"integration-test-destination",
						"integration-test-fragment"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination_fragment.test", "name", "integration-test-fragment"),
						resource.TestMatchResourceAttr("btp_subaccount_destination_fragment.test", "subaccount_id", regexpValidUUID),
						resource.TestCheckNoResourceAttr("btp_subaccount_destination_fragment.test", "service_instance_id"),
						resource.TestCheckNoResourceAttr("btp_subaccount_destination_fragment.test", "fragment_content"),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_subaccount_destination_fragment.test", map[string]knownvalue.Check{
							"subaccount_id":       knownvalue.StringRegexp(regexpValidUUID),
							"name":                knownvalue.StringExact("integration-test-fragment"),
							"service_instance_id": knownvalue.Null(),
						}),
					},
				},
				{
					ResourceName:    "btp_subaccount_destination_fragment.test",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})

	t.Run("happy path: destination fragment with service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_fragment_with_service_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationFragmentWithServiceInstance(
						"test_si",
						"integration-test-destination",
						"integration-test-fragment-with-service-instance",
						"servtest"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination_fragment.test_si", "name", "integration-test-fragment-with-service-instance"),
						resource.TestMatchResourceAttr("btp_subaccount_destination_fragment.test_si", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_destination_fragment.test_si", "service_instance_id", regexpValidUUID),
						resource.TestCheckNoResourceAttr("btp_subaccount_destination_fragment.test_si", "fragment_content"),
					),
				},
			},
		})
	})

	t.Run("happy path - import destination fragment service instance using resource identity", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_fragment_with_service_instance.import_by_resource_identity")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_12_0),
			},
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationFragmentWithServiceInstance(
						"test_si",
						"integration-test-destination",
						"integration-test-fragment-with-service-instance",
						"servtest"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination_fragment.test_si", "name", "integration-test-fragment-with-service-instance"),
						resource.TestMatchResourceAttr("btp_subaccount_destination_fragment.test_si", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_destination_fragment.test_si", "service_instance_id", regexpValidUUID),
						resource.TestCheckNoResourceAttr("btp_subaccount_destination_fragment.test_si", "fragment_content"),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity("btp_subaccount_destination_fragment.test_si", map[string]knownvalue.Check{
							"subaccount_id":       knownvalue.StringRegexp(regexpValidUUID),
							"name":                knownvalue.StringExact("integration-test-fragment-with-service-instance"),
							"service_instance_id": knownvalue.StringRegexp(regexpValidUUID),
						}),
					},
				},
				{
					ResourceName:    "btp_subaccount_destination_fragment.test_si",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})

	t.Run("error path - name and subaccount_id are mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `resource "btp_subaccount_destination_fragment" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "(name|subaccount_id)" is required, but no definition was found.`),
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
					Config:      hclResourceSubaccountDestinationFragment("uut", "integration-test-destination", ""),
					ExpectError: regexp.MustCompile(`Attribute name string length must be at least 1, got: 0`),
				},
			},
		})
	})

}

func hclResourceSubaccountDestinationFragment(resourceName string, subaccountName, fragmentName string) string {
	return fmt.Sprintf(`
	data "btp_subaccounts" "all" {}
	resource "btp_subaccount_destination_fragment" "%s" {
		subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
		name = "%s"	
    }`, resourceName, subaccountName, fragmentName)
}

func hclResourceSubaccountDestinationFragmentWithServiceInstance(resourceName string, subaccountName string, fragmentName string, serviceInstanceName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_instance" "fragment_instance" {
  subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  name          = "%s"
}
resource "btp_subaccount_destination_fragment" "%s" {	
subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
name = "%s"
service_instance_id = data.btp_subaccount_service_instance.fragment_instance.id
}`
	return fmt.Sprintf(template, subaccountName, serviceInstanceName, resourceName, subaccountName, fragmentName)
}

func getIdForSubaccounDestinationFragmentImportId(resourceName, fragmentName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["subaccount_id"], fragmentName), nil
	}
}
