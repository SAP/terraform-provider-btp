package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountDestinationCertificate(t *testing.T) {
	t.Parallel()

	t.Run("happy path - destination certificate - PEM", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/data_source_subaccount_destination_certificate.read_pem_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDataSourceSubaccountDestinationCertificate("uut", "integration-test-destination", "terraform.pem"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_name", "terraform.pem"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "x509_certificate"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination certificate from service instance - PEM", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/data_source_subaccount_destination_certificate.service_instance.read_pem_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDataSourceSubaccountDestinationCertificateFromServiceInstance("uut", "integration-test-destination", "servtest", "terraform.pem"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_name", "terraform.pem"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "x509_certificate"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination certificate - P12", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/data_source_subaccount_destination_certificate.read_p12_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDataSourceSubaccountDestinationCertificate("uut", "integration-test-destination", "terraform.p12"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_name", "terraform.p12"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "private_key"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.format", "PKCS#8"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.algorithm", "RSA"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.1.type", "x509_certificate"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certification_creation_details.generation_method", "import"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination certificate from service instance - P12", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/data_source_subaccount_destination_certificate.service_instance.read_p12_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDataSourceSubaccountDestinationCertificateFromServiceInstance("uut", "integration-test-destination", "servtest", "terraform.p12"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_name", "terraform.p12"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "private_key"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.format", "PKCS#8"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.algorithm", "RSA"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.1.type", "x509_certificate"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certification_creation_details.generation_method", "import"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination certificate - PFX", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/data_source_subaccount_destination_certificate.read_pfx_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDataSourceSubaccountDestinationCertificate("uut", "integration-test-destination", "terraform.pfx"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_name", "terraform.pfx"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "private_key"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.format", "PKCS#8"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.algorithm", "RSA"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.1.type", "x509_certificate"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certification_creation_details.generation_method", "import"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination certificate from service instance - PFX", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/data_source_subaccount_destination_certificate.service_instance.read_pfx_certificate")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDataSourceSubaccountDestinationCertificateFromServiceInstance("uut", "integration-test-destination", "servtest", "terraform.pfx"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_name", "terraform.pfx"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "private_key"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.format", "PKCS#8"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.0.algorithm", "RSA"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certificate_nodes.1.type", "x509_certificate"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_certificate.uut", "certification_creation_details.generation_method", "import"),
					),
				},
			},
		})
	})
}

func hclDataSourceSubaccountDestinationCertificate(resourceName string, subaccountName string, certificateName string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		data "btp_subaccount_destination_certificate" "%s" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
			certificate_name = "%s" 
		}
	`, resourceName, subaccountName, certificateName)
}

func hclDataSourceSubaccountDestinationCertificateFromServiceInstance(resourceName string, subaccountName string, serviceInstanceName string, certificateName string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		data "btp_subaccount_service_instances" "all" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
		}
		data "btp_subaccount_destination_certificate" "%[1]s" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			service_instance_id  = [for sa in data.btp_subaccount_service_instances.all.values : sa.id if sa.name == "%[3]s"][0]
			certificate_name = "%[4]s" 
		}
	`, resourceName, subaccountName, serviceInstanceName, certificateName)
}
