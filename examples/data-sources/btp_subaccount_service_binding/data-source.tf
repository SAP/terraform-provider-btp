# look up service binding by id
data "btp_subaccount_service_binding" "by_id" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  id            = "c2d02852-1678-4c1e-b546-74d5274f1522"
}

# look up service binding by name
data "btp_subaccount_service_binding" "by_name" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "hyperspace-2022-10"
}
