package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountDestinationTrust(t *testing.T) {
	var regexpBase64 = regexp.MustCompile(`^[A-Za-z0-9+/=\s]+$`)
	var regexpX509 = regexp.MustCompile(`-----BEGIN CERTIFICATE-----[\s\S]+-----END CERTIFICATE-----`)

	t.Parallel()
	t.Run("happy path - active destination trust", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_trust_active")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountDestinationTrust("test", "integration-test-acc-static", "ACTIVE"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_trust.test", "trust_type", "ACTIVE"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_trust.test", "base_url", "cfapps.eu12.hana.ondemand.com"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_destination_trust.test", "expiration"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_destination_trust.test", "generated_on"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_trust.test", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_trust.test", "name", "keypair"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_trust.test", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_trust.test", "owner.subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_trust.test", "public_key_base64", regexpBase64),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_trust.test", "x509_public_key_base64", regexpX509),
					),
				},
			},
		})
	})

	t.Run("happy path - passive destination trust", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_trust_passive")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountDestinationTrust("test", "integration-test-acc-static", "PASSIVE"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_trust.test", "trust_type", "PASSIVE"),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_trust.test", "base_url", "cfapps.eu12.hana.ondemand.com"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_destination_trust.test", "expiration"),
						resource.TestCheckResourceAttrSet("data.btp_subaccount_destination_trust.test", "generated_on"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_trust.test", "id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_destination_trust.test", "name", "keypair"),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_trust.test", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_trust.test", "owner.subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_trust.test", "public_key_base64", regexpBase64),
						resource.TestMatchResourceAttr("data.btp_subaccount_destination_trust.test", "x509_public_key_base64", regexpX509),
					),
				},
			},
		})
	})

	t.Run("error path - subaccount is required", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_trust_subaccount_required")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + `
						data "btp_subaccount_destination_trust" "test" {
						  active = false
						  }`,
					ExpectError: regexp.MustCompile(`The argument \"subaccount_id\" is required`),
				},
			},
		})
	})

	t.Run("error path - destination trust doesn't exist", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_destination_trust_not_found")
		defer stopQuietly(rec)
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclDatasourceSubaccountDestinationTrust("test", "integration-test-services-static", "PASSIVE"),
					ExpectError: regexp.MustCompile(`There is no such key pair for this account`), // TODO improve error text
				},
			},
		})
	})

}

func hclDatasourceSubaccountDestinationTrust(resourceName string, Name string, TrustType string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_destination_trust" "%s" {	
subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
trust_type = "%s"
}`
	return fmt.Sprintf(template, resourceName, Name, TrustType)
}
