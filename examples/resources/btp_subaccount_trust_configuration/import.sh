# terraform import btp_subaccount_trust_configuration.<resource_name> <subaccount_id>,<origin>

terraform import btp_subaccount_trust_configuration.trust 6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,sap.custom

#terraform import using id attribute in import block

import {
  to = btp_subaccount_trust_configuration.<resource_name>
  id = "<subaccount_id>,<origin>"
}

import {
  to = btp_subaccount_trust_configuration.<resource_name>
  identity = {
    subaccount_id = "<subaccount_id>"
    origin        = "<origin>"
  }
}