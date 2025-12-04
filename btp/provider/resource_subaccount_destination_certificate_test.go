package provider

import (
	"fmt"
	"testing"

	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceSubaccountDestinationCertificate(t *testing.T) {

	t.Run("happy path - destination certificate - PEM", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_certificate.create_pem_certificate")
		defer stopQuietly(rec)

		certContent := "redacted"

		if rec.IsRecording() {
			var err error
			certContent, err = tfutils.GetBase64EncodedCertificate("pem")
			if err != nil {
				t.Fatalf("Unable to generate certificate content: %s", err)
			}
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationCertificate("uut", "integration-test-destination", "cert.pem", certContent),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_name", "cert.pem"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "x509_certificate"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certification_creation_details.generation_method", "import"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination certificate from service instance - PEM", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_certificate.service_instance.create_pem_certificate")
		defer stopQuietly(rec)

		certContent := "redacted"

		if rec.IsRecording() {
			var err error
			certContent, err = tfutils.GetBase64EncodedCertificate("pem")
			if err != nil {
				t.Fatalf("Unable to generate certificate content: %s", err)
			}
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationCertificateFromServiceInstance("uut", "integration-test-destination", "servtest", "cert.pem", certContent),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_name", "cert.pem"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "x509_certificate"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certification_creation_details.generation_method", "import"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination certificate - P12", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_certificate.create_p12_certificate")
		defer stopQuietly(rec)

		certContent := "redacted"

		if rec.IsRecording() {
			var err error
			certContent, err = tfutils.GetBase64EncodedCertificate("p12")
			if err != nil {
				t.Fatalf("Unable to generate certificate content: %s", err)
			}
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationCertificate("uut", "integration-test-destination", "cert.p12", certContent),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_name", "cert.p12"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "private_key"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.format", "PKCS#8"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.algorithm", "RSA"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.1.type", "x509_certificate"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certification_creation_details.generation_method", "import"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination certificate from service instance - P12", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_certificate.service_instance.create_p12_certificate")
		defer stopQuietly(rec)

		certContent := "redacted"

		if rec.IsRecording() {
			var err error
			certContent, err = tfutils.GetBase64EncodedCertificate("p12")
			if err != nil {
				t.Fatalf("Unable to generate certificate content: %s", err)
			}
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationCertificateFromServiceInstance("uut", "integration-test-destination", "servtest", "cert.p12", certContent),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_name", "cert.p12"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "private_key"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.format", "PKCS#8"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.algorithm", "RSA"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.1.type", "x509_certificate"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certification_creation_details.generation_method", "import"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination certificate - PFX", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_certificate.create_pfx_certificate")
		defer stopQuietly(rec)

		certContent := "redacted"

		if rec.IsRecording() {
			var err error
			certContent, err = tfutils.GetBase64EncodedCertificate("pfx")
			if err != nil {
				t.Fatalf("Unable to generate certificate content: %s", err)
			}
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationCertificate("uut", "integration-test-destination", "cert.pfx", certContent),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_name", "cert.pfx"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "private_key"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.format", "PKCS#8"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.algorithm", "RSA"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.1.type", "x509_certificate"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certification_creation_details.generation_method", "import"),
					),
				},
			},
		})
	})

	t.Run("happy path - destination certificate from service instance - PFX", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_destination_certificate.service_instance.create_pfx_certificate")
		defer stopQuietly(rec)

		certContent := "redacted"

		if rec.IsRecording() {
			var err error
			certContent, err = tfutils.GetBase64EncodedCertificate("pfx")
			if err != nil {
				t.Fatalf("Unable to generate certificate content: %s", err)
			}
		}

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountDestinationCertificateFromServiceInstance("uut", "integration-test-destination", "servtest", "cert.pfx", certContent),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_name", "cert.pfx"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.type", "private_key"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.format", "PKCS#8"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.0.algorithm", "RSA"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certificate_nodes.1.type", "x509_certificate"),
						resource.TestCheckResourceAttr("btp_subaccount_destination_certificate.uut", "certification_creation_details.generation_method", "import"),
					),
				},
			},
		})
	})
}

func hclResourceSubaccountDestinationCertificate(resourceName, subaccountId, certificateName, certificateContent string) string {
	return fmt.Sprintf(`
		data "btp_subaccounts" "all" {}
		resource "btp_subaccount_destination_certificate" "%s" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
			certificate_name     = "%s"
			certificate_content  = "%s"
		}
		`, resourceName, subaccountId, certificateName, certificateContent)
}

func hclResourceSubaccountDestinationCertificateFromServiceInstance(resourceName, subaccountId, serviceInstanceId string, certificateName, certificateContent string) string {
	return fmt.Sprintf(`
	    data "btp_subaccounts" "all" {}
		data "btp_subaccount_service_instances" "all" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
		}
		resource "btp_subaccount_destination_certificate" "%[1]s" {
			subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%[2]s"][0]
			service_instance_id  = [for sa in data.btp_subaccount_service_instances.all.values : sa.id if sa.name == "%[3]s"][0]
			certificate_name     = "%[4]s"
			certificate_content  = "%[5]s"
		}
		`, resourceName, subaccountId, serviceInstanceId, certificateName, certificateContent)
}
