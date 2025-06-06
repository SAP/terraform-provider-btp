# terraform import btp_globalaccount_role_collection.<resource_name> '<name>'

terraform import btp_globalaccount_role_collection.globalaccount_viewer 'Global Account Viewer'

#terraform import using id attribute in import block

import {
  to =  btp_globalaccount_role_collection.<resource_name>
  id = "<name>"
}

# terraform import using identity attribute in import block (supported in terraform version 1.12 or later)

import {
to = btp_globalaccount_role_collection.<resource_name>
identity = {
  name = "<name>"
}