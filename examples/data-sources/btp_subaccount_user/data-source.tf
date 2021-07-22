# look up user details which belong to the default identity provider on subaccount level
data "btp_subaccount_user" "someone" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  user_name     = "john.doe@mycompany.com"
}

# look up user details which belongs to a custom identity provider on subaccount level
data "btp_subaccount_user" "someone_else" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  user_name     = "jane.doe@mycompany.com"
  origin        = "my-custom-idp"
}
