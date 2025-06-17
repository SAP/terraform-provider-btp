# terraform import btp_globalaccount_trust_configuration.<resource_name> <origin>

terraform import btp_globalaccount_trust_configuration.trust sap.custom

#terraform import using id attribute in import block

import {
  to = btp_globalaccount_trust_configuration.<resource_name>
  id = "<origin>"
}
