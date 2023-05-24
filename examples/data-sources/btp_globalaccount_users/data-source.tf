# look up all users which belong to the default identity provider on global account level
data "btp_globalaccount_users" "defaultidp" {}

# look up all users which belong to a custom identity provider on global account level
data "btp_globalaccount_users" "mycustomidp" {
  origin = "my-custom-idp"
}
