# look up offering by ID and subaccount ID
data "btp_subaccount_service_offering" "by_id" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  id            = "4e953cf1-7eda-4ebb-a58c-02c6ebfe45fb"
}

# look up offering by offering name and subaccount ID
data "btp_subaccount_service_offering" "by_name" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name          = "auditlog-management"
}
