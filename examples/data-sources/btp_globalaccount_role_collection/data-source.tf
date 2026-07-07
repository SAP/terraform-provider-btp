data "btp_globalaccount_role_collection" "globalaccount_admin" {
  name = "Global Account Administrator"
}

# Show which groups and attribute rules from identity providers grant this role collection
data "btp_globalaccount_role_collection" "globalaccount_admin_with_attribute_mappings" {
  name                   = "Global Account Administrator"
  show_attribute_mappings = true
}

# Show all users who have this role collection, including those assigned via groups
data "btp_globalaccount_role_collection" "globalaccount_admin_with_user_assignments" {
  name                 = "Global Account Administrator"
  show_user_assignments = true
}
