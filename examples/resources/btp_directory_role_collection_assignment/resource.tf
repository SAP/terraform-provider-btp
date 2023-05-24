# assign a single user to a role collection on directory level
resource "btp_directory_role_collection_assignment" "jd" {
  directory_id         = "ddfc2206-5f11-48ed-a1ec-29010af70050"
  role_collection_name = "Directory Viewer"
  user_name            = "john.doe@mycompany.com"
}

# assign a group to a role collection on directory level
resource "btp_directory_role_collection_assignment" "directory_viewer_group" {
  directory_id         = "ddfc2206-5f11-48ed-a1ec-29010af70050"
  role_collection_name = "Directory Viewer"
  group_name           = "directory-viewer-group"
}
