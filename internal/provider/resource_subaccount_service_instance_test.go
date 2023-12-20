package provider

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

type XsuaaParameters struct {
	Xsappname  string `json:"xsappname"`
	TenantMode string `json:"tenant-mode"`
}

func TestResourceSubaccountServiceInstance(t *testing.T) {
	t.Run("happy path - simple service creation wo parameters", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_instance.wo_parameters")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceInstanceWoParametersBySubaccountByServicePlan("uut", "integration-test-services-static", "tf-test-audit-log", "default", "auditlog-management"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "tf-test-audit-log"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "platform_id", "service-manager"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceInstanceWoParametersBySubaccountByServicePlan("uut", "integration-test-services-static", "TF-TEST-AUDIT-LOG", "default", "auditlog-management"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "TF-TEST-AUDIT-LOG"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "platform_id", "service-manager"),
					),
				},
				{
					ResourceName:      "btp_subaccount_service_instance.uut",
					ImportStateIdFunc: getServiceInstanceIdForImport("btp_subaccount_service_instance.uut"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - simple service creation with timeout", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_instance_with_timeouts")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceInstanceWoParametersWithTimeoutsBySubaccountByServicePlan("uut",
						"integration-test-services-static",
						"tf-test-audit-log",
						"default",
						"auditlog-management",
						"15m",
						"15m",
						"20m"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "tf-test-audit-log"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "platform_id", "service-manager"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "timeouts.create", "15m"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "timeouts.update", "15m"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "timeouts.delete", "20m"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceInstanceWoParametersWithTimeoutsBySubaccountByServicePlan("uut",
						"integration-test-services-static",
						"TF-TEST-AUDIT-LOG",
						"default",
						"auditlog-management",
						"15m",
						"15m",
						"20m"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "TF-TEST-AUDIT-LOG"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "platform_id", "service-manager"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "timeouts.create", "15m"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "timeouts.update", "15m"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "timeouts.delete", "20m"),
					),
				},
			},
		})
	})

	t.Run("happy path - simple service creation with parameters", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_instance.with_parameters")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceInstanceWithParametersBySubaccountByServicePlan("uut", "integration-test-services-static", "tf-test-xsuaa", "application", "xsuaa"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "tf-test-xsuaa"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "platform_id", "service-manager"),
					),
				},
			},
		})
	})

	t.Run("happy path - simple service creation with labels", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_instance.with_labels")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceInstanceWithLabelsBySubaccountByPlan("uut", "integration-test-services-static", "tf-test-malware-scanner", "clamav", "malware-scanner"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "tf-test-malware-scanner"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "labels.foo.0", "bar"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceInstanceWithLabelsBySubaccountByPlan("uut", "integration-test-services-static", "TF-TEST-MALWARE-SCANNER", "clamav", "malware-scanner"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "TF-TEST-MALWARE-SCANNER"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "labels.foo.0", "bar"),
					),
				},
				{
					ResourceName:      "btp_subaccount_service_instance.uut",
					ImportStateIdFunc: getServiceInstanceIdForImport("btp_subaccount_service_instance.uut"),
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})

	t.Run("happy path - simple service creation with labels change", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_instance.with_labels_change")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceInstanceWithLabelsBySubaccountByPlan("uut", "integration-test-services-static", "tf-test-malware-scanner", "clamav", "malware-scanner"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "tf-test-malware-scanner"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "labels.foo.0", "bar"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceInstanceWithLabelsChangedBySubaccountByPlan("uut", "integration-test-services-static", "TF-TEST-MALWARE-SCANNER", "clamav", "malware-scanner"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "serviceplan_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_instance.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "usable", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "name", "TF-TEST-MALWARE-SCANNER"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "labels.foo.0", "BAR"),
						resource.TestCheckResourceAttr("btp_subaccount_service_instance.uut", "labels.bar.0", "foo"),
					),
				},
			},
		})
	})

	t.Run("error path - subacount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountServiceInstanceNoSubaccountId("uut", "tf-test-audit-log", "00000000-0000-0000-0000-000000000000"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - service name mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountServiceInstanceNoServicName("uut", "00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000"),
					ExpectError: regexp.MustCompile(`The argument "name" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - service plan ID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountServiceInstanceNoPlan("uut", "this-is-not-a-uuid", "tf-test-audit-log"),
					ExpectError: regexp.MustCompile(`The argument "serviceplan_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - import failure", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_instance.import_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceInstanceWoParametersBySubaccountByServicePlan("uut", "integration-test-services-static", "TF-TEST-MALWARE-SCANNER", "clamav", "malware-scanner"),
				},
				{
					ResourceName:      "btp_subaccount_service_instance.uut",
					ImportStateIdFunc: getServiceInstanceIdForImportNoServiceInstanceId("btp_subaccount_service_instance.uut"),
					ImportState:       true,
					ImportStateVerify: true,
					ExpectError:       regexp.MustCompile(`Unexpected Import Identifier`),
				},
			},
		})
	})
}

