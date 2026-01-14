# terraform import btp_subaccount_environment_instance.<resource_name> <subaccount_id>,<environment_instance_id>

terraform import btp_subaccount_environment_instance.cloudfoundry 6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,FD9BB73F-F663-4284-A50B-D72EC24FC4E1

#terraform import using id attribute in import block

import {
  to = btp_subaccount_environment_instance.<resource_name>
  id = "subaccount_id>,<environment_instance_id>"
}

import {
  to = btp_subaccount_environment_instance.<resource_name>
  identity = {
    subaccount_id = "<subaccount_id>"
    id = "<environment_instance_id>"
  }
}