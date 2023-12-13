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
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceInstanceBySubaccountIdByIdFilteredFromList("uut", "integration-test-services-static", "tf-testacc-alertnotification-instance"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instance.uut", "id", regexpValidUUID),
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
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceInstanceBySubaccountIdByNameFilteredFromList("uut", "integration-test-services-static", "tf-testacc-alertnotification-instance"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_instance.uut", "id", regexpValidUUID),
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
					Config:      hclDatasourceSubaccountServiceInstanceIdName("uut", "00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", "any-service-instance-name"),
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

func hclDatasourceSubaccountServiceInstanceBySubaccountIdByIdFilteredFromList(resourceName string, subaccountName string, serviceInstanceName string) string {
	template := `
data "btp_subaccounts" "allsas" {}
data "btp_subaccount_service_instances" "allssis" {
  subaccount_id = [for sa in data.btp_subaccounts.allsas.values : sa.id if sa.name == "%[2]s"][0]
}
data "btp_subaccount_service_instance" "%[1]s" {
     subaccount_id = [for sa in data.btp_subaccounts.allsas.values : sa.id if sa.name == "%[2]s"][0]
	 id            = [for ssi in data.btp_subaccount_service_instances.allssis.values : ssi.id if ssi.name == "%[3]s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName, serviceInstanceName)
}

func hclDatasourceSubaccountServiceInstanceBySubaccountIdByNameFilteredFromList(resourceName string, subaccountName string, serviceName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_service_instance" "%s" {
     subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	 name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, serviceName)
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
