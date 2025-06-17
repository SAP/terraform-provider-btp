# terraform import btp_subaccount.<resource_name> <subaccount_id>

terraform import btp_subaccount.my_project 6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f

#terraform import using id attribute in import block

import {
  to = btp_subaccount.<resource_name>
  id = "<subaccount_id>"
}
