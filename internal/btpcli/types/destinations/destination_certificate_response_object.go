package destinations

type DestinationCertificateResponseObject struct {
	Name     string                     `json:"name,omitempty"`
	Nodes    []NodeResponseObject       `json:"nodes,omitempty"`
	Creation CreationDataResponseObject `json:"creation"`
}

type NodeResponseObject struct {
	Type        string `json:"type,omitempty" tfsdk:"type"`
	Format      string `json:"format,omitempty" tfsdk:"format"`
	Algorithm   string `json:"algorithm,omitempty" tfsdk:"algorithm"`
	Alias       string `json:"alias,omitempty" tfsdk:"alias"`
	Subject     string `json:"subject,omitempty" tfsdk:"subject"`
	Issuer      string `json:"issuer,omitempty" tfsdk:"issuer"`
	CommonName  string `json:"commonName,omitempty" tfsdk:"common_name"`
	NotBefore   string `json:"notBefore,omitempty" tfsdk:"not_before"`
	NotAfter    string `json:"notAfter,omitempty" tfsdk:"not_after"`
	Certificate string `json:"certificate,omitempty" tfsdk:"certificate"`
}

type CreationDataResponseObject struct {
	GenerationMethod  string `json:"generation_method,omitempty" tfsdk:"generation_method"`
	CommonName        string `json:"common_name,omitempty" tfsdk:"common_name"`
	HasPassword       bool   `json:"has_password,omitempty" tfsdk:"has_password"`
	AutoRenew         bool   `json:"auto_renew,omitempty" tfsdk:"auto_renew"`
	ValidityDuration  string `json:"validity_duration,omitempty" tfsdk:"validity_duration"`
	ValidityTimeUnits string `json:"validity_time_units,omitempty" tfsdk:"validity_time_units"`
}
