# You created a environment instance for Kyma via the resource btp_subaccount_environment_instance
# The API Server URL of the Kyma instance should be displayed as output
# The labels attribute of the environment instance contains a JSON string with the following structure
# {"APIServerURL":"https://api.kyma.example.com","KubeconfigURL":"https://kyma.example.com/kubeconfig/ABC", "Name":"test-terraform-kyma"}

# This will return the value "https://api.kyma.example.com"
output "kyma_api_server_url" {
  description = "API Server URL of the Kyma environment instance"
  value       = provider::btp::extract_kyma_api_server_url(btp_subaccount_environment_instance.kyma_instance.labels)
}
