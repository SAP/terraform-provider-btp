package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceInstance(t *testing.T) {

	t.Parallel()
	t.Run("happy path - service instance by id", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_instance.by_id")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceInstanceById("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "df532d07-57a7-415e-a261-23a398ef068a"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "id", "df532d07-57a7-415e-a261-23a398ef068a"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "name", "tf-testacc-alertnotification-instance"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "ready", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "serviceplan_id", "f0aac855-474d-4016-9529-61c062efbc7c"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "platform_id", "service-manager"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
			},
		})
	})
	t.Run("happy path - service instance by name", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_instance.by_name")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceInstanceByName("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "tf-testacc-alertnotification-instance"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "id", "df532d07-57a7-415e-a261-23a398ef068a"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "name", "tf-testacc-alertnotification-instance"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "ready", "true"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "serviceplan_id", "f0aac855-474d-4016-9529-61c062efbc7c"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "platform_id", "service-manager"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
			},
		})
	})

	t.Run("error path - specify ID and name", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountServiceInstanceIdName("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "tf-testacc-alertnotification-instance"),
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Combination`),
				},
			},
		})
	})

	t.Run("error path - subaccount id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountServiceInstanceWoSubaccount("uut", "lite"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountServiceInstanceById(resourceName string, subaccountId string, serviceId string) string {
	template := `
data "btp_subaccount_service_instance" "%s" {
     subaccount_id = "%s"
	 id            = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, serviceId)
}

func hclDatasourceSubaccountServiceInstanceByName(resourceName string, subaccountId string, serviceName string) string {
	template := `
data "btp_subaccount_service_instance" "%s" {
     subaccount_id = "%s"
	 name            = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, serviceName)
}

func hclDatasourceSubaccountServiceInstanceIdName(resourceName string, subaccountId string, serviceId string, serviceName string) string {
	template := `
data "btp_subaccount_service_instance" "%s" {
    subaccount_id = "%s"
	id            = "%s"
    name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, serviceId, serviceName)
}

func hclDatasourceSubaccountServiceInstanceWoSubaccount(resourceName string, serviceName string) string {
	template := `
data "btp_subaccount_service_instance" "%s" {
    name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, serviceName)
}
