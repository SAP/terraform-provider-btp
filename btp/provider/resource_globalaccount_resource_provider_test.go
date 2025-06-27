package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceGlobalaccountResourceProvider(t *testing.T) {
	t.Parallel()
	t.Run("happy path - create", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_resource_provider.create")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountResourceProvider("uut",
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
					Config: hclProviderFor(user) + hclResourceGlobalaccountResourceProviderNoDesc("uut",
						"AWS",
						"another_aws_resource_provider",
						"Another AWS Resource Provider",
						"{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"eu-central-1\"}",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "provider_type", "AWS"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "technical_name", "another_aws_resource_provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "display_name", "Another AWS Resource Provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "configuration", "{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"eu-central-1\"}"),
					),
				},
			},
		})
	})

	t.Run("happy path - update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_resource_provider.update")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountResourceProvider("uut",
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
					Config: hclProviderFor(user) + hclResourceGlobalaccountResourceProvider("uut",
						"AWS",
						"my_aws_resource_provider",
						"My New Display Name",
						"",
						"{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"us-east-1\"}",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "provider_type", "AWS"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "technical_name", "my_aws_resource_provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "display_name", "My New Display Name"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "description", ""),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "configuration", "{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"us-east-1\"}"),
					),
				},
			},
		})
	})

	t.Run("happy path - update omitting description", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_globalaccount_resource_provider.update_wo_description")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceGlobalaccountResourceProvider("uut",
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
					Config: hclProviderFor(user) + hclResourceGlobalaccountResourceProviderNoDesc("uut",
						"AWS",
						"my_aws_resource_provider",
						"My New Display Name",
						"{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"us-east-1\"}",
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "provider_type", "AWS"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "technical_name", "my_aws_resource_provider"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "display_name", "My New Display Name"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "description", "My description"),
						resource.TestCheckResourceAttr("btp_globalaccount_resource_provider.uut", "configuration", "{\"access_key_id\":\"AWSACCESSKEY\",\"secret_access_key\":\"AWSSECRETKEY\",\"vpc_id\":\"vpc-test\",\"region\":\"us-east-1\"}"),
					),
				},
			},
		})
	})

	t.Run("error path - provider_type is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: `
resource "btp_globalaccount_resource_provider" "uut" {
	technical_name = "technical_name"
	display_name = "display_name"
	configuration = "configuration"
}`,
					ExpectError: regexp.MustCompile(`The argument "provider_type" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - technical_name is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: `
resource "btp_globalaccount_resource_provider" "uut" {
	provider_type = "provider_type"
	display_name = "display_name"
	configuration = "configuration"
}`,
					ExpectError: regexp.MustCompile(`The argument "technical_name" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - display_name is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: `
resource "btp_globalaccount_resource_provider" "uut" {
	provider_type = "provider_type"
	technical_name = "technical_name"
	configuration = "configuration"
}`,
					ExpectError: regexp.MustCompile(`The argument "display_name" is required, but no definition was found.`),
				},
			},
		})
	})

	t.Run("error path - configuration is mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: `
resource "btp_globalaccount_resource_provider" "uut" {
	provider_type = "provider_type"
	technical_name = "technical_name"
	display_name = "display_name"
}`,
					ExpectError: regexp.MustCompile(`The argument "configuration" is required, but no definition was found.`),
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
