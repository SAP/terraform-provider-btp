# You created a environment instance for Cloud Foundry via the resource btp_subaccount_environment_instance
# The Org ID of the Cloud Foundry instance should be displayed as output
# The labels attribute of the environment instance contains a JSON string with the following structure
# {"API Endpoint":"https://api.cf.example.com","Org Name":"example","Org ID":"8d818824-394a-abcd-0815-7a3c8ce93e57","Org Memory Limit":"1000MB"}


# This will return the value "8d818824-394a-abcd-0815-7a3c8ce93e57"
output "cf_org_id" {
  description = "Org ID of the Cloud Foundry environment instance"
  value = provider::btp::extract_cf_org_id(btp_subaccount_environment_instance.cf_instance.labels)
}
