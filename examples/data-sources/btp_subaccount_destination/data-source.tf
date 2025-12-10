# Read a destination by name for a subaccount
data "btp_subaccount_destination" "test_subaccount" {
  name                = "my-destination"
  subaccount_id       = "6aa64c2f-38c1-12345-a1b2-cf9fea769b7f"
  service_instance_id = "7bb75d3e-49d2-67890-b2c3-df0feb870c8e"
}

# Read a destination by name for a subaccount and service instance
data "btp_subaccount_destination" "test_service_instance" {
  name                = "my-destination"
  subaccount_id       = "6aa64c2f-38c1-12345-a1b2-cf9fea769b7f"
  service_instance_id = "7bb75d3e-49d2-67890-b2c3-df0feb870c8e"
}
