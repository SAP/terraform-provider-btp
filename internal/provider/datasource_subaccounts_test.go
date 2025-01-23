package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestDataSourceSubaccounts(t *testing.T) {
	t.Parallel()
	t.Run("happy path", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/datasource_subaccounts.all")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclDatasourceSubaccounts("uut"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.btp_subaccounts.uut", "values.#", "9"),
					),
				},
			},
		})
	})

}

func hclDatasourceSubaccounts(resourceName string) string {
	template := `data "btp_subaccounts" "%s" {}`
	return fmt.Sprintf(template, resourceName)
}
