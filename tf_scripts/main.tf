# resource "btp_subaccount_service_broker" "sb_sa_services_static" {
#   subaccount_id       = "59cd458e-e66e-4b60-b6d8-8f219379f9a5"
#   name                = "my-broker"
#   url                 = "https://my-broker-sweet-bongo-pm.cfapps.eu12.hana.ondemand.com"
#   username            = "platform"
#   password            = "a-secure-password"
# }


# data "btp_subaccounts" "all" {}
# data "btp_subaccount" "test" {
# 	id = [for sa in data.btp_subaccounts.all.values : sa.id if sa.name == "%s"][0]
# }

# data "btp_directories" "all" {}
# data "btp_directory_role_collections" "test" {
# 	directory_id = [for dir in data.btp_directories.all.values : dir.id if dir.name == "integration-test-2025-dir-se-static"][0]
# }