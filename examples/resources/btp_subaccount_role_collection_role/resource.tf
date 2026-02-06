resource "btp_subaccount_role_collection_role" "my_collection_role" {
  subaccount_id        = "<subaccount_id>"
  name                 = "<collection_name>"
  role_name            = "<role_name>"
  role_template_name   = "<role_template_name>"
  role_template_app_id = "<role_template_app_id>"
}
