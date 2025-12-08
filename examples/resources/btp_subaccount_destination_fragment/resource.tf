# create a subaccount for a subaccount
resource "btp_subaccount_destination_fragment" "sdf" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "my-fragment"
  fragment_content = {
    "key1" = "value1"
  }
}


# create a subaccount for a subaccount and service instance
resource "btp_subaccount_destination_fragment" "sdf" {
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "bc8a216f-1184-49dc-b4b4-17cfe2828965"
  name                = "my-fragment"
  fragment_content = {
    "key1" = "value1"
  }
}
