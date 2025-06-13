# terraform import btp_directory_role.<resource_name> '<directory_id>,<name>,<role_template_name>,<app_id>'

terraform import btp_directory_role.directory_viewer '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,Directory Viewer,Directory_Viewer,cis-central!b13'

#terraform import using id attribute in import block

import {
  to = btp_directory_role.<resource_name>
  id = "<directory_id>,<name>,<role_template_name>,<app_id>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher

import {
to = btp_directory_role.<resource_name>
identity = {
  directory_id = "<directory_id>"
  name = "<name>"
  role_template_name = "<role_template_name>"
  app_id = "<app_id>"
  }
}
