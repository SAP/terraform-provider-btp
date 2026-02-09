#subaccount destination without auth
resource "btp_subaccount_destination_generic" "http_dest" {
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  destination_configuration = jsonencode({
    "Name"           = "project"
    "Type"           = "HTTP"
    "ProxyType"      = "Internet"
    "URL"            = "https://myservice.example.com"
    "Authentication" = "NoAuthentication"
    "Description"    = "trial destination of basic usecase with service instance"

  })
}

#subaccount destination creation without service instance and without labels (additional configuration).
resource "btp_subaccount_destination_generic" "destination" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  destination_configuration = jsonencode({
    "Name"           = "project"
    "Type"           = "HTTP"
    "ProxyType"      = "Internet"
    "URL"            = "https://myservice.example.com"
    "Authentication" = "NoAuthentication"
    "Description"    = "trial destination of basic usecase "
  })
}

#subaccount destination creation without service instance and and labels (additional configuration).
#Note: Auth properties are part of additional configuration.
resource "btp_subaccount_destination_generic" "http_dest_with_destination_configuration_auth" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  destination_configuration = jsonencode({
    "Name"            = "project_with_auth"
    "Type"            = "HTTP"
    "clientId"        = "abc"
    "tokenServiceURL" = "https://myservice.example.com"
    "ProxyType"       = "Internet"
    "URL"             = "https://myservice.example.com"
    "Authentication"  = "OAuth2ClientCredentials"
    "Description"     = "trial destination of basic usecase with service instance and with addditional variables update"
  })
}

#subaccount destination rfc type
resource "btp_subaccount_destination_generic" "rfc_dest" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  destination_configuration = jsonencode({
    "Name"                                  = "rfc_dest"
    "Type"                                  = "RFC"
    "jco.client.ashost"                     = "va4hci"
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
    "jco.destination.description"           = "RFC destination test update"
  })
}

#subaccount destination ldap type
resource "btp_subaccount_destination_generic" "ldap_dest" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  destination_configuration = jsonencode({
    "Name"                = "ldap_dest"
    "Type"                = "LDAP"
    "ldap.url"            = "ldap://ldap.example.com:389"
    "ldap.proxyType"      = "Internet"
    "ldap.description"    = "LDAP destination test update"
    "ldap.authentication" = "BasicAuthentication"
    "ldap.user"           = "abc"
    "ldap.password"       = "abc"
  })
}

#subaccount destination mail type
resource "btp_subaccount_destination_generic" "mail_dest" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  destination_configuration = jsonencode({
    "Name"             = "mail_dest"
    "Type"             = "MAIL"
    "Authentication"   = "BasicAuthentication"
    "ProxyType"        = "OnPremise"
    "mail.description" = "MAIL destination test update"
    "mail.user"        = "user@example.com"
    "mail.password"    = "secret"
  })
}

#subaccount destination tcp type
resource "btp_subaccount_destination_generic" "tcp_dest" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  destination_configuration = jsonencode({
    "Name"        = "tcp_dest"
    "Type"        = "TCP"
    "Address"     = "host:1234"
    "ProxyType"   = "OnPremise"
    "Description" = "TCP destination example update"
  })
}
