# terraform import btp_globalaccount_role.<resource_name> '<name>,<role_template_name>,<app_id>'

terraform import btp_globalaccount_role.globalaccount_auditor 'User and Role Auditor,xsuaa_auditor,xsuaa!t2'

#terraform import using id attribute in import block

import {
  to =  btp_globalaccount_role.<resource_name>
  id = "<name>,<role_template_name>,<app_id>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher

import {
to =  btp_globalaccount_role.<resource_name>
identity = {
  name = "<name>"
  role_template_name = "<role_template_name>"
  app_id = "<app_id>"
  }
}
