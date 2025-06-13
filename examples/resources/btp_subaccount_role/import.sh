# terraform import btp_subaccount_role.<resource_name> '<subaccount_id>,<name>,<role_template_name>,<app_id>'

terraform import btp_subaccount_role.subaccount_viewer '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,Subaccount Viewer,Subaccount_Viewer,cis-local!b2'

#terraform import using id attribute in import block

import {
  to = btp_subaccount_role.<resource_name>
  id = "<subaccount_id>,<name>,<role_template_name>,<app_id>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher

import {
to = btp_subaccount_role.<resource_name>
identity = {
  subaccount_id = "<subaccount_id>"
  name = "<name>"
  role_template_name = "<role_template_name>"
  app_id = "<app_id>"
  }
}
