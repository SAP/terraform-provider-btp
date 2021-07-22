# look up user details which belongs to the default identity provider on global account level
data "btp_globalaccount_user" "someone" {
  user_name = "john.doe@mycompany.com"
}

# look up user details which belongs to a custom identity provider on global account level
data "btp_globalaccount_user" "someone_else" {
  user_name = "jane.doe@mycompany.com"
  origin    = "my-custom-idp"
}
