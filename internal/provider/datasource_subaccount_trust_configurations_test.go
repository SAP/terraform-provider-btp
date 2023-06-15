package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountTrustConfigurations(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_subaccount_trust_configurations.subaccount_exists")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceSubaccountTrustConfigurations("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_trust_configurations.uut", "values.#", "2"),
					),
				},
			},
		})
	})
	t.Run("happy path - subaccount not existing", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_subaccount_trust_configurations.subaccount_not_existing")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config:      hclProvider() + hclDatasourceSubaccountTrustConfigurations("uut", "aaaaaaaa-bbbb-cccc-dddd-caffee00affe"),
					ExpectError: regexp.MustCompile(`Received response with unexpected status \[Status: 404; Correlation ID:\s+[a-f0-9\-]+\]`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountTrustConfigurations(resourceName string, subaccountId string) string {
	template := `data "btp_subaccount_trust_configurations" "%s" { subaccount_id = "%s" }`

	return fmt.Sprintf(template, resourceName, subaccountId)
}
