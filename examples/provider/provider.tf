terraform {
  required_providers {
    btp = {
      source  = "SAP/btp"
      version = "0.2.0-beta1"
    }
  }
}

# Configure the BTP Provider
provider "btp" {
  globalaccount = "my-global-account-subdomain"
}
