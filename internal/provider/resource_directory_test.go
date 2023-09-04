package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestResourceDirectory(t *testing.T) {
	t.Parallel()
	t.Run("happy path - parent directory", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectory("uut", "my-new-directory", "This is a new directory"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "parent_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory.uut", "name", "my-new-directory"),
						resource.TestCheckResourceAttr("btp_directory.uut", "description", "This is a new directory"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceDirectory("uut", "my-updated-directory", "This is a updated directory"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "parent_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory.uut", "name", "my-updated-directory"),
						resource.TestCheckResourceAttr("btp_directory.uut", "description", "This is a updated directory"),
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

	t.Run("happy path - directory with features", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory_features")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryWithFeatures("uut", "my-new-directory-feat", "This is a new directory with features"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "parent_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory.uut", "name", "my-new-directory-feat"),
						resource.TestCheckResourceAttr("btp_directory.uut", "description", "This is a new directory with features"),
						resource.TestCheckResourceAttr("btp_directory.uut", "features.#", "3"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceDirectoryWithFeatures("uut", "my-updated-directory-feat", "This is a updated directory with features"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "parent_id", regexpValidUUID),
						resource.TestCheckResourceAttr("btp_directory.uut", "name", "my-updated-directory-feat"),
						resource.TestCheckResourceAttr("btp_directory.uut", "description", "This is a updated directory with features"),
						resource.TestCheckResourceAttr("btp_directory.uut", "features.#", "3"),
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

func hclResourceDirectory(resourceName string, displayName string, description string) string {
	return fmt.Sprintf(`resource "btp_directory" "%s" {
        name        = "%s"
        description = "%s"
    }`, resourceName, displayName, description)
}

func hclResourceDirectoryWithFeatures(resourceName string, displayName string, description string) string {
	return fmt.Sprintf(`resource "btp_directory" "%s" {
        name        = "%s"
        description = "%s"
		features    = ["DEFAULT","ENTITLEMENTS","AUTHORIZATIONS"]
    }`, resourceName, displayName, description)
}
