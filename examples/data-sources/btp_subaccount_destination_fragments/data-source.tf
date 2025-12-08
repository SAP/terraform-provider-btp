# Read destination fragments by subaccount
data "btp_subaccount_destination_fragment" "dest_fragment_by_subaccount" {
  subaccount_id = "6aa64c2f-38c1-12a3-b2e8-cf9fea769b7f"
}

# Read destination fragments by subaccount and service instance
data "btp_subaccount_destination_fragment" "dest_fragment_by_subaccount_and_service_instance" {
  subaccount_id       = "6aa64c2f-38c1-12a3-b2e8-cf9fea769b7f"
  service_instance_id = "bc8a216f-1184-12a3-b4b4-17cfe2828965"
}
