package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestResourceSubaccountServiceBroker(t *testing.T) {
	t.Run("happy path - simple service_broker", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_broker")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBroker("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-dummy-test-broker", "a description", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com", "platform", "a-secure-password"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "name", "my-dummy-test-broker"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "description", "a description"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "url", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "username", "platform"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "password", "a-secure-password"),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{ // rename and update the description
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBroker("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-dummy-test-broker-new-name", "another description", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com", "platform", "a-secure-password"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "name", "my-dummy-test-broker-new-name"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "description", "another description"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "url", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "username", "platform"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "password", "a-secure-password"),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{
					ResourceName:            "btp_subaccount_service_broker.uut",
					ImportStateIdFunc:       getServiceBrokerIdForImport("btp_subaccount_service_broker.uut"),
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"name", "username", "password", "mtls", "cert", "key"},
				},
			},
		})
	})

	t.Run("happy path - service_broker with mtls only", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_broker_with_mtls")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBrokerWithMTLS("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-dummy-test-broker-with-mtls", "a description", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com", true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "name", "my-dummy-test-broker-with-mtls"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "description", "a description"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "url", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "mtls", "true"),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{ // rename and update the description
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBrokerWithMTLS("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-dummy-test-broker-with-mtls", "another description", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com", true),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "name", "my-dummy-test-broker-with-mtls"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "description", "another description"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "url", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "mtls", "true"),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{
					ResourceName:            "btp_subaccount_service_broker.uut",
					ImportStateIdFunc:       getServiceBrokerIdForImport("btp_subaccount_service_broker.uut"),
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"name", "username", "password", "mtls", "cert", "key"},
				},
			},
		})
	})

	t.Run("happy path - service_broker with mtls and username and password", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_broker_with_mtls_username_password")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBrokerWithMTLSUserPwd("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-dummy-test-broker-with-mtls-and-basic-auth", "a description", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com", true, "platform", "a-secure-password"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "name", "my-dummy-test-broker-with-mtls-and-basic-auth"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "description", "a description"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "url", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "mtls", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "username", "platform"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "password", "a-secure-password"),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{ // rename and update the description
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBrokerWithMTLSUserPwd("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-dummy-test-broker-with-mtls-and-basic-auth", "another description", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com", true, "platform", "a-secure-password"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "name", "my-dummy-test-broker-with-mtls-and-basic-auth"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "description", "another description"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "url", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "mtls", "true"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "username", "platform"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "password", "a-secure-password"),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{
					ResourceName:            "btp_subaccount_service_broker.uut",
					ImportStateIdFunc:       getServiceBrokerIdForImport("btp_subaccount_service_broker.uut"),
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"name", "username", "password", "mtls", "cert", "key"},
				},
			},
		})
	})

	t.Run("happy path - service_broker with cert and key", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_broker_with_cert_and_key")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBrokerWithCertKey("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-dummy-test-broker-with-mtls-and-basic-auth", "a description", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "name", "my-dummy-test-broker-with-mtls-and-basic-auth"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "description", "a description"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "url", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com"),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{ // rename and update the description
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBrokerWithCertKey("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-dummy-test-broker-with-mtls-and-basic-auth", "another description", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "name", "my-dummy-test-broker-with-mtls-and-basic-auth"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "description", "another description"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "url", "https://my-dummy-broker.cfapps.eu12.hana.ondemand.com"),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{
					ResourceName:            "btp_subaccount_service_broker.uut",
					ImportStateIdFunc:       getServiceBrokerIdForImport("btp_subaccount_service_broker.uut"),
					ImportState:             true,
					ImportStateVerify:       true,
					ImportStateVerifyIgnore: []string{"name", "username", "password", "mtls", "cert", "key"},
				},
			},
		})
	})

	t.Run("error path - mtls set to true with cert and key provided", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountServiceBrokerWithMTLSCertKey("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-broker", "a description", "https://my-broker-bogus-ratel-yb.cfapps.eu12.hana.ondemand.com", true, "dummy-cert", "dummy-key"),
					ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
				},
			},
		})
	})

	t.Run("error path - mtls not set and no auth data", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountServiceBrokerWithoutBasicAuth("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-broker", "a description", "https://my-broker-bogus-ratel-yb.cfapps.eu12.hana.ondemand.com"),
					ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
				},
			},
		})
	})

	t.Run("error path - mtls set to false with auth data", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountServiceBrokerWithMTLSUserPwd("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-broker", "a description", "https://my-broker-bogus-ratel-yb.cfapps.eu12.hana.ondemand.com", false, "platform", "a-secure-password"),
					ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
				},
			},
		})
	})
}

