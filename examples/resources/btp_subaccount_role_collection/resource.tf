resource "btp_subaccount_role_collection" "my_collection" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "My own role collection"
  description   = "A description of what the role collection is supposed to do."

  roles = [
    {
      name                 = "Subaccount Admin"
      role_template_app_id = "cis-local!b2"
      role_template_name   = "Subaccount_Admin"
    }
  ]
}
