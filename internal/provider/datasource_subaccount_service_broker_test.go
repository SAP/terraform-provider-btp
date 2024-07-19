package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountServiceBroker(t *testing.T) {

	t.Parallel()
	t.Run("happy path - service brokers by id", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_service_broker.by_id")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountServiceBrokerBySubaccountNameByBrokerName("uut", "integration-test-services-static", "my-broker"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_broker.uut", "name", regexp.MustCompile("^my-broker-.+")),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_broker.uut", "url", "https://my.broker.test"),
						resource.TestCheckResourceAttr("data.btp_subaccount_service_broker.uut", "ready", "true"),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_broker.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_subaccount_service_broker.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
			},
		})

	})

	t.Run("error path - subaccount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountServiceBrokerNoSubaccount("uut", "any-sb-name"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - no ID or name", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountServiceBrokerNoIdOrName("uut", "00000000-0000-0000-0000-000000000000"),
					ExpectError: regexp.MustCompile(`Error: Invalid Attribute Combination`),
				},
			},
		})
	})

	t.Run("error path - subaccount_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountServiceBrokerByNameBySubaccountIdByBrokerName("uut", "this-is-not-a-uuid", "any-sb-name"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountServiceBrokerBySubaccountNameByBrokerName(resourceName string, subaccountName string, brokerName string) string {
	template := `
data "btp_subaccounts" "allsas" {}
data "btp_subaccount_service_brokers" "allsbs" {
  subaccount_id = [for sa in data.btp_subaccounts.allsas.values : sa.id if sa.name == "%[2]s"][0]
}
data "btp_subaccount_service_broker" "%[1]s" {
	subaccount_id = [for sa in data.btp_subaccounts.allsas.values : sa.id if sa.name == "%[2]s"][0]
	id            = [for sb in data.btp_subaccount_service_brokers.allsbs.values : sb.id if startswith(sb.name, "%[3]s")][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName, brokerName)
}

func hclDatasourceSubaccountServiceBrokerByNameBySubaccountIdByBrokerName(resourceName string, subaccountId string, brokerName string) string {
	template := `data "btp_subaccount_service_broker" "%s" {
	subaccount_id = "%s"
	name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, brokerName)
}

func hclDatasourceSubaccountServiceBrokerNoSubaccount(resourceName string, brokerName string) string {
	template := `data "btp_subaccount_service_broker" "%s" {
	name          = "%s"
}`
	return fmt.Sprintf(template, resourceName, brokerName)
}

func hclDatasourceSubaccountServiceBrokerNoIdOrName(resourceName string, subaccountId string) string {
	template := `data "btp_subaccount_service_broker" "%s" {
	subaccount_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId)
}
