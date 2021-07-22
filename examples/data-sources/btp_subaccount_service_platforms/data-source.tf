# look up all service platforms
data "btp_subaccount_service_platforms" "all" {
  subaccount_id = "98cbd1c8-49e2-42d5-8266-980e3e8728a4"
}

# look up all service platforms of type kubernetes
data "btp_subaccount_service_platforms" "k8s" {
  subaccount_id = "98cbd1c8-49e2-42d5-8266-980e3e8728a4"
  fields_filter = "type eq 'kubernetes'"
}
