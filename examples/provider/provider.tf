terraform {
  required_providers {
    btp = {
      source  = "SAP/btp"
      version = "1.28.1"
    }
  }
}

# Configure the BTP Provider
provider "btp" {
  globalaccount = "my-global-account-subdomain"
}
