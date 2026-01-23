package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountDestinationsNames(t *testing.T) {
	t.Parallel()

	t.Run("happy path without service instance only names", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destinations_names_without_service_instance_v2")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDestinationsNamesV2("uut", "integration-test-destination"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destinations_names.uut", "destination_names.#", "4"),
					),
				},
			},
		})
	})
	t.Run("happy path with service instance only names", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destinations_names_with_service_instance_v2")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDestinationsSINamesV2("uut", "integration-test-destination", "servtest"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destinations_names.uut", "destination_names.#", "2"),
					),
				},
			},
		})
	})
	t.Run("error path - subaccount not provided", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_subaccount_destinations_names" "test" {}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclDatasourceDestinationsNamesV2(datasourceName string, subaccount string) string {
	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_destinations_names" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	}`
	return fmt.Sprintf(template, datasourceName, subaccount)
}

func hclDatasourceDestinationsSINamesV2(datasourceName string, subaccountName string, serviceInstanceName string) string {
	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_service_instance" "dest" {
  		subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  		name          = "%s"
	}
	data "btp_subaccount_destinations_names" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	service_instance_id = data.btp_subaccount_service_instance.dest.id
	}`
	return fmt.Sprintf(template, subaccountName, serviceInstanceName, datasourceName, subaccountName)
}
