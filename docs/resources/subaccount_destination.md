---
page_title: "btp_subaccount_destination Resource - terraform-provider-btp"
subcategory: ""
description: |-
  Manages a destination in a SAP BTP subaccount or in the scope of a specific service instance.
  This resource must be preferred only for HTTP destinations. We recommend using the resource 'btp_subaccount_destination_generic' to accommodate all types.
  Tip:
  You must have the appropriate connectivity and destination permissions, such as:
  Subaccount Administrator
  Destination Administrator
  Connectivity and Destination Administrator
  Scope:
  Subaccount-level destination: Specify only the 'subaccount_id' and 'name' attribute.Service instance-level destination: Specify the 'subaccount_id', 'service_instance_id' and 'name' attributes.
  Notes:
  'service_instance_id' is optional. When omitted, the destination is created at the subaccount level.
---

# btp_subaccount_destination (Resource)

Manages a destination in a SAP BTP subaccount or in the scope of a specific service instance.
							  This resource must be preferred only for HTTP destinations. We recommend using the resource 'btp_subaccount_destination_generic' to accommodate all types.

__Tip:__
You must have the appropriate connectivity and destination permissions, such as:

Subaccount Administrator
Destination Administrator
Connectivity and Destination Administrator
__Scope:__
- **Subaccount-level destination**: Specify only the 'subaccount_id' and 'name' attribute.
- **Service instance-level destination**: Specify the 'subaccount_id', 'service_instance_id' and 'name' attributes.

__Notes:__
- 'service_instance_id' is optional. When omitted, the destination is created at the subaccount level.

## Example Usage

```terraform
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
    "jco.client.ashost"                     = "abc"
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
    "jco.destination.description"           = "RFC destination test"
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
  type          = "MAIL"
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

  additional_configuration = jsonencode({
    "mail.smtp.host"          = "smtp.example.com"
    "mail.smtp.port"          = "587"
    "mail.smtp.auth"          = "true"
    "mail.description"        = "MAIL destination test"
    "mail.user"               = "user@example.com"
    "mail.password"           = "secret"
    "mail.transport.protocol" = "smtp"
    "ProxyType"               = "Internet"
    "Authentication"          = "BasicAuthentication"
  })
}

#subaccount destination resource with TCP type.
resource "btp_subaccount_destination" "tcp_dest" {
  name          = "tcp_dest"
  type          = "TCP"
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  additional_configuration = jsonencode({
    "Address"     = "host:1234",
    "ProxyType"   = "OnPremise",
    "Description" = "TCP destination example"
  })
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The descriptive name of the destination for subaccount
- `subaccount_id` (String) The ID of the subaccount.
- `type` (String) The type of request from destination.

### Optional

- `additional_configuration` (String) The additional configuration parameters for the destination.
- `authentication` (String) The authentication of the destination.
- `description` (String) The description of the destination.
- `proxy_type` (String) The proxytype of the destination.
- `service_instance_id` (String) The service instance that becomes part of the path used to access the destination of the subaccount.
- `url` (String) The url of the destination.

### Read-Only

- `creation_time` (String) The date and time when the resource was created in
- `etag` (String) The etag for the destination resource
- `id` (String, Deprecated) The ID of the destination used for import operations.
- `modification_time` (String) The date and time when the resource was modified

## Import

Import is supported using the following syntax:

```terraform
# To import destination on the subaccount level, use the following syntax:
# terraform import btp_subaccount_destination.<resource_name> '<subaccount_id>,<name>'
terraform import btp_subaccount_destination.abc '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,test'

# To import destination on the service instance level, use the following syntax:
# terraform import btp_subaccount_destination.<resource_name> '<subaccount_id>,<name>,<service_instance_id>'
terraform import btp_subaccount_destination.abc '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,test,6a55f158-41b5-4e63-aa77-84089fa0ab98'

#terraform import using id attribute in import block
#On Subaccount Level
import {
  to = btp_subaccount_destination.<resource_name>
  id = "<subaccount_id>,<name>"
}
#On Service Instance Level
import {
  to = btp_subaccount_destination.<resource_name>
  id = "<subaccount_id>,<name>,<service_instance_id>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher
#On Subaccount Level
import {
to =  btp_subaccount_destination.<resource_name>
identity = {
  name          = "<name>"
  subaccount_id = "<subaccount_id>"
  }
}

#On Service Instance Level
import {
to =  btp_subaccount_destination.<resource_name>
identity = {
  name                = "<name>"
  subaccount_id       = "<subaccount_id>"
  service_instance_id = "<service_instance_id>"
  }
}
```
