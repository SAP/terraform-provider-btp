data "btp_directory_entitlement" "directory_entitlement" {
  directory_id = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"
  service_name = "hana-cloud"
  plan_name    = "hana"
}

# When multiple plans share the same plan_name but have different plan_unique_identifier,
# specify plan_unique_identifier to select a specific variant.
data "btp_directory_entitlement" "directory_entitlement_with_uid" {
  directory_id           = "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"
  service_name           = "APPLICATION_RUNTIME"
  plan_name              = "MEMORY"
  plan_unique_identifier = "applicationruntime"
}
