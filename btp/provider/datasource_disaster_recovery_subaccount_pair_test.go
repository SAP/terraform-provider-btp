package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceDisasterRecoverySubaccountPair(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_disaster_recovery_subaccount_pair")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{ // normal directory
					Config: hclProviderFor(user) + hclDatasourceDisasterRecoverySubaccountPair("uut", "integration-test-dr-subaccount-eu10-canary"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_disaster_recovery_subaccount_pair.uut", "pair_id", regexpValidUUID),
						resource.TestCheckResourceAttrSet("data.btp_disaster_recovery_subaccount_pair.uut", "created_by"),
						resource.TestMatchResourceAttr("data.btp_disaster_recovery_subaccount_pair.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("data.btp_disaster_recovery_subaccount_pair.uut", "globalaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_disaster_recovery_subaccount_pair.uut", "subaccounts.#", "2"),
						resource.TestMatchResourceAttr("data.btp_disaster_recovery_subaccount_pair.uut", "subaccounts.0.id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_disaster_recovery_subaccount_pair.uut", "subaccounts.1.id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_disaster_recovery_subaccount_pair.uut", "subaccounts.0.region", "eu10-canary"),
						resource.TestCheckResourceAttr("data.btp_disaster_recovery_subaccount_pair.uut", "subaccounts.1.region", "eu12"),
					),
				},
			},
		})
	})

	t.Run("error path - invalid subaccount ID", func(t *testing.T) {
		// See: https://github.com/SAP/terraform-provider-btp/issues/1210
		rec, user := setupVCR(t, "fixtures/datasource_disaster_recovery_subaccount_pair.invalid_id")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{ // normal directory
					Config:      hclProviderFor(user) + hclDatasourceDisasterRecoverySubaccountPairById("uut", "00000000-0000-0000-0000-000000000000"),
					ExpectError: regexp.MustCompile(`An error was encountered when reading the subaccount pair with ID`),
				},
			},
		})
	})

	t.Run("error path - subaccount id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      `data "btp_disaster_recovery_subaccount_pair" "uut" {}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found.`),
				},
			},
		})
	})
}

func hclDatasourceDisasterRecoverySubaccountPair(resourceName string, subaccountName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_disaster_recovery_subaccount_pair" "%s" {
    subaccount_id = [for sub in data.btp_subaccounts.all.values : sub.id if sub.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName)
}

func hclDatasourceDisasterRecoverySubaccountPairById(resourceName string, subaccountId string) string {
	template := `
data "btp_disaster_recovery_subaccount_pair" "%s" {
    subaccount_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId)
}
