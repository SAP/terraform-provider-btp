package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestActionUpdateGlobalAccount(t *testing.T) {
	t.Parallel()

	t.Run("happy path - update description", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/action_update_global_account")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclActionUpdateGlobalAccountDescription("Global Account for running tests of Terraform Provider for SAP BTP"),
				},
			},
		})
	})

	t.Run("happy path - enable subaccount force deletion", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/action_update_global_account_enable_force_deletion")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclActionUpdateGlobalAccountForceDeletion(true),
				},
			},
		})
	})

	t.Run("error path - no attributes provided", func(t *testing.T) {
		t.Parallel()
		rec, user := setupVCR(t, "fixtures/action_update_global_account_fail_no_params")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			TerraformVersionChecks: []tfversion.TerraformVersionCheck{
				tfversion.SkipBelow(tfversion.Version1_14_0),
			},
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProviderFor(user) + hclActionUpdateGlobalAccountNoParams(),
					ExpectError: regexp.MustCompile(`At least one of`),
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
func hclActionUpdateGlobalAccountDescription(description string) string {
	return fmt.Sprintf(`
resource "terraform_data" "test" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.btp_update_global_account.test]
    }
  }
}

action "btp_update_global_account" "test" {
  config {
    description = "%s"
  }
}
`, description)
}

func hclActionUpdateGlobalAccountForceDeletion(enable bool) string {
	return fmt.Sprintf(`
resource "terraform_data" "test" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.btp_update_global_account.test]
    }
  }
}

action "btp_update_global_account" "test" {
  config {
    enable_subaccount_force_deletion = %v
  }
}
`, enable)
}

func hclActionUpdateGlobalAccountNoParams() string {
	return `
resource "terraform_data" "test" {
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.btp_update_global_account.test]
    }
  }
}

action "btp_update_global_account" "test" {
  config {}
}`
}
