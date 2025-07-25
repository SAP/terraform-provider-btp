# terraform import btp_directory_role_collection.<resource_name> '<directory_id>,<name>'

terraform import btp_directory_role_collection.directory_viewer '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,Directory Viewer'

#terraform import using id attribute in import block

import {
  to = btp_directory_role_collection.<resource_name>
  id = "<directory_id>,<name>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher

import {
to = btp_directory_role_collection.<resource_name>
identity = {
  directory_id = "<directory_id>"
  name = "<name>"
  }
}
