package xsuaa_settings

// Binding binding
type Binding struct {

	// app Id
	// Example: product-list!t1000
	AppID string `json:"appId,omitempty"`

	// binding Id
	// Example: 3a3a9aa2-bf44-4cd3-904a-b22d13a17fa4
	BindingID string `json:"bindingId,omitempty"`

	// credential type
	// Enum: [instance-secret binding-secret X.509]
	CredentialType string `json:"credentialType,omitempty"`

	// service Id
	// Example: a53672a7-7f94-48e8-bd24-0099b822abd7
	ServiceID string `json:"serviceId,omitempty"`

	// tenant Id
	// Example: 4a3a6a53-ae93-4450-ba12-4becb1c345ab
	TenantID string `json:"tenantId,omitempty"`
}
