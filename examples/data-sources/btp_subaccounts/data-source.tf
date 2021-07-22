# look up all available subaccounts of a global account
data "btp_subaccounts" "all" {}

# look up all available subaccounts of a global acount that have a specific label attached
data "btp_subaccounts" "filtered" {
  labels_filter = "my-label=my-value"
}
