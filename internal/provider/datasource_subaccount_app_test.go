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
					Config: hclProviderFor(user) + hclDatasourceSubaccountAppById("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5", "cas-ui-xsuaa-prod!t216"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_app.uut", "subaccount_id", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
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
					Config:      hclDatasourceSubaccountAppNoId("uut", "59cd458e-e66e-4b60-b6d8-8f219379f9a5"),
					ExpectError: regexp.MustCompile(`The argument "id" is required, but no definition was found`),
				},
			},
		})
	})
	t.Run("error path - subaccount_id not a valid UUID", func(t *testing.T) {
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

	t.Run("error path - subaccount_id not a valid UUID", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []resource.TestStep{
				{
					Config:      hclDatasourceSubaccountAppById("uut", "this-is-not-a-uuid", "cas-ui-xsuaa-prod!t216"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
}

func hclDatasourceSubaccountAppById(resourceName string, subaccountId string, appId string) string {
	template := `data "btp_subaccount_app" "%s" {
	subaccount_id = "%s"
	id 	          = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, appId)
}

func hclDatasourceSubaccountAppNoSubaccountId(resourceName string, appId string) string {
	template := `data "btp_subaccount_app" "%s" {
	id 	          = "%s"
}`
	return fmt.Sprintf(template, resourceName, appId)
}

func hclDatasourceSubaccountAppNoId(resourceName string, subaccountId string) string {
	template := `data "btp_subaccount_app" "%s" {
	subaccount_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId)
}
