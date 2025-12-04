# You created a environment instance for Kyma via the resource btp_subaccount_environment_instance
# The Download URL of the kubeconfig for the Kyma instance should be displayed as output
# The labels attribute of the environment instance contains a JSON string with the following structure
# {"APIServerURL":"https://api.kyma.example.com","KubeconfigURL":"https://kyma.example.com/kubeconfig/ABC", "Name":"test-terraform-kyma"}

# This will return the value "https://kyma.example.com/kubeconfig/ABC"
output "kyma_kubeconfig_url" {
  description = "Kubeconfig URL of the Kyma environment instance"
  value       = provider::btp::extract_kyma_kubeconfig_url(btp_subaccount_environment_instance.kyma_instance.labels)
}
