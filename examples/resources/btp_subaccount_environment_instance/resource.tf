# creates a cloud foundry environment in a given account
resource "btp_subaccount_environment_instance" "cloudfoundry" {
  subaccount_id    = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name             = "my-cf-environment"
  environment_type = "cloudfoundry"
  service_name     = "cloudfoundry"
  plan_name        = "standard"
  # ATTENTION: some regions offer multiple environments of a kind and you must explicitly select the target environment in which
  # the instance shall be created using the parameter landscape label. 
  # available environments can be looked up using the btp_subaccount_environments datasource
  parameters = jsonencode({
    instance_name = "my-cf-org-name"
  })
}


# creates a cloud foundry environment in a given account
# in additon add a custom timeout for the create and delete operation
resource "btp_subaccount_environment_instance" "cloudfoundry" {
  subaccount_id    = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name             = "my-cf-environment"
  environment_type = "cloudfoundry"
  service_name     = "cloudfoundry"
  plan_name        = "standard"
  # ATTENTION: some regions offer multiple environments of a kind and you must explicitly select the target environment in which
  # the instance shall be created using the parameter landscape label. 
  # available environments can be looked up using the btp_subaccount_environments datasource
  parameters = jsonencode({
    instance_name = "my-cf-org-name"
  })
  timeouts = {
    create = "1h"
    update = "35m"
    delete = "30m"
  }
}


# creates a cloud foundry environment in a given account
# and the dedicted target landscape cf-us10
resource "btp_subaccount_environment_instance" "cloudfoundry" {
  subaccount_id    = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name             = "my-cf-environment"
  environment_type = "cloudfoundry"
  service_name     = "cloudfoundry"
  landscape_label  = "cf-us10"
  plan_name        = "standard"
  parameters = jsonencode({
    instance_name = "my-cf-org-name"
  })
}


# creates a kyma environment in a given account
# NOTE: for the available parameter values, check https://help.sap.com/docs/btp/sap-business-technology-platform/available-plans-in-kyma-environment
resource "btp_subaccount_environment_instance" "kyma" {
  subaccount_id    = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name             = "my-kyma-environment"
  environment_type = "kyma"
  service_name     = "kymaruntime"
  plan_name        = "aws"

  parameters = jsonencode({
    name          = "my-kyma-environment"
    region        = "us-east-1"
    machineType   = "mx5.xlarge" #smallest option
    autoScalerMin = 3
    autoScalerMax = 20
  })
  timeouts = {
    create = "1h"
    update = "35m"
    delete = "1h"
  }
}
