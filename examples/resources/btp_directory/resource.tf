# Create a parent directory without features enabled
resource "btp_directory" "parent" {
  name        = "my-parent-directory"
  description = "This is a parent directory."
}

# Create a child directory underneath a parent directory without features enabled
resource "btp_directory" "child" {
  parent_id   = btp_directory.parent.id
  name        = "my-child-directory"
  description = "This is a child directory."
}

# Create a directory with ENTITLEMENT and AUTHORIZATIONS features enabled
resource "btp_directory" "dir_with_features" {
  name        = "my-feat-directory"
  description = "This is a directory with features."
  features    = ["DEFAULT","ENTITLEMENTS","AUTHORIZATIONS"]
}
