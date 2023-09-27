package tfutils

import "github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"

func RemoveComputedlabels(labels servicemanager.ServiceManagerLabels) servicemanager.ServiceManagerLabels {
	// This method is intended to filter computed labels from service manager response
	// The removal is necessary as the computed label jeopardizes the Terraform state

	removeComputedLabelKeys := []string{"subaccount_id"}

	for _, keyToRemove := range removeComputedLabelKeys {
		delete(labels, keyToRemove)
	}

	return labels
}
