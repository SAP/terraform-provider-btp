action "btp_update_global_account" "update" {
  config {
    display_name                     = "My Global Account"
    description                      = "Global account for managing SAP BTP resources across all subaccounts"
    enable_subaccount_force_deletion = false
  }
}