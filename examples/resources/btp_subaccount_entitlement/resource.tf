# entitle service plan without quota in a subaccount
resource "btp_subaccount_entitlement" "alert_notification_service" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_name  = "alert-notification"
  plan_name     = "free"
}

# entitle service plan with quota in a subaccount
resource "btp_subaccount_entitlement" "uas_reporting" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_name  = "uas"
  plan_name     = "reporting-directory"
  amount        = 1
}
