# assign a single ser to a role collection on subaccount level
resource "btp_subaccount_role_collection_assignment" "jd" {
  subaccount_id        = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  role_collection_name = "Subaccount Viewer"
  user_name            = "john.doe@mycompany.com"
}

# assign a group to a role collection on subaccount level
resource "btp_subaccount_role_collection_assignment" "subaccount_viewer_group" {
  subaccount_id        = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  role_collection_name = "Subaccount Viewer"
  group_name           = "subaccount-viewer-group"
}
