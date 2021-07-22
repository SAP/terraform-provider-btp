package provisioning

import (
	"encoding/json"
)

type ResourceProviderResponseObject struct {
	// Unique technical name of the resource.
	ResourceTechnicalName string `json:"technicalName,omitempty"`
	// Type of the resource.
	ResourceType string `json:"resourceType,omitempty"`
	// Provider of the requested resource. For example, IaaS provider: AWS.
	ResourceProvider string `json:"resourceProvider,omitempty"`
	// Descriptive name of the resource for customer-facing UIs.
	DisplayName string `json:"displayName,omitempty"`
	// Description of the resource.
	Description string `json"description,omitempty"`
	// Any relevant information about the resource that is not provided by other parameter values.
	AdditionalInfo *json.RawMessage `json:"additionalInfo,omitempty"`
}
