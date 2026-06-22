data "btp_subaccount_entitlement" "entitlement" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_name  = "hana-cloud"
  plan_name     = "hana"
}

# When multiple plans share the same plan_name but have different plan_unique_identifier,
# specify plan_unique_identifier to select a specific variant.
data "btp_subaccount_entitlement" "entitlement_with_uid" {
  subaccount_id          = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_name           = "APPLICATION_RUNTIME"
  plan_name              = "MEMORY"
  plan_unique_identifier = "applicationruntime"
}
