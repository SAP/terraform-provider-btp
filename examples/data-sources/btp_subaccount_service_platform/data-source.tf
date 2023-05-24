# look up service platform by ID and subaccount ID
data "btp_subaccount_service_platform" "by_id" {
  subaccount_id = "98cbd1c8-49e2-42d5-8266-980e3e8728a4"
  id            = "76765dca-6683-473a-8f42-809e33a2ea68"
}

# look up service platform by name and subaccount ID
data "btp_subaccount_service_platform" "by_name" {
  subaccount_id = "98cbd1c8-49e2-42d5-8266-980e3e8728a4"
  name          = "default-76765dca-6683-473a-8f42-809e33a2ea68"
}
