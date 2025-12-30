#subaccount destination with service instance and labels (additional configuration).
resource "btp_subaccount_destination" "destination" {
  name                = "destination"
  type                = "HTTP"
  proxy_type          = "Internet"
  url                 = "https://myservice.example.com"
  authentication      = "NoAuthentication"
  description         = "resource"
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  additional_configuration = jsonencode({
    Abc = "good"
  })
}

#subaccount destination creation without service instance and labels (additional configuration). 
resource "btp_subaccount_destination" "destination-without-service-instance" {
  name           = "destination-without-service-instance"
  type           = "HTTP"
  proxy_type     = "Internet"
  url            = "https://myservice.example.com"
  authentication = "NoAuthentication"
  description    = "resource"
  subaccount_id  = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  additional_configuration = jsonencode({
    Abc = "good"
  })
}

#subaccount destination creation without service instance and without labels (additional configuration). 
resource "btp_subaccount_destination" "destination-without-additional-configuration" {
  name                = "destination-without-additional-configuration"
  type                = "HTTP"
  proxy_type          = "Internet"
  url                 = "https://myservice.example.com"
  authentication      = "NoAuthentication"
  description         = "resource"
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}

#subaccount destination creation without service instance and and labels (additional configuration). 
#Note: Auth properties are part of additional configuration.
resource "btp_subaccount_destination" "destination-with-additional-configuration" {
  name           = "destination-with-additional-configuration"
  type           = "HTTP"
  proxy_type     = "Internet"
  url            = "https://myservice.example.com"
  authentication = "OAuth2ClientCredentials"
  description    = "resource"
  subaccount_id  = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  additional_configuration = jsonencode({
    ClientID            = "<clientid>"
    ClientSecret        = "<clientsecret>"
    TokenServiceURL     = "https://tokenurl"
    TokenServiceURLType = "dedicated"
  })
}

#subaccount destination resource with RFC type.
resource "btp_subaccount_destination" "rfc_dest" {
  name          = "rfc_dest"
  type          = "RFC"
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

  additional_configuration = jsonencode({
    "jco.client.ashost"                     = "abcd"
    "jco.client.client"                     = "001"
    "jco.client.delta"                      = "1"
    "jco.client.network"                    = "LAN"
    "jco.client.passwd"                     = "Welcome1"
    "jco.client.serialization_format"       = "rowBased"
    "jco.client.sysnr"                      = "00"
    "jco.client.trace"                      = "0"
    "jco.client.user"                       = "SAPIPS"
    "jco.destination.auth_type"             = "CONFIGURED_USER"
    "jco.destination.pool_check_connection" = "0"
    "jco.destination.proxy_type"            = "OnPremise"
    "jco.destination.description"            = "RFC destination test"
  })
}

#subaccount destination resource with LDAP type.
resource "btp_subaccount_destination" "ldap_dest" {
  name          = "ldap_dest"
  type          = "LDAP"
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

  additional_configuration = jsonencode({
    "ldap.url"            = "ldap://ldap.example.com:389"
    "ldap.proxyType"      = "Internet"
    "ldap.description"    = "LDAP destination test"
    "ldap.authentication" = "BasicAuthentication"
    "ldap.user"           = "abc"
    "ldap.password"       = "abc"
  })
}

#subaccount destination resource with MAIL type.
resource "btp_subaccount_destination" "mail_dest" {
  name          = "mail_dest"
  type          = "TCP"
  proxy_type    = "Internet"
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

  additional_configuration = jsonencode({
    "mail.smtp.host"     = "smtp.example.com"
    "mail.smtp.port"     = "587"
    "mail.smtp.auth"     = "true" 
    "mail.description"   = "MAIL destination test"
    "mail.user"          = "user@example.com"
    "mail.password"      = "secret"
    "mail.transport.protocol" = "smtp"
  })
}

#subaccount destination resource with TCP type.
resource "btp_subaccount_destination" "tcp_dest" {
  name          = "tcp_dest"
  type          = "TCP"
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  proxy_type     = "OnPremise"
  description   = "TCP destination example"
  additional_configuration = jsonencode({
    "Address"= "host:1234"
  })
}