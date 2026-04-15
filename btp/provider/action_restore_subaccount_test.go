package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestActionRestoreSubaccount(t *testing.T) {
	t.Parallel()
	t.Run("happy path - successful restore", func(t *testing.T) {
		/*
			ATTENTION:
			BEFORE recording the test, make sure that the subaccount is in the state pending deletion
			AFTER recording the test, make sure that the subaccount is active
		*/
		rec, user := setupVCR(t, "fixtures/action_restore_subaccount")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclActionRestoreSubaccount("integration-test-pending-deletion"),
					// Post-apply is not needed as Action checks post conditions internally!
				},
			},
		})
	})

	t.Run("error path - restore active subaccount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/action_restore_subaccount_fail_active")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclActionRestoreSubaccount("integration-test-services-static"),
					ExpectError: regexp.MustCompile(`No pending deletion: The subaccount with ID`),
				},
			},
		})
	})

	t.Run("error path - restore non-existing subaccount", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/action_restore_subaccount_fail_non_existing")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclActionRestoreNonExistingSubaccount(),
					ExpectError: regexp.MustCompile(`API Error Reading Subaccount:`),
				},
			},
		})
	})
}

/*
IMPORTANT: Using a lifecycle.action_trigger block with config mode is currently the recommended way to test an action.
If there isn't a managed resource that fits to include in your acceptance test of an action,
the built-in terraform_data resource can be used as a replacement.
See: https://developer.hashicorp.com/terraform/plugin/framework/actions/testing
*/
func hclActionRestoreSubaccount(subaccountName string) string {
	return fmt.Sprintf(`
data "btp_subaccounts" "all" {}

resource "terraform_data" "test" {
  depends_on = [ data.btp_subaccounts.all ]
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.btp_restore_subaccount.test]
    }
  }
}

action "btp_restore_subaccount" "test" {
  config {
    subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
  }
}
`, subaccountName)
}

func hclActionRestoreNonExistingSubaccount() string {
	return `
resource "terraform_data" "test" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.btp_restore_subaccount.test]
    }
  }
}

action "btp_restore_subaccount" "test" {
  config {
    subaccount_id = "00000000-0000-0000-0000-000000000000"
  }
}`
}
