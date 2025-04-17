terraform {
  required_providers {
    btp = {
      source  = "SAP/btp"
      version = "1.11.0"
    }
  }
}

# Configure the BTP Provider
provider "btp" {
  globalaccount = "terraformintprod"
}

data "btp_directory_entitlement" "pui" {
  directory_id = "aebacff5-d7d9-484c-8ac7-a382273059c4"
  service_name = "hana-cloud"
  plan_name    = "hana"
}

resource "btp_directory_entitlement" "addpui" {
  directory_id           = data.btp_directory_entitlement.pui.directory_id
  service_name           = data.btp_directory_entitlement.pui.service_name
  plan_name              = data.btp_directory_entitlement.pui.plan_name
  plan_unique_identifier = data.btp_directory_entitlement.pui.plan_unique_identifier
}

output "pop" {
  value = data.btp_directory_entitlement.pui.plan_unique_identifier
}