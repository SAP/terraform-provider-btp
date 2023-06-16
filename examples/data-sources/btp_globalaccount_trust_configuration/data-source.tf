# default identity provider
data "btp_globalaccount_trust_configuration" "default" {
  origin = "sap.default"
}

# custom identity provider
data "btp_globalaccount_trust_configuration" "custom" {
  origin = "terraformint-platform"
}
