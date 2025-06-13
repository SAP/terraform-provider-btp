# terraform import btp_subaccount_entitlement.<resource_name> <subaccount_id>,<service_name>,<plan_name>

terraform import btp_subaccount_entitlement.alert_notification_service 6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,alert-notification,free

#terraform import using id attribute in import block

import {
  to = btp_subaccount_entitlement.<resource_name>
  id = "<subaccount_id>,<service_name>,<plan_name>"
}

# terraform import using identity attribute in import block (supported in terraform version 1.12 or later)

import {
to = btp_subaccount_entitlement.<resource_name>
identity = {
  subaccount_id = "<subaccount_id>"
  service_name = "<service_name>"
  plan_name = "<plan_name>"
  }
}
