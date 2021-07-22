terraform {
  required_providers {
    btp = {
      source  = "sap/btp"
      version = "~> 0.3"
    }
  }
}

# Configure the BTP Provider
provider "btp" {
  globalaccount = "795b53bb-a3f0-4769-adf0-26173282a975"
}
