# terraform import btp_directory_entitlement.<resource_name> <directory_id>,<service_name>,<plan_name>

terraform import btp_directory_entitlement.alert_notification_service 6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,alert-notification,free

#terraform import using id attribute in import block

import {
  to = btp_directory_entitlement.<resource_name>
  id = "<directory_id>,<service_name>,<plan_name>"
}

# terraform import using identity attribute in import block (supported in terraform version 1.12 or later)

import {
to = btp_directory_entitlement.<resource_name>
identity = {
  directory_id = "<directory_id>"
  service_name = "<service_name>"
  plan_name = "<plan_name>"
  }
}
