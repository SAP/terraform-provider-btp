data "btp_subaccount_role_collection" "subaccount_admin" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "Subaccount Administrator"
}

# Show which groups and attribute rules from identity providers grant this role collection
data "btp_subaccount_role_collection" "subaccount_admin_with_attribute_mappings" {
  subaccount_id           = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name                    = "Subaccount Administrator"
  show_attribute_mappings = true
}

# Show all users who have this role collection, including those assigned via groups
data "btp_subaccount_role_collection" "subaccount_admin_with_user_assignments" {
  subaccount_id         = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name                  = "Subaccount Administrator"
  show_user_assignments = true
}
