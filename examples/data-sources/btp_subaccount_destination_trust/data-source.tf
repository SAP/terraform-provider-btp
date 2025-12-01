data "btp_subaccount_destination_trust" "subaccount_dt_active" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  active        = true // to fetch active destination trust
}
