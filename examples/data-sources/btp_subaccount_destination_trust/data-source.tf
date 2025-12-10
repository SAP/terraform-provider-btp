# Read BTP Subaccount Destination Trust information for a specific subaccount and origin
data "btp_subaccount_destination_trust" "subaccount_dt_active" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  origin        = "my.origin"
}
