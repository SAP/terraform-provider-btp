package provider

import (
	"fmt"
	"regexp"
	"strings"
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
					Config: hclProviderFor(user) + hclDatasourceDestination("uut", "destination-resource", "ba268910-81e6-4ac1-9016-cae7ed196889", ""),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "subaccount_id", "ba268910-81e6-4ac1-9016-cae7ed196889"),
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
					Config: hclProviderFor(user) + hclDatasourceDestination("uut", "destination-resource-with-service-instance", "ba268910-81e6-4ac1-9016-cae7ed196889", "499614ab-4e09-4836-9cf8-889f0a39b4b3"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "subaccount_id", "ba268910-81e6-4ac1-9016-cae7ed196889"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "authentication", "OAuth2ClientCredentials"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination.uut", "creation_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "description", "testing resource for destination update"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination.uut", "modification_time", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "name", "destination-resource-with-service-instance"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "proxy_type", "Internet"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination.uut", "service_instance_id", "499614ab-4e09-4836-9cf8-889f0a39b4b3"),
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
					Config:      `data "btp_subaccount_destination" "test" {subaccount_id = "ba268910-81e6-4ac1-9016-cae7ed196889"}`,
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
					Config:      hclProviderFor(user) + hclDatasourceDestination("uut", "invalid_destination", "ba268910-81e6-4ac1-9016-cae7ed196889", "499614ab-4e09-4836-9cf8-889f0a39b4b3"),
					ExpectError: regexp.MustCompile(`Configuration with the specified name was not found`),
				},
			},
		})
	})
}

func hclDatasourceDestination(datasourceName string, name string, subaccount string, serviceInstance string) string {
	var serviceInstanceLine string
	if strings.TrimSpace(serviceInstance) != "" {
		serviceInstanceLine = fmt.Sprintf(`service_instance_id = "%s"`, serviceInstance)
	}
	template := `data "btp_subaccount_destination" "%s" {
	name = "%s"
	subaccount_id = "%s"
	%s
	}`
	return fmt.Sprintf(template, datasourceName, name, subaccount, serviceInstanceLine)
}
