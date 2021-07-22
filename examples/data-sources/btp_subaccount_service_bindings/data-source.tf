# look up all service bindings of a given subaccount
data "btp_subaccount_service_bindings" "all" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

# look up service bindings of a given subaccount wich have certain label assigned
data "btp_subaccount_service_bindings" "labeled" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  labels_filter = "subaccount_id eq '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f'"
}
