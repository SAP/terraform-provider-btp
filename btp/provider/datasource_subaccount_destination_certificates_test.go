package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceSubaccountDestinationCertificates(t *testing.T) {

	t.Run("happy path - destination certificates", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/data_source_subaccount_destination_certificates.read_certificates")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationCertificates("uut", "integration-test-destination", "servtest"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificates.uut", "subaccount_level_certificates.#", "3"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificates.uut", "service_instance_level_certificates.#", "3"),
					),
				},
			},
		})
	})

}

func hclResourceSubaccountDestinationCertificates(resourceName, subaccountName, serviceInstanceName string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		data "btp_subaccount_service_instances" "all" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
		}
		data "btp_subaccount_destination_certificates" "%[1]s" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			service_instance_id  = [for sa in data.btp_subaccount_service_instances.all.values : sa.id if sa.name == "%[3]s"][0]
		}
		`, resourceName, subaccountName, serviceInstanceName)
}
