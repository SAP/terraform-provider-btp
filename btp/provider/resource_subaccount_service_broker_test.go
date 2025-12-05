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

	t.Run("error path - mtls true with cert and key provided", func(t *testing.T) {

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclResourceSubaccountServiceBrokerWithMTLS("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-broker", "a description", "https://my-broker-bogus-ratel-yb.cfapps.eu12.hana.ondemand.com", true, "dummy-cert", "dummy-key"),
					ExpectError: regexp.MustCompile("When `mtls` is true, `cert` and `key` must NOT be provided."),
				},
			},
		})
	})

	t.Run("error path - mtls set to false without basic auth", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_subaccount_service_broker_error")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclResourceSubaccountServiceBrokerWithoutBasicAuth("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-broker", "a description", "https://my-broker-bogus-ratel-yb.cfapps.eu12.hana.ondemand.com"),
					ExpectError: regexp.MustCompile("Invalid Configuration"),
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
			name		  = "%s"
			description   = "%s"

			url			  = "%s"
			username	  = "%s"
			password      = "%s"
		}
	`, resourceName, subaccountId, name, description, url, username, password)
}

func hclResourceSubaccountServiceBrokerWithoutBasicAuth(resourceName string, subaccountId string, name string, description string, url string) string {
	return fmt.Sprintf(`
		resource "btp_subaccount_service_broker" "%s" {
			subaccount_id = "%s"
			name		  = "%s"
			description   = "%s"

			url			  = "%s"
		}
	`, resourceName, subaccountId, name, description, url)
}

func hclResourceSubaccountServiceBrokerWithMTLS(resourceName string, subaccountId string, name string, description string, url string, mtls bool, cert string, key string) string {
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
