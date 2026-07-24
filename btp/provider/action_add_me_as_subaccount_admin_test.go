package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestActionAddMeAsSubaccountAdmin(t *testing.T) {
	t.Parallel()
	t.Run("happy path - successful admin assignment", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/action_add_me_as_subaccount_admin")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclActionAddMeAsSubaccountAdmin("77395f6a-a601-4c9e-8cd0-c1fcefc7f60f"), //integration-test-acc-static
				},
			},
		})
	})

	t.Run("error path - non-existing subaccount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/action_add_me_as_subaccount_admin_fail_non_existing")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclActionAddMeAsSubaccountAdmin("00000000-0000-0000-0000-000000000000"),
					ExpectError: regexp.MustCompile(`API Error Adding Current User as Subaccount Admin`),
				},
			},
		})
	})
}

func hclActionAddMeAsSubaccountAdmin(subaccountId string) string {
	return fmt.Sprintf(`
resource "terraform_data" "test" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.btp_add_me_as_subaccount_admin.test]
    }
  }
}

action "btp_add_me_as_subaccount_admin" "test" {
  config {
    subaccount_id = "%s"
  }
}
`, subaccountId)
}
