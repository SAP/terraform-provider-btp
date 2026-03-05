package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

func TestResourceDisasterRecoverySubaccountPair(t *testing.T) {
	t.Run("happy path - create subaccount pair", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_disaster_recovery_subaccount_pair")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDisasterRecoverySubaccountPair("uut", "f59b5902-d24c-446c-b245-92c814faa0d9", "2dc1ecf1-786c-4f92-91f2-26650ab3ad28"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_disaster_recovery_subaccount_pair.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_disaster_recovery_subaccount_pair.uut", "paired_subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_disaster_recovery_subaccount_pair.uut", "pair_id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_disaster_recovery_subaccount_pair.uut", "created_date", regexpValidRFC3999Format),
						resource.TestCheckResourceAttrSet("btp_disaster_recovery_subaccount_pair.uut", "created_by"),
						resource.TestMatchResourceAttr("btp_disaster_recovery_subaccount_pair.uut", "globalaccount_id", regexpValidUUID),
					),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectIdentity(
							"btp_disaster_recovery_subaccount_pair.uut",
							map[string]knownvalue.Check{
								"subaccount_id":        knownvalue.StringRegexp(regexpValidUUID),
								"paired_subaccount_id": knownvalue.StringRegexp(regexpValidUUID),
							},
						),
					},
				},
				{
					ResourceName:    "btp_disaster_recovery_subaccount_pair.uut",
					ImportState:     true,
					ImportStateKind: resource.ImportBlockWithResourceIdentity,
				},
			},
		})
	})

	t.Run("error path - create subaccount pair existing", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_disaster_recovery_subaccount_pair.error_existing")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclResourceDisasterRecoverySubaccountPair("uut", "f59b5902-d24c-446c-b245-92c814faa0d9", "badbcbf8-eca7-4472-9f66-cb9887ba7c3d"),
					ExpectError: regexp.MustCompile(`INVALID_INPUT: subaccount`),
				},
			},
		})
	})
}

func hclResourceDisasterRecoverySubaccountPair(resourceName string, subaccountId string, pairedSubaccountId string) string {
	return fmt.Sprintf(`resource "btp_disaster_recovery_subaccount_pair" "%s" {
        subaccount_id          = "%s"
        paired_subaccount_id   = "%s"
    }`, resourceName, subaccountId, pairedSubaccountId)
}
