resource "btp_directory" "parent" {
  name        = "my-parent-directory"
  description = "This is a parent directory."
}

resource "btp_directory" "child" {
  parent_id = btp_directory.parent.id
  name        = "my-child-directory"
  description = "This is a child directory."
}
