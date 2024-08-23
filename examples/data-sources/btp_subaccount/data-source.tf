# Read a subaccount by ID
data "btp_subaccount" "my_account_byid" {
  id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

# Read a subaccount by region and subdomain
data "btp_subaccount" "my_account_bysubdomain" {
  region    = "eu10"
  subdomain = "my-subaccount-subdomain"
}