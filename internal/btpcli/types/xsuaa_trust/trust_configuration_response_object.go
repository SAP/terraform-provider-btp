package xsuaa_trust

type TrustConfigurationResponseObject struct {
	Name                         string `json:"name,omitempty"`
	OriginKey                    string `json:"originKey,omitempty"`
	TypeOfTrust                  string `json:"typeOfTrust,omitempty"`
	Status                       string `json:"status,omitempty"`
	Description                  string `json:"description,omitempty"`
	IdentityProvider             string `json:"identityProvider,omitempty"`
	Domain                       string `json:"domain,omitempty"`
	LinkTextForUserLogon         string `json:"linkTextForUserLogon,omitempty"`
	AvailableForUserLogon        string `json:"availableForUserLogon,omitempty"`
	CreateShadowUsersDuringLogon string `json:"createShadowUsersDuringLogon,omitempty"`
	Protocol                     string `json:"protocol,omitempty"`
	ReadOnly                     bool   `json:"readOnly,omitempty"`
}
