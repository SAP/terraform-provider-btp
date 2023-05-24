package xsuaa_trust

type TrustConfigurationResponseObject struct {
	// The name of the identity provider.
	Name string `json:"name,omitempty"`
	// The origin of the identity provider.
	OriginKey   string `json:"originKey,omitempty"`
	TypeOfTrust string `json:"typeOfTrust,omitempty"`
	// Whether the identity provider is currently active or not.
	Status string `json:"status,omitempty"`
	// A description for the identity provider.
	Description string `json:"description,omitempty"`
	// The protocol used to establish trust with the identity provider.
	Protocol string `json:"protocol,omitempty"`
	// Whether the trust configuration can be modified.
	ReadOnly bool `json"readOnly,omitempty"`
	// Name of the identity provider
	IdentityProvider string `json:"identityProvider,omitempty"`
}
