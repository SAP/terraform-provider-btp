# look up all service offerings in a given subaccount 
data "btp_subaccount_service_offerings" "all" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

# look up services offerings available on sapbtp environment in a given subaccount
data "btp_subaccount_service_offerings" "sapbtp" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  environment   = "sapbtp"
}

# look up services offerings available on kubernetes environment in a given subaccount
data "btp_subaccount_service_offerings" "k8s" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  environment   = "kubernetes"
}

# look up services offerings wich have a specifc label assigned in a given subaccount
data "btp_subaccount_service_offerings" "labeled" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  labels_filter = "a_label eq 'label-value'"
}
