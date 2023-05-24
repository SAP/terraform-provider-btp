# look up a service instance by its ID and subaccount ID
data "btp_subaccount_service_instance" "by_id" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  id            = "bc8a216f-1184-49dc-b4b4-17cfe2828965"
}

# look up a service instance by its name and subaccount ID
data "btp_subaccount_service_instance" "by_name" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "my-xsuaa-application"
}
