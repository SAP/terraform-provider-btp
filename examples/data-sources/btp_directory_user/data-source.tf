# look up user details which belong to the default identity provider on directory level
data "btp_directory_user" "someone" {
  directory_id = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0" # directory must be security enabled!
  user_name    = "john.doe@mycompany.com"
}

# look up user details which belong to a custom identity provider on directory level
data "btp_directory_user" "someone_else" {
  directory_id = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0" # directory must be security enabled!
  user_name    = "jane.doe@mycompany.com"
  origin       = "my-custom-idp"
}
