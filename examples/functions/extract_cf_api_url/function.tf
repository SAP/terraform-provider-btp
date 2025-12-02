# You created a environment instance for Cloud Foundry via the resource btp_subaccount_environment_instance
# The API URL of the Cloud Foundry instance should be displayed as output
# The labels attribute of the environment instance contains a JSON string with the following structure
# {"API Endpoint":"https://api.cf.example.com","Org Name":"example","Org ID":"8d818824-394a-abcd-0815-7a3c8ce93e57","Org Memory Limit":"1000MB"}

# This will return the value "https://api.cf.example.com"
output "cf_api_url" {
  description = "API URL of the Cloud Foundry environment instance"
  value = provider::btp::extract_cf_api_url(btp_subaccount_environment_instance.cf_instance.labels)
}
