package provider

import (
	"fmt"
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
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBroker("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-broker", "a description", "https://my-broker-bogus-ratel-yb.cfapps.eu12.hana.ondemand.com", "platform", "a-secure-password"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "name", "my-broker"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "description", "a description"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "url", "https://my-broker-bogus-ratel-yb.cfapps.eu12.hana.ondemand.com"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "username", "platform"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "password", "a-secure-password"),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "last_modified", regexpValidRFC3999Format),
					),
				},
				{ // rename and update the description
					Config: hclProviderFor(user) + hclResourceSubaccountServiceBroker("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "my-broker-with-a-new-name", "another description", "https://my-broker-bogus-ratel-yb.cfapps.eu12.hana.ondemand.com", "platform", "a-secure-password"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_subaccount_service_broker.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "name", "my-broker-with-a-new-name"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "description", "another description"),
						resource.TestCheckResourceAttr("btp_subaccount_service_broker.uut", "url", "https://my-broker-bogus-ratel-yb.cfapps.eu12.hana.ondemand.com"),
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
					ImportStateVerifyIgnore: []string{"name", "username", "password"},
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
