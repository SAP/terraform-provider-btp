# create a subscription to workzone
resource "btp_subaccount_subscription" "workzone" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  app_name      = "SAPLaunchpad"
  plan_name     = "free"
}


# create a subscription to workzone with a timeout
resource "btp_subaccount_subscription" "workzone" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  app_name      = "SAPLaunchpad"
  plan_name     = "free"
  timeouts = {
    create = "25m"
    delete = "15m"
  }
}
