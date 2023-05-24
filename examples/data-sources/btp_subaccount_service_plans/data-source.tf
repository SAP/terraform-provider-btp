# look up all service plans of a given subaccount
data "btp_subaccount_service_plans" "all" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

# look up all services available on sapbtp environment in a given subaccount
data "btp_subaccount_service_plans" "sapbtp" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  environment   = "sapbtp"
}

# look up all services available on kubernetes environment in a given subaccount
data "btp_subaccount_service_plans" "k8s" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  environment   = "kubernetes"
}

# look up all services wich have certain label assigned in a given subaccount
data "btp_subaccount_service_plans" "labeled" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  labels_filter = "commercial_name eq 'application'"
}
