# Read all destination for a given subaccount
data "btp_subaccount_destinations" "dest_by_subaccount" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

# Read all destination for a given subaccount with names only
data "btp_subaccount_destinations" "dest_by_subaccount_names_only" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  names_only    = true
}

# Read all destination for a given subaccount and service instance
data "btp_subaccount_destinations" "dest_by_subaccount_and_service_instance" {
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "5aa54c2f-47a1-49a9-b2e8-cf9fea769b7f"
}
