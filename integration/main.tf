###
# Setup of names in accordance to naming convention
###

locals {
  prefix_integration_test                       = "integration-test-"
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
  disclaimer_description                        = "Please don't modify. This is used for integration tests."
}

###
# Creation of subaccounts
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

resource "btp_subaccount" "sa_security-settings" {
  name      = local.integration_test_security_settings
  subdomain = local.integration_test_security_settings
  region    = var.region
}

###
# Creation of directories
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

