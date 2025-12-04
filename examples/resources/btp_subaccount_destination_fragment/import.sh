# To import destination fragment on the subaccount level, use the following syntax:
# terraform import btp_subaccount_destination_fragment.<resource_name> '<subaccount_id>,<name>'
terraform import btp_subaccount_destination_fragment.sdf '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,test_fragment'

# To import destination fragment on the service instance level, use the following syntax:
# terraform import btp_subaccount_destination_fragment.<resource_name> '<subaccount_id>,<name>,<service_instance_id>'
terraform import btp_subaccount_destination_fragment.sdf '6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f,test_fragment,6a55f158-41b5-4e63-aa77-84089fa0ab98'

#terraform import using id attribute in import block
#On Subaccount Level
import {
  to = btp_subaccount_destination_fragment.<resource_name>
  id = "<subaccount_id>,<name>"
}
#On Service Instance Level
import {
  to = btp_subaccount_destination_fragment.<resource_name>
  id = "<subaccount_id>,<name>,<service_instance_id>"
}

# this resource supports import using identity attribute from Terraform version 1.12 or higher
#On Subaccount Level
import {
to =  btp_subaccount_destination_fragment.<resource_name>
identity = {
  name  = "<name>"
  subaccount_id = "<subaccount_id>"
  }
}

#On Service Instance Level
import {
to =  btp_subaccount_destination_fragment.<resource_name>
identity = {
  name  = "<name>"
  subaccount_id = "<subaccount_id>"
  service_instance_id = "<service_instance_id>"
  }
}
