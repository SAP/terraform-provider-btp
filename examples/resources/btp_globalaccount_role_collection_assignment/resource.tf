# assign a user to a role collection on global account level
resource "btp_globalaccount_role_collection_assignment" "jd" {
  role_collection_name = "Global Account Viewer"
  user_name            = "john.doe@mycompany.com"
}

# assign a group to a role collection on global account level
resource "btp_globalaccount_role_collection_assignment" "globalaccount_viewer_group" {
  role_collection_name = "Global Account Viewer"
  group_name           = "globalaccount-viewer-group"
}
