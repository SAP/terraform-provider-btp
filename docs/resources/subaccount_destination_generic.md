---
page_title: "btp_subaccount_destination_generic Resource - terraform-provider-btp"
subcategory: ""
description: |-
  Manages a destination in a SAP BTP subaccount or in the scope of a specific service instance.
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

# btp_subaccount_destination_generic (Resource)

Manages a destination in a SAP BTP subaccount or in the scope of a specific service instance.
		
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
#subaccount destination without auth
resource "btp_subaccount_destination_generic" "http_dest" {
  subaccount_id  = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  service_instance_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  additional_configuration = jsonencode({
    "Name"           = "project"
    "Type"           = "HTTP"
    "ProxyType"      = "Internet"
    "URL"            = "https://myservice.example.com"
    "Authentication" = "NoAuthentication"
    "Description"    = "trial destination of basic usecase with service instance"
  
  })
}

#subaccount destination creation without service instance and without labels (additional configuration). 
resource "btp_subaccount_destination" "destination" {
  subaccount_id       = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
    additional_configuration = jsonencode({
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
resource "btp_subaccount_destination_generic" "http_dest_with_additional_configuration_auth" {
  subaccount_id     = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  additional_configuration = jsonencode({
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
  additional_configuration = jsonencode({
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
  additional_configuration = jsonencode({
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
  additional_configuration = jsonencode({
    "Name"                    = "mail_dest"
    "Type"                    = "MAIL"
    "Authentication"          = "BasicAuthentication"
    "ProxyType"               = "OnPremise"
    "mail.description"        = "MAIL destination test update"
    "mail.user"               = "user@example.com"
    "mail.password"           = "secret"
  })
}

#subaccount destination tcp type
resource "btp_subaccount_destination_generic" "tcp_dest" {
  subaccount_id = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  additional_configuration = jsonencode({
    "Name"          = "tcp_dest"
    "Type"          = "TCP"
    "Address"       = "host:1234"
    "ProxyType"     = "OnPremise"
    "Description"   = "TCP destination example update"
  })
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `additional_configuration` (String) The additional configuration parameters for the destination.
- `subaccount_id` (String) The ID of the subaccount.

### Optional

- `service_instance_id` (String) The service instance that becomes part of the path used to access the destination of the subaccount.

### Read-Only

- `creation_time` (String) The date and time when the resource was created in
- `etag` (String) The etag for the destination resource
- `id` (String, Deprecated) The ID of the destination used for import operations.
- `modification_time` (String) The date and time when the resource was modified
- `name` (String) The descriptive name of the destination for subaccount

## Import

Import is supported using the following syntax:

```terraform
# To import destination on the subaccount level, use the following syntax:
# terraform import btp_subaccount_destination_generic.<resource_name> '<subaccount_id>,<name>'
terraform import btp_subaccount_destination_generic.abc '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,test'

# To import destination on the service instance level, use the following syntax:
# terraform import btp_subaccount_destination_generic.<resource_name> '<subaccount_id>,<name>,<service_instance_id>'
terraform import btp_subaccount_destination_generic.abc '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,test,6a55f158-41b5-4e63-aa77-84089fa0ab98'

#terraform import using id attribute in import block
#On Subaccount Level
import {
  to = btp_subaccount_destination_generic.<resource_name>
  id = "<subaccount_id>,<name>"
}
#On Service Instance Level
import {
  to = btp_subaccount_destination_generic.<resource_name>
  id = "<subaccount_id>,<name>,<service_instance_id>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher
#On Subaccount Level
import {
to =  btp_subaccount_destination_generic.<resource_name>
identity = {
  name          = "<name>"
  subaccount_id = "<subaccount_id>"
  }
}

#On Service Instance Level
import {
to =  btp_subaccount_destination_generic.<resource_name>
identity = {
  name                = "<name>"
  subaccount_id       = "<subaccount_id>"
  service_instance_id = "<service_instance_id>"
  }
}
```
