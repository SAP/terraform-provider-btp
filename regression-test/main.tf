###
# setup of names in accordance to naming convention
###

locals {
  prefix_integration_test                       = "regression-test-"
  prefix_integration_test_dir                   = "${local.prefix_integration_test}dir-"
  prefix_integration_test_account               = "${local.prefix_integration_test}acc-"
  integration_test_account_static               = "${local.prefix_integration_test_account}static"
  integration_test_account_entitlements_stacked = "${local.prefix_integration_test_account}entitlements-stacked"
  integration_test_services_static              = "${local.prefix_integration_test}services-static"
  integration_test_security_settings            = "${local.prefix_integration_test}security-settings"
  integration_test_dir_static                   = "${local.prefix_integration_test_dir}static"
  integration_test_dir_se_static                = "${local.prefix_integration_test_dir}se-static"
  integration_test_dir_entitlements             = "${local.prefix_integration_test_dir}entitlements"
  integration_test_dir_entitlements_stacked     = "${local.prefix_integration_test_dir}entitlements-stacked"
  disclaimer_description                        = "Don't modify. This is used for regression tests."
}

###
# subaccounts
###

resource "btp_subaccount" "sa_acc_static" {
  name        = local.integration_test_account_static
  description = local.disclaimer_description
  subdomain   = local.integration_test_account_static
  region      = var.region
  labels = {
    label1 = [
      "label text 1"
    ]
    label2 = []
  }
}

resource "btp_subaccount" "sa_acc_entitlements_stacked" {
  parent_id = btp_directory.dir_entitlements_stacked.id
  name      = local.integration_test_account_entitlements_stacked
  subdomain = local.integration_test_account_entitlements_stacked
  region    = var.region
}

resource "btp_subaccount" "sa_services_static" {
  name        = local.integration_test_services_static
  subdomain   = local.integration_test_services_static
  region      = var.region
  description = "Subaccount to test:\n- Service Instances\n- Service Bindings\n- App Subscriptions"
}

resource "btp_subaccount" "sa_security_settings" {
  name      = local.integration_test_security_settings
  subdomain = local.integration_test_security_settings
  region    = var.region
}

###
# directories
###

resource "btp_directory" "dir_entitlements" {
  name     = local.integration_test_dir_entitlements
  features = ["DEFAULT", "ENTITLEMENTS", "AUTHORIZATIONS"]
}

resource "btp_directory" "dir_entitlements_stacked" {
  parent_id = btp_directory.dir_entitlements.id
  name      = local.integration_test_dir_entitlements_stacked
}

resource "btp_directory" "dir_static" {
  name        = local.integration_test_dir_static
  description = local.disclaimer_description
}

resource "btp_directory" "dir_se_static" {
  name        = local.integration_test_dir_se_static
  description = local.disclaimer_description
  features    = ["DEFAULT", "ENTITLEMENTS", "AUTHORIZATIONS"]
  labels = {
    my-label-1 = [
      "Label text 1"
    ]
    my-label-2 = []
  }
}

###
# directory role collection assignments
###

resource "btp_directory_role_collection_assignment" "drca_dir_se_static_jenny_directory_viewer" {
  directory_id         = btp_directory.dir_se_static.id
  role_collection_name = "Directory Viewer"
  user_name            = "jenny.doe@test.com"
}

###
# global account role collection assignments
###

resource "btp_globalaccount_role_collection_assignment" "grca_jenny_ga_viewer" {
  role_collection_name = "Global Account Viewer"
  user_name            = "jenny.doe@test.com"
}

###
# subaccount role collection assignments
###
resource "btp_subaccount_role_collection_assignment" "srca_sa_acc_static_subaccount_administrators" {
  subaccount_id        = btp_subaccount.sa_acc_static.id
  for_each             = toset(["BTP Terraform Administrator", "BTP Terraform Developer"])
  role_collection_name = "Subaccount Administrator"
  group_name           = each.value
}

resource "btp_subaccount_role_collection_assignment" "srca_sa_acc_entitlements_stacked_subaccount_administrators" {
  subaccount_id        = btp_subaccount.sa_acc_entitlements_stacked.id
  for_each             = toset(["BTP Terraform Administrator", "BTP Terraform Developer"])
  role_collection_name = "Subaccount Administrator"
  group_name           = each.value
}


