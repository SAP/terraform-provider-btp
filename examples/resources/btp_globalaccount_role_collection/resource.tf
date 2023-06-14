resource "btp_globalaccount_role_collection" "my_collection" {
  name        = "My own role collection"
  description = "A description of what the role collection is supposed to do."

  roles = [
    {
      name                 = "Global Account Admin"
      role_template_app_id = "cis-central!b13"
      role_template_name   = "GlobalAccount_Admin"
    }
  ]
}
