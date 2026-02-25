terraform {
  required_providers {
    btp = {
      source  = "SAP/btp"
      version = "1.24.0"
    }
  }
}

# Configure the BTP Provider
provider "btp" {
  globalaccount  = "3048e80f-931a-4b05-b7e1-5abd3146639c"
//  cli_server_url = "https://staging.cli.btp.int.sap"
}

/*
import {
  to = btp_subaccount.my_project
  id = "249ca2ea-12ac-4e69-8cb6-b9652524cf2e"
}
*/

resource "btp_subaccount" "my_project" {
  name      = "test-delayed-deletion-2"
  subdomain = "test-delayed-deletion-112233"
  region    = "us10"
}

/*
resource "btp_subaccount_environment_instance" "cloudfoundry" {
  subaccount_id    = btp_subaccount.my_project.id
  name             = "my-cf-environment-cl"
  environment_type = "cloudfoundry"
  service_name     = "cloudfoundry"
  plan_name        = "standard"
  landscape_label  = "cf-us10"
  parameters = jsonencode({
    instance_name = "my-cf-org-name"
  })
}
*/
