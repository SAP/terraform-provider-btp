package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceDestination(t *testing.T) {
	t.Parallel()
	t.Run("happy path with and without service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDestination("uut", "destination-resource", "integration-test-destination"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "authentication", "NoAuthentication"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination.uut", "creation_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "description", "testing resource for destination update"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination.uut", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "name", "destination-resource"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "proxy_type", "Internet"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "type", "HTTP"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "additional_configuration", "{\"abc\":\"good\"}"),
					),
				},
			},
		})
	})
	t.Run("happy path with and without service instance", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceDestinationSI("uut", "destination-resource-with-service-instance", "integration-test-destination", "servtest"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "authentication", "OAuth2ClientCredentials"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination.uut", "creation_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "description", "testing resource for destination update"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination.uut", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "name", "destination-resource-with-service-instance"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "proxy_type", "Internet"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination.uut", "service_instance_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "type", "HTTP"),
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
					Config:      `data "btp_subaccount_destination" "test" {subaccount_id = "integration-test-destination"}`,
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
					Config:      `data "btp_subaccount_destination" "test" {name = "res1"}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})
	t.Run("error path - destination not found", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceDestination("uut", "invalid_destination", "integration-test-destination"),
					ExpectError: regexp.MustCompile(`Configuration with the specified name was not found`),
				},
			},
		})
	})
}

func hclDatasourceDestination(datasourceName string, name string, subaccount string) string {
	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_destination" "%s" {
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	}`
	return fmt.Sprintf(template, datasourceName, name, subaccount)
}

func hclDatasourceDestinationSI(datasourceName string, name string, subaccountName string, serviceInstanceName string) string {
	template := `
	data "btp_subaccounts" "all" {}
	data "btp_subaccount_service_instance" "dest" {
  		subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  		name          = "%s"
	}
	data "btp_subaccount_destination" "%s" {
	name = "%s"
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	service_instance_id = data.btp_subaccount_service_instance.dest.id
	}`
	return fmt.Sprintf(template, subaccountName, serviceInstanceName, datasourceName, name, subaccountName)
}
