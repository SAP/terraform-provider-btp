package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestFunctionExtractCfOrgId(t *testing.T) {

	t.Parallel()
	// Test happy case only, the error handling is mostly covered in function_helper_extract_environment_label_test.go
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: getProviders(nil),
		Steps: []resource.TestStep{
			{
				Config: hclFunctionExtractCfOrgId(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("cf_org_id", knownvalue.StringExact("8d818824-394a-4bae-9088-7a3c8ce93e57")),
				},
			},
		},
	})

}

func hclFunctionExtractCfOrgId() string {
	return `output "cf_org_id" {
		value = provider::btp::extract_cf_org_id("{\"API Endpoint\":\"https://api.cf.example.com\",\"Org Name\":\"cf-terraform-test\",\"Org ID\":\"8d818824-394a-4bae-9088-7a3c8ce93e57\",\"Org Memory Limit\":\"0MB\"}")
  }`
}
