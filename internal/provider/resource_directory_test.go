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
		rec, user := setupVCR(t, "fixtures/resource_directory.with_features")
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

	t.Run("happy path full config with update", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory.full_config")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryAll("uut", "my-new-directory-full-config", "This is a new directory"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_directory.uut", "name", "my-new-directory-full-config"),
						resource.TestCheckResourceAttr("btp_directory.uut", "description", "This is a new directory"),
						resource.TestCheckResourceAttr("btp_directory.uut", "labels.foo.0", "bar"),
						resource.TestCheckResourceAttr("btp_directory.uut", "features.#", "3"),
					),
				},
				{
					// Update name wo change of usage but omit optional parameters
					Config: hclProviderFor(user) + hclResourceDirectory("uut", "my-updated-directory-full-config", "This is a updated directory"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_directory.uut", "name", "my-updated-directory-full-config"),
						resource.TestCheckResourceAttr("btp_directory.uut", "description", "This is a updated directory"),
						resource.TestCheckNoResourceAttr("btp_directory.uut", "labels"),
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

	t.Run("happy path - change directory features", func(t *testing.T) {
		rec, user := setupVCR(t, "fixtures/resource_directory.change_with_features")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectory("uut", "my-new-directory-features", "This is a new directory"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_directory.uut", "name", "my-new-directory-features"),
						resource.TestCheckResourceAttr("btp_directory.uut", "description", "This is a new directory"),
					),
				},
				{
					Config: hclProviderFor(user) + hclResourceDirectoryWithFeatures("uut", "my-updated-directory-features", "This is an updated directory"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_directory.uut", "name", "my-updated-directory-features"),
						resource.TestCheckResourceAttr("btp_directory.uut", "description", "This is an updated directory"),
						resource.TestCheckResourceAttr("btp_directory.uut", "features.#", "3"),
					),
				},
			},
		})
	})

	t.Run("happy path - directory hierarchy", func(t *testing.T) {
		// When recroding this test, make sure that your are not Global Account Admin, but Directory Admin of the parent directory
		rec, user := setupVCR(t, "fixtures/resource_directory.with_hierarchy")
		defer stopQuietly(rec)

		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []resource.TestStep{
				{
					Config: hclProviderFor(user) + hclResourceDirectoryWithParent("uut", "my-new-directory", "This is a new directory", "2613212d-a51e-4e7e-858c-7f96c15d67e7"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestMatchResourceAttr("btp_directory.uut", "id", regexpValidUUID),
						resource.TestMatchResourceAttr("btp_directory.uut", "created_date", regexpValidRFC3999Format),
						resource.TestMatchResourceAttr("btp_directory.uut", "last_modified", regexpValidRFC3999Format),
						resource.TestCheckResourceAttr("btp_directory.uut", "name", "my-new-directory"),
						resource.TestCheckResourceAttr("btp_directory.uut", "description", "This is a new directory"),
						resource.TestCheckResourceAttr("btp_directory.uut", "parent_id", "2613212d-a51e-4e7e-858c-7f96c15d67e7"),
					),
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

func hclResourceDirectoryWithParent(resourceName string, displayName string, description string, parentId string) string {
	return fmt.Sprintf(`resource "btp_directory" "%s" {
        name        = "%s"
        description = "%s"
				parent_id   = "%s"
    }`, resourceName, displayName, description, parentId)
}

func hclResourceDirectoryWithFeatures(resourceName string, displayName string, description string) string {
	return fmt.Sprintf(`resource "btp_directory" "%s" {
        name        = "%s"
        description = "%s"
		features    = ["DEFAULT","ENTITLEMENTS","AUTHORIZATIONS"]
    }`, resourceName, displayName, description)
}

func hclResourceDirectoryAll(resourceName string, displayName string, description string) string {
	return fmt.Sprintf(`resource "btp_directory" "%s" {
        name        = "%s"
        description = "%s"
		features    = ["DEFAULT","ENTITLEMENTS","AUTHORIZATIONS"]
		labels = {"foo" = ["bar"]}
    }`, resourceName, displayName, description)
}
