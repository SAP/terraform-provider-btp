terraform {
  required_providers {
    btp = {
      source  = "SAP/btp"
      version = "1.10.0"
    }
  }
}

# Configure the BTP Provider
provider "btp" {
  globalaccount = "my-global-account-subdomain"
}
