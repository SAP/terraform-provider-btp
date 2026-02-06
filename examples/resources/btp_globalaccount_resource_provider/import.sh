# terraform import btp_globalaccount_resource_provider.<resource_name> <resource_provider>,<unique_technical_name>

terraform import btp_globalaccount_resource_provider.azure AZURE,my_azure_provider

#terraform import using id attribute in import block

import {
  to = btp_globalaccount_resource_provider.<resource_name>
  id = "<resource_provider>,<unique_technical_name>"
}

import {
  to = btp_globalaccount_resource_provider.<resource_name>
  identity = {
    provider_type  = "<resource_provider>"
    technical_name = "<unique_technical_name>"
  }
}
