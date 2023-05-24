# create a new simple trust configuration for a global account
resource "btp_globalaccount_trust_configuration" "simple" {
  identity_provider = "cloudorchestration.accounts.ondemand.com"
}

# create a new fully customized trust configuration for a global account
resource "btp_globalaccount_trust_configuration" "fully_customized" {
  identity_provider = "cloudorchestration.accounts.ondemand.com"
  name              = "my-name"
  description       = "my-description"
  origin            = "my-own-origin-platform"
}
