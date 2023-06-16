# default identity provider
data "btp_subaccount_trust_configuration" "default" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  origin        = "sap.default"
}

# custom identity provider
data "btp_subaccount_trust_configuration" "custom" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  origin        = "terraformint-platform"
}
