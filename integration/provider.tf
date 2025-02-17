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
  cli_server_url = "https://canary.cli.btp.int.sap"
}
