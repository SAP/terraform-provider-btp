# look up all users which belong to the default identity provider on directory level
data "btp_directory_users" "defaultidp" {
  directory_id = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"
}

# look up all users which belong to a custom identity provider on directory level
data "btp_directory_users" "mycustomidp" {
  directory_id = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"
  origin       = "my-custom-idp"
}
