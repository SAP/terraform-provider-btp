package xsuaa_settings

// SamlConfigSettingsResp SamlConfigSettingsResp
//
// swagger:model SamlConfigSettingsResp
type SamlConfigSettingsResp struct {

	// The ID of the key to be used for signing metadata and assertions.
	// Example: default-saml-key-99999
	ActiveKeyID string `json:"activeKeyId,omitempty"`

	// If true, this zone doesn't validate the `InResponseToField` part of an incoming identity provider assertion.
	DisableInResponseToCheck bool `json:"disableInResponseToCheck,omitempty"`

	// The parameter contains a globally unique name for an identity provider or a service provider.
	// Example: https://example-tenant.authentication.eu10.hana.ondemand.com
	EntityID string `json:"entityID,omitempty"`

	// keys
	Keys *SamlKey `json:"keys,omitempty"`
}
