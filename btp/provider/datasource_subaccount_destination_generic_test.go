package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountDestinationGeneric(t *testing.T) {
	t.Parallel()
	t.Run("happy path without service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_generic_without_service_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDestinationGeneric("uut", "destination-resource", "integration-test-destination"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_generic.uut", "creation_time", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_generic.uut", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_generic.uut", "name", "destination-resource"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_generic.uut", "destination_configuration", "{\"Authentication\":\"NoAuthentication\",\"Description\":\"testing resource for destination update\",\"ProxyType\":\"Internet\",\"Type\":\"HTTP\",\"URL\":\"https://myservice.example.com\",\"abc\":\"good\"}"),
					),
				},
			},
		})
	})
	t.Run("happy path with service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_generic_with_service_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDestinationGenericSI("uut", "destination-resource-with-service-instance", "integration-test-destination", "servtest"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_generic.uut", "creation_time", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_generic.uut", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_generic.uut", "name", "destination-resource-with-service-instance"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_generic.uut", "service_instance_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_generic.uut", "destination_configuration", "{\"Authentication\":\"NoAuthentication\",\"ProxyType\":\"Internet\",\"Type\":\"HTTP\",\"URL\":\"https://myservice2.example.com\",\"description\":\"testing resource for destination update\"}"),
					),
				},
			},
		})
	})
	t.Run("error path - name not provided", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_subaccount_destination_generic" "test" {subaccount_id = "integration-test-destination"}`,
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found.`),
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
					Config:      `data "btp_subaccount_destination_generic" "test" {name = "res1"}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})
	t.Run("error path - destination not found", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_not_found_generic")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceDestinationGeneric("uut", "invalid_destination", "integration-test-destination"),
					ExpectError: regexp.MustCompile(`Configuration with the specified name was not found`),
				},
			},
		})
	})
}

func hclDatasourceDestinationGeneric(datasourceName string, name string, subaccountName string) string {
	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_destination_generic" "%s" {
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	}`
	return fmt.Sprintf(template, datasourceName, name, subaccountName)
}

func hclDatasourceDestinationGenericSI(datasourceName string, name string, subaccountName string, serviceInstanceName string) string {
	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_service_instance" "dest" {
  		subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  		name          = "%s"
	}
	data "btp_subaccount_destination_generic" "%s" {
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	service_instance_id = data.btp_subaccount_service_instance.dest.id
	}`
	return fmt.Sprintf(template, subaccountName, serviceInstanceName, datasourceName, name, subaccountName)
}
