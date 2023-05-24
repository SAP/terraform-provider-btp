# look up all available service instances for a given subaccount
data "btp_subaccount_service_instances" "all" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

# look up all available service instances in a subaccount which are ready to be used
data "btp_subaccount_service_instances" "ready" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  fields_filter = "ready eq 'true'"
}

# look up all available service instance in a subaccount which have a certain label assigned
data "btp_subaccount_service_instances" "by_label" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  labels_filter = "labelname eq 'labelvalue'"
}
