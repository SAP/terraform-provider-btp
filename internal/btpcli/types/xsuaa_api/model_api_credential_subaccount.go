package xsuaa_api

type ApiCredential struct {
	TenantMode         string `json:"tenantmode"`
	SubaccountId       string `json:"subaccountid"`
	CredentialType     string `json:"credential-type"`
	ClientId           string `json:"clientid"`
	ClientX509Enabled  bool   `json:"clientx509enabled,omitempty"`
	Certificate        string `json:"certificate,omitempty"`
	CertUrl            string `json:"certurl,omitempty"`
	CertificatePinning bool   `json:"certificate-pinning,omitempty"`
	Key                string `json:"key,omitempty"`
	TokenUrl           string `json:"tokenurl"`
	XsAppname          string `json:"xsappname"`
	ClientSecret       string `json:"clientsecret,omitempty"`
	ServiceInstanceId  string `json:"serviceInstanceId"`
	Url                string `json:"url"`
	UaaDomain          string `json:"uaadomain"`
	ApiUrl             string `json:"apiurl"`
	IdentityZone       string `json:"identityzone"`
	IdentityZoneId     string `json:"identityzoneid"`
	TenantId           string `json:"tenantid"`
	Name               string `json:"name"`
	ZoneId             string `json:"zoneid"`
	ReadOnly           bool   `json:"read-only"`
}
