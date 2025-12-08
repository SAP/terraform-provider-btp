package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountDestinations(t *testing.T) {
	t.Parallel()
	t.Run("happy path without service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destinations_without_service_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDestinations("uut", "integration-test-destination"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destinations.uut", "values.#", "3"),
					),
				},
			},
		})
	})
	t.Run("happy path with service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destinations_with_service_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDestinationsSI("uut", "integration-test-destination", "servtest"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destinations.uut", "values.#", "2"),
					),
				},
			},
		})
	})
	t.Run("happy path without service instance only names", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destinations_names_without_service_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDestinationsNames("uut", "integration-test-destination"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destinations.uut", "destination_names.#", "3"),
					),
				},
			},
		})
	})
	t.Run("happy path with service instance only names", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destinations_names_with_service_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDestinationsSINames("uut", "integration-test-destination", "servtest"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destinations.uut", "destination_names.#", "2"),
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
					Config:      `data "btp_subaccount_destinations" "test" {}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclDatasourceDestinations(datasourceName string, subaccount string) string {
	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_destinations" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	}`
	return fmt.Sprintf(template, datasourceName, subaccount)
}
func hclDatasourceDestinationsNames(datasourceName string, subaccount string) string {
	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_destinations" "%s" {
	names_only = true
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	}`
	return fmt.Sprintf(template, datasourceName, subaccount)
}

func hclDatasourceDestinationsSI(datasourceName string, subaccountName string, serviceInstanceName string) string {
	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_service_instance" "dest" {
  		subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  		name          = "%s"
	}
	data "btp_subaccount_destinations" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	service_instance_id = data.btp_subaccount_service_instance.dest.id
	}`
	return fmt.Sprintf(template, subaccountName, serviceInstanceName, datasourceName, subaccountName)
}

func hclDatasourceDestinationsSINames(datasourceName string, subaccountName string, serviceInstanceName string) string {
	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_service_instance" "dest" {
  		subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  		name          = "%s"
	}
	data "btp_subaccount_destinations" "%s" {
	names_only = true
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	service_instance_id = data.btp_subaccount_service_instance.dest.id
	}`
	return fmt.Sprintf(template, subaccountName, serviceInstanceName, datasourceName, subaccountName)
}
