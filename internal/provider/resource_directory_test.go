package provider

/*
import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

*/
/*

func TestResourceDirectory(t *testing.T) {
	t.Run("happy path - parent directory", func(t *testing.T) {
		rec := setupVCR(t, "fixtures/resource_directory.parent_directory")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProvider() + hclResourceDirectoryParent("uut", "my-parent-folder", "This is a parent folder"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "parent_id", regexpValidUUID),

						resource.TestCheckResourceAttr("btp_directory.uut", "name", "my-parent-folder"),
						resource.TestCheckResourceAttr("btp_directory.uut", "description", "This is a parent folder"),
						resource.TestCheckResourceAttr("btp_directory.uut", "subdomain", ""),
						resource.TestCheckResourceAttr("btp_directory.uut", "created_by", "john.doe@int.test"),
						resource.TestCheckResourceAttr("btp_directory.uut", "state", "OK"),

						resource.TestCheckResourceAttr("btp_directory.uut", "features.#", "1"),
						resource.TestCheckResourceAttr("btp_directory.uut", "labels.#", "0"),
					),
				},
				{
					ResourceName:      "btp_directory.uut",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})
}

func hclResourceDirectoryParent(resourceName string, displayName string, description string) string {
	return fmt.Sprintf(`resource "btp_directory" "%s" {
        name        = "%s"
        description = "%s"
    }`, resourceName, displayName, description)
}

*/
