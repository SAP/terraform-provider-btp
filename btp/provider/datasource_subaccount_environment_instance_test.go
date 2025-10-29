package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountEnvironmentInstance(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_environment_instance")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountEnvironmentInstanceByInstanceId("uut", "integration-test-acc-static", "DA2883C7-0FAF-4D4A-80BB-A0B54AC9743D"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_environment_instance.uut", "subaccount_id", regexpValidUUID),
						resource.TestMatchResourceAttr("data.btp_subaccount_environment_instance.uut", "id", regexpValidUUID),
					),
				},
			},
		})
	})
	t.Run("error path - subaccount_id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: `
					data "btp_subaccount_environment_instance" "instance" {
						id = "DA2883C7-0FAF-4D4A-80BB-A0B54AC9743D"
					}`,
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - id mandatory", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config: `
					data "btp_subaccounts" "all" {}
					data "btp_subaccount_environment_instance" "instance" {
						subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "integration-test-acc-static"][0]
					}`,
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found`),
				},
			},
		})
	})

	t.Run("error path - id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountEnvironmentInstanceByInstanceId("uut", "integration-test-acc-static", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountEnvironmentInstanceByInstanceId(resourceName string, subaccountName string, instanceId string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_environment_instance" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	id = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, instanceId)
}
