# create an instance of the alert-notification service (no configuration necessary)
resource "btp_subaccount_service_instance" "alert_notification_free" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  # The service plan ID can be looked up via the data source btp_subaccount_service_plan
  serviceplan_id = "b50d1b0b-2059-4f21-a014-2ea87752eb48" # alert-notification - free
  name           = "my-alert-notification-instance-new"
}

# create an xsuaa service instance with additional configurations
resource "btp_subaccount_service_instance" "xsuaa_application" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  # The service plan ID can be looked up via the data source btp_subaccount_service_plan
  serviceplan_id = "bd5e893f-81dd-4d10-8343-e02975e8ecd8" # xsuaa - application
  name           = "my-application"
  parameters = jsonencode({
    xsappname   = "my-application"
    tenant-mode = "dedicated"
  })
}

# create an instance of the alert-notification service (no configuration necessary)
# in additon add a custom timeout for the create and update operation
resource "btp_subaccount_service_instance" "alert_notification_free" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  # The service plan ID can be looked up via the data source btp_subaccount_service_plan
  serviceplan_id = "b50d1b0b-2059-4f21-a014-2ea87752eb48" # alert-notification - free
  name           = "my-alert-notification-instance-new"
  timeouts = {
    create = "25m"
    update = "15m"
    delete = "15m"
  }
}

# create an instance of the xsuaa service and also share the instance
resource "btp_subaccount_service_instance" "xsuaa_application" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  # The service plan ID can be looked up via the data source btp_subaccount_service_plan
  serviceplan_id = "bd5e893f-81dd-4d10-8343-e02975e8ecd8" # xsuaa - application
  name           = "my-application"
  shared         = true
}