# look up a sercvice plan by ID and subaccount ID
data "btp_subaccount_service_plan" "by_id" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  id            = "b50d1b0b-2059-4f21-a014-2ea87752eb48"
}

# look up a sercvice plan by plan name and offering name and subaccount ID
data "btp_subaccount_service_plan" "by_name" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "free"
  offering_name = "alert-notification"
}
