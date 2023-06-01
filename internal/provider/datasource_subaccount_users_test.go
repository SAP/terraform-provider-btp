package provider

import (
	"fmt"
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
					Config: hclProvider() + hclDatasourceSubaccountUsersDefaultIdp("defaultidp", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccount_users.defaultidp", "subaccount_id", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
						resource.TestCheckResourceAttr("data.btp_subaccount_users.defaultidp", "values.#", "3"),
					),
				},
			},
		})
	})

	/* TBD: make this code work. Currently we get the following error message: cannot unmarshal object into Go value of type []string */
	/*
		t.Run("happy path with custom idp", func(t *testing.T) {
			rec := setupVCR(t, "fixtures/datasource_subaccount_users.custom_idp")
			defer stopQuietly(rec)

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
				Steps: []resource.TestStep{
					{
						Config: hclProvider() + hclDatasourceSubaccountUsersCustomIdp("mycustomidp", "ef23ace8-6ade-4d78-9c1f-8df729548bbf", "terraformint"),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.btp_subaccount_users.mycustomidp", "subaccount_id", "ef23ace8-6ade-4d78-9c1f-8df729548bbf"),
							resource.TestCheckResourceAttr("data.btp_subaccount_users.mycustomidp", "values.#", "3"),
						),
					},
				},
			})
		})
	*/

}

func hclDatasourceSubaccountUsersDefaultIdp(resourceName string, subaccountId string) string {
	template := `data "btp_subaccount_users" "%s" {
	subaccount_id = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId)
}

func hclDatasourceSubaccountUsersCustomIdp(resourceName string, subaccountId string, origin string) string {
	template := `data "btp_subaccount_users" "%s" {
	subaccount_id = "%s"
	origin = "%s"
}`
	return fmt.Sprintf(template, resourceName, subaccountId, origin)
}
