# terraform import btp_subaccount_role_collection_assignment.<resource_name> '<subaccount_id>,<role_collection_name>,<user_name>,<origin>'

terraform import btp_subaccount_role_collection_assignment.jd '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,Subaccount Viewer,john.doe@mycompany.com,ldap'

# for group assignments use the group: prefix
terraform import btp_subaccount_role_collection_assignment.jd '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,Subaccount Viewer,group:subaccount-viewer-group,my-idp'

# for attribute assignments use the attribute: prefix with / separator
terraform import btp_subaccount_role_collection_assignment.jd '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,Subaccount Viewer,attribute:cost_center/1234,my-idp'

# terraform import using id attribute in import block

import {
  to = btp_subaccount_role_collection_assignment.<resource_name>
  id = "<subaccount_id>,<role_collection_name>,<user_name>,<origin>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher

# import a user assignment
import {
  to = btp_subaccount_role_collection_assignment.<resource_name>
  identity = {
    subaccount_id        = "<subaccount_id>"
    role_collection_name = "<role_collection_name>"
    user_name            = "<user_name>"
    origin               = "<origin>"
  }
}

# import a group assignment
import {
  to = btp_subaccount_role_collection_assignment.<resource_name>
  identity = {
    subaccount_id        = "<subaccount_id>"
    role_collection_name = "<role_collection_name>"
    group_name           = "<group_name>"
    origin               = "<origin>"
  }
}

# import an attribute assignment
import {
  to = btp_subaccount_role_collection_assignment.<resource_name>
  identity = {
    subaccount_id        = "<subaccount_id>"
    role_collection_name = "<role_collection_name>"
    attribute_name       = "<attribute_name>"
    attribute_value      = "<attribute_value>"
    origin               = "<origin>"
  }
}
