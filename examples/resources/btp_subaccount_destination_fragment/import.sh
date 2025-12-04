# this resource supports import using identity attribute from Terraform version 1.12 or higher

import {
to =  btp_subaccount_destination_fragment.<resource_name>
identity = {
  name  = "<name>"
  subaccount_id = "<subaccount_id>"
  service_instance_id = "<service_instance_id>"
  }
}
