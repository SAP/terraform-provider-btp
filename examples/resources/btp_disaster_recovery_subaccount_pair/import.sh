# terraform import btp_disaster_recovery_subaccount_pair.<resource_name> '<subaccount_id>,<paired_subaccount_id>'
terraform import btp_disaster_recovery_subaccount_pair.dr_pair 'dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0, 2dc1ecf1-786c-4f92-91f2-26650ab3ad28'

# terraform import using id attribute in import block
import {
  to = btp_disaster_recovery_subaccount_pair.<resource_name>
  id = "<subaccount_id>,<paired_subaccount_id>"
}

# terraform import using resource identity
import {
to = btp_disaster_recovery_subaccount_pair.<resource_name>
identity = {
  subaccount_id        = "<subaccount_id>"
  paired_subaccount_id = "<paired_subaccount_id>"
  }
}
