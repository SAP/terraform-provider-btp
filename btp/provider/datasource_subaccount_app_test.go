package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountApp(t *testing.T) {

	t.Parallel()
	t.Run("happy path - app by id", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccount_app.by_id")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccountAppByAppId("uut", "integration-test-services-static", "cas-ui-xsuaa-prod!t216"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("data.btp_subaccount_app.uut", "subaccount_id", regexpValidUUID),
						resource.TestCheckResourceAttr("data.btp_subaccount_app.uut", "id", "cas-ui-xsuaa-prod!t216"),
						resource.TestCheckResourceAttr("data.btp_subaccount_app.uut", "xsappname", "cas-ui-xsuaa-prod"),
						resource.TestCheckResourceAttr("data.btp_subaccount_app.uut", "plan_name", "application"),
						resource.TestCheckResourceAttr("data.btp_subaccount_app.uut", "plan_id", "ThGdx5loQ6XhvcdY6dLlEXcTgQD7641pDKXJfzwYGLg="),
						resource.TestCheckResourceAttr("data.btp_subaccount_app.uut", "tenant_mode", "shared"),
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
					Config:      hclDatasourceSubaccountAppNoAppId("uut", "integration-test-services-static"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found`),
				},
			},
		})
	})
	t.Run("error path - subaccount_id is required", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountAppNoSubaccountId("uut", "cas-ui-xsuaa-prod!t216"),
					ExpectError: regexp.MustCompile(`The argument "subaccount_id" is required, but no definition was found`),
				},
			},
		})
	})

}

func hclDatasourceSubaccountAppByAppId(resourceName string, subaccountName string, appId string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_app" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
	id 	          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountName, appId)
}

func hclDatasourceSubaccountAppNoSubaccountId(resourceName string, appId string) string {
	template := `data "btp_subaccount_app" "%s" {
	id 	          = "%s"
}`
	return fmt.Sprintf(template, resourceName, appId)
}

func hclDatasourceSubaccountAppNoAppId(resourceName string, subaccountName string) string {
	template := `
data "btp_subaccounts" "all" {}
data "btp_subaccount_app" "%s" {
	subaccount_id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
}`
	return fmt.Sprintf(template, resourceName, subaccountName)
}