func hclResourceSubaccountServiceInstanceWoParametersBySubaccountByServicePlan(resourceName string, subaccountName string, name string, servicePlanName string, serviceOfferingName string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		data "btp_subaccount_service_plans" "all" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
		}
		data "btp_subaccount_service_offering" "so" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			name =  "%[5]s"
		}
		resource "btp_subaccount_service_instance" "%[1]s"{
		    subaccount_id    = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			name             = "%[3]s"
		    serviceplan_id   = [for ssp in data.btp_subaccount_service_plans.all.values : ssp.id if ssp.name == "%[4]s" && ssp.serviceoffering_id == data.btp_subaccount_service_offering.so.id][0]
		}`, resourceName, subaccountName, name, servicePlanName, serviceOfferingName)
}

func hclResourceSubaccountServiceInstanceWoParametersWithTimeoutsBySubaccountByServicePlan(resourceName string, subaccountName string, name string, servicePlanName string, serviceOfferingName string, createTimeout string, updateTimeout string, deleteTimeout string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		data "btp_subaccount_service_plans" "all" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
		}
		data "btp_subaccount_service_offering" "so" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			name =  "%[5]s"
		}
		resource "btp_subaccount_service_instance" "%[1]s"{
		    subaccount_id    = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			name             = "%[3]s"
			serviceplan_id   = [for ssp in data.btp_subaccount_service_plans.all.values : ssp.id if ssp.name == "%[4]s" && ssp.serviceoffering_id == data.btp_subaccount_service_offering.so.id][0]
			timeouts = {
				create = "%[6]s"
				update = "%[7]s"
				delete = "%[8]s"
			  }
		}`, resourceName, subaccountName, name, servicePlanName, serviceOfferingName, createTimeout, updateTimeout, deleteTimeout)
}

func hclResourceSubaccountServiceInstanceWithParametersBySubaccountByServicePlan(resourceName string, subaccountName string, name string, servicePlanName string, serviceOfferingName string) string {
	xsuaaParameters := XsuaaParameters{
		Xsappname:  "bookshop-demo",
		TenantMode: "dedicated",
	}
	xsuaaParametersJson, _ := json.Marshal(xsuaaParameters)

	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		data "btp_subaccount_service_plans" "all" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
		}
		data "btp_subaccount_service_offering" "so" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			name =  "%[5]s"
		}
		resource "btp_subaccount_service_instance" "%[1]s"{
		    subaccount_id    = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			name             = "%[3]s"
			serviceplan_id   = [for ssp in data.btp_subaccount_service_plans.all.values : ssp.id if ssp.name == "%[4]s" && ssp.serviceoffering_id == data.btp_subaccount_service_offering.so.id][0]
			parameters       = %[6]q
		}`, resourceName, subaccountName, name, servicePlanName, serviceOfferingName, string(xsuaaParametersJson))
}

func hclResourceSubaccountServiceInstanceNoSubaccountId(resourceName string, name string, servicePlanId string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_instance" "%s"{
		    name             = "%s"
			serviceplan_id   = "%s"
		}`, resourceName, name, servicePlanId)
}

func hclResourceSubaccountServiceInstanceNoServicName(resourceName string, subaccountId string, servicePlanId string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_instance" "%s"{
		    subaccount_id    = "%s"
			serviceplan_id   = "%s"
		}`, resourceName, subaccountId, servicePlanId)
}

func hclResourceSubaccountServiceInstanceNoPlan(resourceName string, subaccountId string, name string) string {

	return fmt.Sprintf(`
		resource "btp_subaccount_service_instance" "%s"{
		    subaccount_id    = "%s"
			name             = "%s"
		}`, resourceName, subaccountId, name)
}

func hclResourceSubaccountServiceInstanceWithLabelsBySubaccountByPlan(resourceName string, subaccountName string, name string, servicePlanName string, serviceOfferingName string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		data "btp_subaccount_service_plans" "all" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
		}
		data "btp_subaccount_service_offering" "so" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			name =  "%[5]s"
		}
		resource "btp_subaccount_service_instance" "%[1]s" {
		    subaccount_id    = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			name             = "%[3]s"
			serviceplan_id   = [for ssp in data.btp_subaccount_service_plans.all.values : ssp.id if ssp.name == "%[4]s" && ssp.serviceoffering_id == data.btp_subaccount_service_offering.so.id][0]
			labels           = {"foo" = ["bar"]}
		}`, resourceName, subaccountName, name, servicePlanName, serviceOfferingName)
}

func hclResourceSubaccountServiceInstanceWithLabelsChangedBySubaccountByPlan(resourceName string, subaccountName string, name string, servicePlanName string, serviceOfferingName string) string {

	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		data "btp_subaccount_service_plans" "all" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
		}
		data "btp_subaccount_service_offering" "so" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			name =  "%[5]s"
		}
		resource "btp_subaccount_service_instance" "%[1]s"{
		    subaccount_id    = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			name             = "%[3]s"
			serviceplan_id   = [for ssp in data.btp_subaccount_service_plans.all.values : ssp.id if ssp.name == "%[4]s" && ssp.serviceoffering_id == data.btp_subaccount_service_offering.so.id][0]
			labels           = {"foo" = ["BAR"], "bar" = ["foo"]}
		}`, resourceName, subaccountName, name, servicePlanName, serviceOfferingName)
}

func getServiceInstanceIdForImport(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["subaccount_id"], rs.Primary.ID), nil
	}
}

func getServiceInstanceIdForImportNoServiceInstanceId(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s", rs.Primary.Attributes["subaccount_id"]), nil
	}
}
