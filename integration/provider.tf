terraform {
  required_providers {
    btp = {
      source  = "SAP/btp"
      version = "0.6.0-beta2"
    }
  }
}

# Configure the BTP Provider
provider "btp" {
  cli_server_url = "https://canary.cli.btp.int.sap"
}

