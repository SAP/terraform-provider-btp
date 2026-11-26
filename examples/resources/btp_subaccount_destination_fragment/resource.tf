# create a subaccount for a subaccount
resource "btp_subaccount_destination_fragment" "dest_fragment_for_sa" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "my-fragment"
  fragment_content = {
    "key1" = "value1"
  }
}


# create a destination fragment for a subaccount and service instance
resource "btp_subaccount_destination_fragment" "dest_fragment_for_si" {
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "bc8a216f-1184-49dc-b4b4-17cfe2828965"
  name                = "my-fragment"
  fragment_content = {
    "key1" = "value1"
  }
}

# create a destination fragment for a subaccount and service instance on a CF extension landscape
resource "btp_subaccount_destination_fragment" "dest_fragment_for_si_on_cf_ext" {
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "ac8a216f-1184-49dc-b4b4-17cfe2828965"
  extension_landscape = "cf-eu12-001"
  name                = "my-fragment"
  fragment_content = {
    "key1" = "value1"
  }
}
