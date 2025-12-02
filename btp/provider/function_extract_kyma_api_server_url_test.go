package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestFunctionExtractKymaApiServerUrl(t *testing.T) {

	t.Parallel()
	// Test happy case only, the error handling is mostly covered in function_helper_extract_environment_label_test.go
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: getProviders(nil),
		Steps: []resource.TestStep{
			{
				Config: hclFunctionExtractKymaApiUrl(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("kyma_server_api_url", knownvalue.StringExact("https://api.a-123.kyma.ondemand.com")),
				},
			},
		},
	})

}

func hclFunctionExtractKymaApiUrl() string {
	return `output "kyma_server_api_url" {
		value = provider::btp::extract_kyma_api_server_url("{\"APIServerURL\":\"https://api.a-123.kyma.ondemand.com\",\"KubeconfigURL\":\"https://kyma-env-broker.cp.kyma.cloud.sap/kubeconfig/ABCDEB98-AABB-1234-ABAB-186489849D6F\",\"Name\":\"test-terraform-function-1234567\"}")
  }`
}
