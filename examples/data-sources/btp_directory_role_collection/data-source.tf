data "btp_directory_role_collection" "directory_admin" {
  directory_id = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"
  name         = "Directory Administrator"
}

# Show which groups and attribute rules from identity providers grant this role collection
data "btp_directory_role_collection" "directory_admin_with_attribute_mappings" {
  directory_id           = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"
  name                   = "Directory Administrator"
  show_attribute_mappings = true
}

# Show all users who have this role collection, including those assigned via groups
data "btp_directory_role_collection" "directory_admin_with_user_assignments" {
  directory_id         = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"
  name                 = "Directory Administrator"
  show_user_assignments = true
}
