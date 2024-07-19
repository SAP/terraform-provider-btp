resource "btp_subaccount_service_broker" "my_broker" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "my-broker"
  description   = "Service broker for provisioning example services."

  url      = "https://my.broker.com"
  username = "platform_user"
  password = "platform_password"
}
