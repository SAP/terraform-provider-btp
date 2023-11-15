package xsuaa_settings

// TenantSettingsResp TenantSettingsResp
//
// swagger:model TenantSettingsResp
type TenantSettingsResp struct {

	// credential type infos
	CredentialTypeInfos []*Binding `json:"CredentialTypeInfos"`

	// Lists the custom e-mail domains supported by this tenant.
	// Example: internal.test, mail.invalid
	CustomEmailDomains []string `json:"customEmailDomains"`

	// The parameter displays the default identity provider (IdP) of the current tenant.
	// Example: sap.default
	DefaultIdp string `json:"defaultIdp,omitempty"`

	// By default, login pages of the service can't be framed by other applications in different domains for security reasons. The service trusts the domains listed here to embed the login page. The entire list can't exceed 2048 characters. For more information, see [Implications of Using IFrames](https://help.sap.com/docs/btp/sap-business-technology-platform/security-considerations-for-sap-authorization-and-trust-management-service#implications-of-using-iframes).
	// Example: https://store.example.com
	IframeDomains string `json:"iframeDomains,omitempty"`

	// links
	Links *LinksSettings `json:"links,omitempty"`

	// saml config settings
	SamlConfigSettings *SamlConfigSettingsResp `json:"samlConfigSettings,omitempty"`

	// token policy settings
	TokenPolicySettings *TokenPolicySettingsResp `json:"tokenPolicySettings,omitempty"`

	// Indicates whether the fallback at logon is enabled or not that if the logon ID provided in the token of the identity provider is unknown, the service attempts to log on the user with the e-mail address from the token. When false, the service attempts to create a missing user if user creation at logon is allowed. Note that before you can switch this parameter from false to true again, ensure that e-mail addresses are unique among your shadow users.
	TreatUsersWithSameEmailAsSameUser bool `json:"treatUsersWithSameEmailAsSameUser,omitempty"`
}
