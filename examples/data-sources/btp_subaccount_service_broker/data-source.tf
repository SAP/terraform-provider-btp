# look up service broker by id and subaccount id
data "btp_subaccount_service_broker" "by_id" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  id            = "9ff44f1b-b2a8-43ae-9072-32bd1dce60e4"
}

# look up service broker by name and subaccount id
data "btp_subaccount_service_broker" "by_name" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "my-broker"
}
