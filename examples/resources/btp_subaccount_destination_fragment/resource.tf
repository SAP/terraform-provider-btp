resource "btp_subaccount_destination_fragment" "sdf" {
  subaccount_id        = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name                 = "my-fragment"
  fragment_content = {
    "key1" = "value1"
  }
}
