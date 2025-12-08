# register a AZURE project as resource provider
resource "btp_globalaccount_resource_provider" "azure" {
  display_name   = "my_azure_provider"
  provider_type  = "AZURE"
  technical_name = "AZURE"
  description    = "My Azure Resource Provider"
  configuration = jsonencode({
    region              = "westeurope"
    client_id           = "AZURECLIENTID"
    client_secret       = "AZURECLIENTSECRET"
    tenant_id           = "42x7676x-f455-423x-82x6-xx2d99791xx7"
    subscription_id     = "x1x9567x-8560-44xx-x4fx-741xx0x08x58"
    resource_group_name = "rg-landscape-azure-example"
  })
}

# register an AWS account as resource provider
resource "btp_globalaccount_resource_provider" "aws" {
  display_name   = "my_aws_provider"
  provider_type  = "AWS"
  technical_name = "AWS"
  description    = "My AWS Resource Provider"
  configuration = jsonencode({
    access_key_id     = "AWSACCESSKEY"
    secret_access_key = "AWSSECRETKEY"
    vpc_id            = "vpc-test"
    region            = "eu-central-1"
  })
}
