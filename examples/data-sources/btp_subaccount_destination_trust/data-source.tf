# Read BTP Subaccount Destination Trust information for a specific subaccount
data "btp_subaccount_destination_trust" "subaccount_dt_active" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

# Read BTP Subaccount Destination Trust information for a specific subaccount and trust_type
data "btp_subaccount_destination_trust" "subaccount_dt_active_trust_type" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  trust_type    = "ACTIVE"
}
