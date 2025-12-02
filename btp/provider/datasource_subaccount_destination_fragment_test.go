package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountDestinationFragment(t *testing.T) {

	t.Parallel()
	t.Run("happy path - destination fragment", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_fragment")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountDestinationFragment("test", "integration-test-destination", "test_destination_fragment"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_fragment.test", "fragment_content.FragmentName", "test_destination_fragment"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_fragment.test", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_fragment.test", "name", "test_destination_fragment"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_fragment.test", "id", regexpValidUUID),
						resource.TestCheckNoResourceAttr("data.btp_subaccount_destination_fragment.test", "service_instance_id"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination fragment with service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_fragment_with_service_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountDestinationFragmentWithServiceInstance("test", "integration-test-destination", "test_destination_fragment_service_instance", "servtest"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_fragment.test", "fragment_content.FragmentName", "test_destination_fragment_service_instance"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_fragment.test", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_fragment.test", "name", "test_destination_fragment_service_instance"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_fragment.test", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_fragment.test", "service_instance_id", regexpValidUUID),
					),
				},
			},
		})
	})

	t.Run("error path - subaccount is required", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_fragment_subaccount_required")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + `
						data "btp_subaccount_destination_fragment" "test" {
						  name = "test_destination_fragment"
						  }`,
					ExpectError: regexp.MustCompile(`The argument \"subaccount_id\" is required`),
				},
			},
		})
	})

	t.Run("error path - name is required", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_fragment_subaccount_required")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + `
						data "btp_subaccount_destination_fragment" "test" {
						  subaccount_id = "some-subaccount-id"
						  }`,
					ExpectError: regexp.MustCompile(`The argument \"name\" is required`),
				},
			},
		})
	})

	t.Run("error path - destination fragment doesn't exist on subaccount level", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_fragment_not_found_subaccount_level")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceSubaccountDestinationFragment("test", "integration-test-destination", "test_non_existent_fragment"),
					ExpectError: regexp.MustCompile(`test_non_existent_fragment`), // TODO improve error text
				},
			},
		})
	})

	t.Run("error path - destination fragment doesn't exist on service instance level", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_fragment_not_found_service_instance_level")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceSubaccountDestinationFragmentWithServiceInstance("test", "integration-test-destination", "test_non_existent_fragment", "servtest"),
					ExpectError: regexp.MustCompile(`test_non_existent_fragment`), // TODO improve error text
				},
			},
		})
	})

}

func hclDatasourceSubaccountDestinationFragment(resourceName string, subaccountName string, fragmentName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_destination_fragment" "%s" {	
subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
name = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, fragmentName)
}

func hclDatasourceSubaccountDestinationFragmentWithServiceInstance(resourceName string, subaccountName string, fragmentName string, serviceInstanceName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_instance" "fragment_instance" {
  subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  name          = "%s"
}
data "btp_subaccount_destination_fragment" "%s" {	
subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
name = "%s"
service_instance_id = data.btp_subaccount_service_instance.fragment_instance.id
}`
	return fmt.Sprintf(template, subaccountName, serviceInstanceName, resourceName, subaccountName, fragmentName)
}