func getServiceBrokerIdForImport(resourceName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s,%s", rs.Primary.Attributes["subaccount_id"], rs.Primary.ID), nil
	}
}

func hclResourceSubaccountServiceBroker(resourceName string, subaccountId string, name string, description string, url string, username string, password string) string {
	return fmt.Sprintf(`
		resource "btp_subaccount_service_broker" "%s" {
			subaccount_id = "%s"
			name		      = "%s"
			description   = "%s"
			url			      = "%s"
			username	    = "%s"
			password      = "%s"
		}
	`, resourceName, subaccountId, name, description, url, username, password)
}

func hclResourceSubaccountServiceBrokerWithoutBasicAuth(resourceName string, subaccountId string, name string, description string, url string) string {
	return fmt.Sprintf(`
		resource "btp_subaccount_service_broker" "%s" {
			subaccount_id = "%s"
			name		      = "%s"
			description   = "%s"
			url			      = "%s"
		}
	`, resourceName, subaccountId, name, description, url)
}

func hclResourceSubaccountServiceBrokerWithMTLSCertKey(resourceName string, subaccountId string, name string, description string, url string, mtls bool, cert string, key string) string {
	return fmt.Sprintf(`
		resource "btp_subaccount_service_broker" "%s" {
			subaccount_id = "%s"
			name          = "%s"
			description   = "%s"
			url           = "%s"
			mtls          = %t
			cert          = "%s"
			key           = "%s"
		}
	`, resourceName, subaccountId, name, description, url, mtls, cert, key)
}

func hclResourceSubaccountServiceBrokerWithMTLSUserPwd(resourceName string, subaccountId string, name string, description string, url string, mtls bool, username string, password string) string {
	return fmt.Sprintf(`
		resource "btp_subaccount_service_broker" "%s" {
			subaccount_id = "%s"
			name          = "%s"
			description   = "%s"
			url           = "%s"
			mtls          = %t
			username      = "%s"
			password      = "%s"
		}
	`, resourceName, subaccountId, name, description, url, mtls, username, password)
}

func hclResourceSubaccountServiceBrokerWithMTLS(resourceName string, subaccountId string, name string, description string, url string, mtls bool) string {
	return fmt.Sprintf(`
		resource "btp_subaccount_service_broker" "%s" {
			subaccount_id = "%s"
			name          = "%s"
			description   = "%s"
			url           = "%s"
			mtls          = %t
		}
	`, resourceName, subaccountId, name, description, url, mtls)
}

