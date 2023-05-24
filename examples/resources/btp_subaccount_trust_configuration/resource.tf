# create a new simple trust configuration for a subaccount
resource "btp_subaccount_trust_configuration" "simple" {
  subaccount_id     = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  identity_provider = "cloudorchestration.accounts.ondemand.com"
}

# create a new fully customized trust configuration for a subaccount
resource "btp_subaccount_trust_configuration" "fully_customized" {
  subaccount_id     = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  identity_provider = "cloudorchestration.accounts.ondemand.com"
  name              = "my-name"
  description       = "my-description"
  origin            = "my-own-origin-platform"
}
