terraform {
  required_providers {
    btp = {
      source = "SAP/btp"
    }
  }
}

# Configure the BTP Provider
provider "btp" {
  globalaccount  = var.globalaccount
  cli_server_url = var.cli_server
}
