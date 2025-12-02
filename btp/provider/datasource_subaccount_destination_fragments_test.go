package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountDestinationFragments(t *testing.T) {

	t.Parallel()
	t.Run("happy path - destination fragment", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_fragments")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountDestinationFragments("test", "integration-test-destination"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_fragments.test", "values.#", "1"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_fragments.test", "values.0.fragment_content.FragmentName", "test_destination_fragment"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_fragments.test", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_fragments.test", "id", regexpValidUUID),
						resource.TestCheckNoResourceAttr("data.btp_subaccount_destination_fragments.test", "service_instance_id"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination fragment with service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_fragments_with_service_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountDestinationFragmentsWithServiceInstance("test", "integration-test-destination", "servtest"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_fragments.test", "values.#", "1"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_fragments.test", "values.0.fragment_content.FragmentName", "test_destination_fragment_service_instance"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_fragments.test", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_fragments.test", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_fragments.test", "service_instance_id", regexpValidUUID),
					),
				},
			},
		})
	})

	t.Run("error path - subaccount is required", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_fragments_subaccount_required")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + `
						data "btp_subaccount_destination_fragments" "test" {}`,
					ExpectError: regexp.MustCompile(`The argument \"subaccount_id\" is required`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountDestinationFragments(resourceName string, subaccountName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_destination_fragments" "%s" {	
subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName)
}

func hclDatasourceSubaccountDestinationFragmentsWithServiceInstance(resourceName string, subaccountName string, serviceInstanceName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_instance" "fragment_instance" {
  subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  name          = "%s"
}
data "btp_subaccount_destination_fragments" "%s" {	
subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
service_instance_id = data.btp_subaccount_service_instance.fragment_instance.id
}`
	return fmt.Sprintf(template, subaccountName, serviceInstanceName, resourceName, subaccountName)
}
