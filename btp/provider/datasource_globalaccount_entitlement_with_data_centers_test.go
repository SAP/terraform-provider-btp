package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountEntitlementWithDataCenters(t *testing.T) {
	t.Parallel()

	t.Run("happy path - retrieve data centers for entitlement", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_entitlement_with_data_centers")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountEntitlementWithDataCenters("uut", "hana-cloud", "hana"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlement_with_data_centers.uut", "service_name", "hana-cloud"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlement_with_data_centers.uut", "plan_name", "hana"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlement_with_data_centers.uut", "datacenter_information.%", "6"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlement_with_data_centers.uut", "datacenter_information.eu12.dc_region", "eu12"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlement_with_data_centers.uut", "datacenter_information.eu12.dc_name", "cf-eu12"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlement_with_data_centers.uut", "datacenter_information.eu12.dc_display_name", "cf-eu12"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlement_with_data_centers.uut", "datacenter_information.eu12.dc_iaas_provider", "AWS"),
					),
				},
			},
		})
	})

	t.Run("happy path - no entry found", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_entitlement_with_data_centers_empty")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountEntitlementWithDataCenters("uut", "some-service", "some-service-plan"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlement_with_data_centers.uut", "service_name", "some-service"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlement_with_data_centers.uut", "plan_name", "some-service-plan"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_entitlement_with_data_centers.uut", "datacenter_information.%", "0"),
					),
				},
			},
		})
	})

	t.Run("error path - service_name is required", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_entitlement_with_data_centers_required_service_name")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + `
						data "btp_globalaccount_entitlement_with_data_centers" "uut" {
						  plan_name    = "hana"
						}
					`,
					ExpectError: regexp.MustCompile(`The argument \"service_name\" is required`),
				},
			},
		})
	})

	t.Run("error path - plan_name is required", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_entitlement_with_data_centers_required_plan_name")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + `
						data "btp_globalaccount_entitlement_with_data_centers" "uut" {
						  service_name    = "hana"
						}
					`,
					ExpectError: regexp.MustCompile(`The argument \"plan_name\" is required`),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountEntitlementWithDataCenters(resourceName, serviceName, planName string) string {
	return fmt.Sprintf(`

	data "btp_globalaccount_entitlement_with_data_centers" "%s" {
	  service_name  = "%s"
	  plan_name     = "%s"
	}
	`, resourceName, serviceName, planName)
}
