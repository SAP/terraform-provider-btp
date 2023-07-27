package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGlobalaccountResourceProvider(t *testing.T) {
	t.Parallel()
	t.Run("happy path - create", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_globalaccount_resource_provider.create")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceGlobalaccountResourceProvider("uut",
						"AWS",
						"my_aws_resource_provider",
						"My AWS Resource Provider",
						"My description",
						"{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"eu-central-1\"}",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "provider_type", "AWS"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "technical_name", "my_aws_resource_provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "id", "my_aws_resource_provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "display_name", "My AWS Resource Provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "description", "My description"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "configuration", "{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"eu-central-1\"}"),
					),
				},
				{
					Config: hclProvider() + hclResourceGlobalaccountResourceProviderNoDesc("uut",
						"AWS",
						"another_resource_provider",
						"Another Resource Provider",
						"{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"eu-central-1\"}",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "provider_type", "AWS"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "technical_name", "another_resource_provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "display_name", "Another Resource Provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "configuration", "{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"eu-central-1\"}"),
					),
				},
			},
		})
	})

}

func hclResourceGlobalaccountResourceProvider(resourceName string, provider string, technicalName string, displayName string, description string, configuration string) string {
	return fmt.Sprintf(`
resource "btp_globalaccount_resource_provider" "%s" {
	provider_type = "%s"
	technical_name = "%s"
	display_name = "%s"
	description = "%s"
	configuration = %q
}`, resourceName, provider, technicalName, displayName, description, configuration)
}

func hclResourceGlobalaccountResourceProviderNoDesc(resourceName string, provider string, technicalName string, displayName string, configuration string) string {
	return fmt.Sprintf(`
resource "btp_globalaccount_resource_provider" "%s" {
	provider_type = "%s"
	technical_name = "%s"
	display_name = "%s"
	configuration = %q
}`, resourceName, provider, technicalName, displayName, configuration)
}
