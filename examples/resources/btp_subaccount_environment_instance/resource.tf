# creates a cloud foundry environment in a given account
resource "btp_subaccount_environment_instance" "cloudfoundry" {
  subaccount_id    = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name             = "my-cf-environment"
  environment_type = "cloudfoundry"
  service_name     = "cloudfoundry"
  plan_name        = "standard"

  # some regions offer multiple environments of a kind and you must explicitly select the target environment in which
  # the instance shall be created. 
  # available environments can be looked up using the btp_subaccount_environments datasource
  parameters = jsonencode({
    instance_name = "my-cf-org-name"
  })
}


# creates a cloud foundry environment in a given account and grant the orchestration user admin access to it
# In additon add a custom timeout for the create and delete operation
resource "btp_subaccount_environment_instance" "cloudfoundry" {
  subaccount_id    = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  name             = "my-cf-environment"
  environment_type = "cloudfoundry"
  service_name     = "cloudfoundry"
  plan_name        = "standard"
  timeouts {
    create = "1h"
    delete = "30m"
  }
  # some regions offer multiple environments of a kind and you must explicitly select the target environment in which
  # the instance shall be created. 
  # available environments can be looked up using the btp_subaccount_environments datasource
  parameters = jsonencode({
    instance_name = "my-cf-org-name"
  })
}