func hclResourceSubaccountServiceBrokerWithCertKey(resourceName string, subaccountId string, name string, description string, url string) string {
	dummyCert := `-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIUOYCo4GLDaOj9NPGZ51hw6s6FU8AwDQYJKoZIhvcNAQEL
BQAwXjELMAkGA1UEBhMCREUxEzARBgNVBAgMClNvbWUtU3RhdGUxDTALBgNVBAoM
BFRlc3QxDTALBgNVBAsMBFRlc3QxHDAaBgkqhkiG9w0BCQEWDXRlc3RAdGVzdC5j
b20wHhcNMjUxMjA5MTA1NTQxWhcNMjYxMjA5MTA1NTQxWjBeMQswCQYDVQQGEwJE
RTETMBEGA1UECAwKU29tZS1TdGF0ZTENMAsGA1UECgwEVGVzdDENMAsGA1UECwwE
VGVzdDEcMBoGCSqGSIb3DQEJARYNdGVzdEB0ZXN0LmNvbTCCAiIwDQYJKoZIhvcN
AQEBBQADggIPADCCAgoCggIBAKc7zXzLAgNlWDGm4lAw757WiGgdOAObkWZaYU25
/TXXIGnb21p03Vf168uh/57RhTESz90++9f5rLkqPFyPrKezs2RPBiq61U6G8It2
wduQz5qeycIssIADkHA4M/OZZ7uL/uzDhpHyRV1Bf+owmmXdCIKvtt+UDIgZ4jKQ
P6SSDRhH4eEki4Zi5AQwJw4G33ocQ80fQhe6mYyOVOMMYUfoWL3uM+51VEFBlUwE
Ht2+/9YdlyX+/LGBrQNKv1qJIVXZtRBycq57yqfoV2NMUc915ThZReSS2XejFcW0
yMXtf1Isut6Vt0dV5ZevMRItyiy63hetIF+6m3M05muQDlbgQlgqcfr6OKMeZFWC
ufdVtbek2HzJ93l/H0Byl7xNUZXVTA2DnHRtmOVlRnzm4pT/z3ccE+XgFWfydnaG
j2izC+skmeuGRq5WBYpM4Nl+10GeWd10iAM6c7hSOmaZCfzFRZK+MsXGpXTdLpDs
fCfjtsywL/wyM5efLhGQIRN5grxggFf9ZX/mghD0ykcXDM/aSNHcKa4XZ3Umn9AY
QUUWVjJ05e/CRlbBwQ7us5bm4Q5bpz7979Q7rrwkh5TEXzmp21hR9OleMLl2aDLu
HBYIn6FqBXwJC8VooUKlX+9H/nGDArUWhOoy4yRomRpX5N+Rq5mw5UUl7yiPYZzf
Svo9AgMBAAGjITAfMB0GA1UdDgQWBBSjfQcYKQIDDWvRxpYalGBo39G8yzANBgkq
hkiG9w0BAQsFAAOCAgEAjmuExBa2/hkWJ0cwci08ddNQC8/zsdpxBvZmhXm06RW/
kWyNLXbvuzDk0gurAe7NwnuLKUFnP2OjNcx5mNwN/6VURvAvnqp/odI0NmUn6gXf
TjTq3zu22msa/ODZeU/FpAzLlkUAIcyReA3pXzOARszNUflYbdHSQhaNq5IPrJ0x
9jx3Jbdma7m3pFqsPsA33qK7VwEjEPDusBn9C2cDtKxZ/cb1MjEsjSt3a9KQzpku
HYiFiKmzPe3MCumpK5exuTA3fIi1Gkw4anI6kRssqcKSzTvwiaKRh7zbz/0Jc7WQ
9FUTVqygtfqSi3809hfA9tLHDgvqK9moAclvC3nS37FDaEDuWgbj9+XqnOzehf1F
m64lU90jLJsQowFpN6Jft+tHWgiNdC+g22VyvgTgy5bt0my6pKR5JRa2wTYvyY8f
GcrQETPv9NhD9naKVEWW/7vMpu7AZriPZ0mp30HBzNBUaTc6EeGq5mQgZ2dA4nA9
kD354t0c64+7P41FNsddJsaeHJvFPLjtxRqP1AEdhRlYS4kWH6EAgx3AEfeAlLAv
gc9GEi84+0Mq66VUaXdMwAHjBjQxlgfiJsXaaEAhaRoVYR9YDBfdoeWaybHK0GlF
OluDXl5F1QkPepfHf2p6MljYQgiqnXDMZLjXkE//UgZ6IqTH10a0dBpST5nPoP8=
-----END CERTIFICATE-----`

	dummyKey := `-----BEGIN PRIVATE KEY-----
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
vZL2qDLwQqC2e5n0/U7a1TJZg9QlvA==
-----END PRIVATE KEY-----`

	return fmt.Sprintf(`
		resource "btp_subaccount_service_broker" "%s" {
			subaccount_id = "%s"
			name          = "%s"
			description   = "%s"
			url           = "%s"
			cert          = <<CERT
%s
CERT
			key           = <<KEY
%s
KEY
		}
	`, resourceName, subaccountId, name, description, url, dummyCert, dummyKey)
}
