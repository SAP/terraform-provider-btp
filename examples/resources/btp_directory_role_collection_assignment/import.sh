# terraform import btp_directory_role_collection_assignment.<resource_name> '<directory_id>,<role_collection_name>,<user_name>,<origin>'

terraform import btp_directory_role_collection_assignment.jd 'ddfc2206-5f11-48ed-a1ec-29010af70050,Directory Viewer,john.doe@mycompany.com,ldap'

# for group assignments use the group: prefix
terraform import btp_directory_role_collection_assignment.jd 'ddfc2206-5f11-48ed-a1ec-29010af70050,Directory Viewer,group:directory-viewer-group,my-idp'

# for attribute assignments use the attribute: prefix with / separator
terraform import btp_directory_role_collection_assignment.jd 'ddfc2206-5f11-48ed-a1ec-29010af70050,Directory Viewer,attribute:cost_center/1234,my-idp'

# terraform import using id attribute in import block

import {
  to = btp_directory_role_collection_assignment.<resource_name>
  id = "<directory_id>,<role_collection_name>,<user_name>,<origin>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher

# import a user assignment
import {
  to = btp_directory_role_collection_assignment.<resource_name>
  identity = {
    directory_id         = "<directory_id>"
    role_collection_name = "<role_collection_name>"
    user_name            = "<user_name>"
    origin               = "<origin>"
  }
}

# import a group assignment
import {
  to = btp_directory_role_collection_assignment.<resource_name>
  identity = {
    directory_id         = "<directory_id>"
    role_collection_name = "<role_collection_name>"
    group_name           = "<group_name>"
    origin               = "<origin>"
  }
}

# import an attribute assignment
import {
  to = btp_directory_role_collection_assignment.<resource_name>
  identity = {
    directory_id         = "<directory_id>"
    role_collection_name = "<role_collection_name>"
    attribute_name       = "<attribute_name>"
    attribute_value      = "<attribute_value>"
    origin               = "<origin>"
  }
}
