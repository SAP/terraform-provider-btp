# look up all users which belong to the default identity provider on subaccount level
data "btp_subaccount_users" "defaultidp" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

# look up all users which belong to a custom identity provider on subaccount level
data "btp_subaccount_users" "mycustomidp" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  origin        = "my-custom-idp"
}
