package provider
/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceGlobalaccountResourceProvider(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_globalaccount_resource_provider")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceGlobalaccountResourceProvider("uut",
						"AWS",
						"tf_test_resource_provider"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_provider.uut", "provider_type", "AWS"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_provider.uut", "technical_name", "tf_test_resource_provider"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_provider.uut", "id", "tf_test_resource_provider"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_provider.uut", "display_name", "Test AWS Resource Provider"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_provider.uut", "description", "Description of the resource provider"),
						resource.TestCheckResourceAttr("data.btp_globalaccount_resource_provider.uut", "configuration", "{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"eu-central-1\"}"),
					),
				},
			},
		})
	})
}

func hclDatasourceGlobalaccountResourceProvider(resourceName string, provider string, technicalName string) string {
	return fmt.Sprintf(`
data "btp_globalaccount_resource_provider" "%s" {
	provider_type = "%s"
	technical_name = "%s"
}`, resourceName, provider, technicalName)
}
*/