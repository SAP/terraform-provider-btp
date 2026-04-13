###
# setup of names in accordance to naming convention
###

locals {
  prefix_regression_test                                = "tf-regression-test-"
  regression_test_services_static                       = "${local.prefix_regression_test}services-static"
  regression_test_services_static_extended              = "${local.regression_test_services_static}-${var.subaccount_subdomain_extension}"
}

###
# subaccounts
###

resource "btp_subaccount" "sa_services_static" {
  name        = local.regression_test_services_static
  subdomain   = local.regression_test_services_static_extended
  region      = var.region
  description = "Subaccount to test:\n- Environment Instance\n- Service Instances\n- Service Bindings\n- App Subscriptions"
}

###
# subaccount role collection assignments
###

resource "btp_subaccount_role_collection_assignment" "srca_sa_services_static_destination_creator" {
  subaccount_id        = btp_subaccount.sa_services_static.id
  role_collection_name = "Destination Administrator"
  user_name            = "jenny.doe@test.com"
}

###
# subaccount entitlements
###

resource "btp_subaccount_entitlement" "se_sa_services_static_bas" {
  subaccount_id = btp_subaccount.sa_services_static.id
  service_name  = "sapappstudio"
  plan_name     = "standard-edition"
}

resource "btp_subaccount_entitlement" "se_sa_services_static_destination" {
  subaccount_id = btp_subaccount.sa_services_static.id
  service_name  = "destination"
  plan_name     = "lite"
}

###
# subaccount subscriptions
###

resource "btp_subaccount_subscription" "sas_sa_services_static_bas" {
  subaccount_id = btp_subaccount.sa_services_static.id
  app_name      = "sapappstudio"
  plan_name     = "standard-edition"
  depends_on    = [btp_subaccount_entitlement.se_sa_services_static_bas]
}

###
# service instances
###
data "btp_subaccount_service_plan" "ssp_sa_services_static_destination" {
  subaccount_id = btp_subaccount.sa_services_static.id
  name          = "lite"
  offering_name = "destination"
  depends_on = [
    btp_subaccount_entitlement.se_sa_services_static_destination
  ]
}

resource "btp_subaccount_service_instance" "ssi_sa_services_destination" {
  subaccount_id  = btp_subaccount.sa_services_static.id
  serviceplan_id = data.btp_subaccount_service_plan.ssp_sa_services_static_destination.id
  name           = "tf-testacc-destination-instance"
}

###
# subaccount destinations
###

# subaccount destination with service instance
resource "btp_subaccount_destination_generic" "http_dest" {
  subaccount_id       = btp_subaccount.sa_services_static.id
  service_instance_id = btp_subaccount_service_instance.ssi_sa_services_destination.id
  destination_configuration = jsonencode({
    "Name"           = "destination-with-service-instance"
    "Type"           = "HTTP"
    "ProxyType"      = "Internet"
    "URL"            = "https://myservice.example.com"
    "Authentication" = "NoAuthentication"
    "Description"    = "trial destination of basic usecase with service instance"
  })
  depends_on = [
    btp_subaccount_role_collection_assignment.srca_sa_services_static_destination_creator
  ]
}

# subaccount destination without service instance and with additional variables for authentication
resource "btp_subaccount_destination_generic" "http_dest_with_destination_configuration_authentication" {
  subaccount_id = btp_subaccount.sa_services_static.id
  destination_configuration = jsonencode({
    "Name"            = "destination-with-authentication"
    "Type"            = "HTTP"
    "clientId"        = "abc"
    "tokenServiceURL" = "https://myservice.example.com"
    "ProxyType"       = "Internet"
    "URL"             = "https://myservice.example.com"
    "Authentication"  = "OAuth2ClientCredentials"
    "Description"     = "trial destination of basic usecase with additional variables "
  })
  depends_on = [
    btp_subaccount_role_collection_assignment.srca_sa_services_static_destination_creator
  ]
}