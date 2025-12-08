# Read destination fragment by subaccount and name
data "btp_subaccount_destination_fragment" "dest_fragment_by_subaccount" {
  subaccount_id = "6aa64c2f-38c1-12a3-b2e8-cf9fea769b7f"
  name          = "test-destination-fragment"
}

# Read destination fragment by subaccount, service instance and name
data "btp_subaccount_destination_fragment" "dest_fragment_by_subaccount" {
  subaccount_id       = "6aa64c2f-38c1-12a3-b2e8-cf9fea769b7f"
  name                = "test-destination-fragment"
  service_instance_id = "bc8a216f-1184-12a3-b4b4-17cfe2828965"
}
