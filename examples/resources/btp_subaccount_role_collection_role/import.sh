# terraform import btp_subaccount_role_collection_role.<resource_name> '<subaccount_id>,<name>,<role_name>,<role_template_app_id>,<role_template_name>'

terraform import btp_subaccount_role_collection_role.destination_admin '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,Destination Administrator,Destination Administrator,destination-xsappname!b5,Destination_Administrator'

#terraform import using id attribute in import block

import {
  to = btp_subaccount_role_collection_role.<resource_name>
  id = "<subaccount_id>,<name>,<role_name>,<role_template_app_id>,<role_template_name>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher

import {
  to       = btp_subaccount_role_collection_role.<resource_name>
  identity = {
    name                 = "<collection_name>"
    role_name            = "<role_name>"
    role_template_app_id = "<role_template_app_id>"
    role_template_name   = "<role_template_name>"
    subaccount_id        = "<subaccount_id>"
  }
}