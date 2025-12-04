package connectivity

type DestinationTrust struct {
	Name                string                 `json:"Name,omitempty"`
	BaseURL             string                 `json:"baseURL,omitempty"`
	Active              bool                   `json:"active,omitempty"`
	Expiration          int64                  `json:"expirationTimestamp,omitempty"`
	Owner               *DestinationTrustOwner `json:"Owner,omitempty"`
	GeneratedOn         string                 `json:"generatedOn,omitempty"`
	X509PublicKeyBase64 string                 `json:"x509PublicKeyBase64,omitempty"`
}

type DestinationTrustOwner struct {
	InstanceID   string `json:"InstanceId,omitempty"`
	SubaccountID string `json:"SubaccountId,omitempty"`
}
