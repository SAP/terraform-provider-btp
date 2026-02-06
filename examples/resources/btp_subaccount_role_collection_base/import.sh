# terraform import btp_subaccount_role_collection_base.<resource_name> '<subaccount_id>,<name>'

terraform import btp_subaccount_role_collection_base.destination_admin '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,Destination Administrator'

#terraform import using id attribute in import block

import {
  to = btp_subaccount_role_collection_base.<resource_name>
  id = "<subaccount_id>,<name>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher

import {
to = btp_subaccount_role_collection_base.<resource_name>
identity = {
  subaccount_id = "<subaccount_id>"
  name          = "<name>"
  }
}
