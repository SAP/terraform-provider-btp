package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccountUsers(t *testing.T) {
	t.Parallel()
	t.Run("happy path with default idp", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_subaccount_users.default_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceSubaccountUsersDefaultIdp("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_users.uut", "subaccount_id", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
						resource.TestCheckResourceAttr("data.btp_subaccount_users.uut", "values.#", "3"),
					),
				},
			},
		})
	})
	t.Run("happy path with custom idp", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/datasource_subaccount_users.custom_idp")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclDatasourceSubaccountUsersWithCustomIdp("uut", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "terraformint-platform"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_users.uut", "subaccount_id", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
						resource.TestCheckResourceAttr("data.btp_subaccount_users.uut", "values.#", "1"),
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
					Config:      hclProvider() + `data "btp_subaccount_users" "uut" {}`,
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
					Config:      hclProvider() + hclDatasourceSubaccountUsersDefaultIdp("uut", "this-is-not-a-uuid"),
					ExpectError: regexp.MustCompile(`Attribute subaccount_id value must be a valid UUID, got: this-is-not-a-uuid`),
				},
			},
		})
	})
	// TODO: error path with non existing idp
}

func hclDatasourceSubaccountUsersDefaultIdp(resourceName string, subaccountId string) string {
	template := `data "btp_subaccount_users" "%s" {
	subaccount_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId)
}

func hclDatasourceSubaccountUsersWithCustomIdp(resourceName string, subaccountId string, origin string) string {
	template := `data "btp_subaccount_users" "%s" {
	subaccount_id = "%s"
	origin = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, origin)
}