###
# directory entitlements
###
resource "btp_directory_entitlement" "de_dir_entitlements_hana_cloud" {
  directory_id = btp_directory.dir_entitlements.id
  service_name = "hana-cloud"
  plan_name    = "hana"
}

###
# global account resource provider
###

resource "btp_globalaccount_resource_provider" "grp_aws" {
  technical_name = "tf_regtest_resource_provider"
  display_name   = "Test AWS Resource Provider"
  description    = "Description of the resource provider"
  provider_type  = "AWS"
  configuration = jsonencode({
    access_key_id     = "AWSACCESSKEY"
    secret_access_key = "AWSSECRETKEY"
    vpc_id            = "vpc-test"
    region            = "eu-central-1"
  })
}

###
# global account trust configuration
###

#resource "btp_globalaccount_trust_configuration" "gtc_idp_terraformint" {
#  identity_provider = "terraformint.accounts400.ondemand.com"
#  name              = "terraformint-platform"
#  description       = "Custom Platform Identity Provider"
#  origin            = "terraformint-platform"
#}

###
# subaccount subscriptions
###

resource "btp_subaccount_subscription" "sas_sa_services_static_content_agent_ui" {
  subaccount_id = btp_subaccount.sa_services_static.id
  app_name      = "content-agent-ui"
  plan_name     = "free"
}

###
# subaccount entitlements
###

resource "btp_subaccount_entitlement" "se_sa_services_static_auditlog_viewer" {
  subaccount_id = btp_subaccount.sa_services_static.id
  service_name  = "auditlog-viewer"
  plan_name     = "free"
}

resource "btp_subaccount_entitlement" "se_sa_services_static_alert_notification" {
  subaccount_id = btp_subaccount.sa_services_static.id
  service_name  = "alert-notification"
  plan_name     = "free"
}

resource "btp_subaccount_entitlement" "se_sa_services_static_malware_scanner" {
  subaccount_id = btp_subaccount.sa_services_static.id
  service_name  = "malware-scanner"
  plan_name     = "clamav"
}

###
# subaccount service instances
###

data "btp_subaccount_service_plan" "ssp_sa_services_static_alert_notification_free" {
  subaccount_id = btp_subaccount.sa_services_static.id
  name          = "free"
  offering_name = "alert-notification"
  depends_on = [
    btp_subaccount_entitlement.se_sa_services_static_alert_notification
  ]
}

resource "btp_subaccount_service_instance" "ssi_sa_services_static_alert_notification_free" {
  subaccount_id  = btp_subaccount.sa_services_static.id
  serviceplan_id = data.btp_subaccount_service_plan.ssp_sa_services_static_alert_notification_free.id
  name           = "tf-regtest-alertnotification-instance"
}

data "btp_subaccount_service_plan" "ssp_sa_services_static_malware_scanner_default" {
  subaccount_id = btp_subaccount.sa_services_static.id
  name          = "clamav"
  offering_name = "malware-scanner"
  depends_on = [
    btp_subaccount_entitlement.se_sa_services_static_malware_scanner
  ]
}

resource "btp_subaccount_service_instance" "ssi_sa_services_static_malware_scanner_default" {
  subaccount_id  = btp_subaccount.sa_services_static.id
  serviceplan_id = data.btp_subaccount_service_plan.ssp_sa_services_static_malware_scanner_default.id
  name           = "tf-regtest-malware-scanner-sample"
  labels = {
    org = [
      "testvalue"
    ]
  }
}

###
# subaccount service bindings
###

resource "btp_subaccount_service_binding" "binding_sa_services_static_alert_notification_free_sb_test" {
  subaccount_id       = btp_subaccount.sa_services_static.id
  service_instance_id = btp_subaccount_service_instance.ssi_sa_services_static_alert_notification_free.id
  name                = "regtest-service-binding"
}

resource "btp_subaccount_service_binding" "binding_sa_services_static_alert_notification_free_sb_test_two" {
  subaccount_id       = btp_subaccount.sa_services_static.id
  service_instance_id = btp_subaccount_service_instance.ssi_sa_services_static_alert_notification_free.id
  name                = "regtest-service-binding-two"
}
