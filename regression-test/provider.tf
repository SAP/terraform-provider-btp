terraform {
  required_providers {
    btp = {
      source  = "SAP/btp"
      version = "1.0.0-rc1"
    }
  }
}

# Configure the BTP Provider
provider "btp" {
  globalaccount  = var.globalaccount
  cli_server_url = var.cli_server
}