package xsuaa_settings

// TokenPolicySettingsResp TokenPolicySettingsResp
//
// swagger:model TokenPolicySettingsResp
type TokenPolicySettingsResp struct {

	// Time in seconds between when a access token is issued and when it expires. The value ranges from 1800 seconds to 86,400 seconds, in other words, from 30 minutes to 24 hours. Keep token validity as short as possible, but not less than 30 minutes. The default value is 43,000 seconds or 12 hours. The value `-1` means that the token uses the default setting. Token policy settings apply to all service instances in the subaccount that haven't set a specific value in the application security descriptor (xs-security.json). For more information, see [Setting Token Policy](https://help.sap.com/docs/BTP/65de2977205c403bbc107264b8eccf4b/f117cab6b92d438cb2a0b5204713994b.html#setting-token-policy).
	AccessTokenValidity int32 `json:"accessTokenValidity,omitempty"`

	// The ID of the key to use for signing metadata and assertions.
	// Example: default-jwt-key--9988843812
	ActiveKeyID string `json:"activeKeyId,omitempty"`

	// key ids
	KeyIds []string `json:"keyIds"`

	// If true, the service only issues one refresh token per client_id and user_id combination.
	RefreshTokenUnique bool `json:"refreshTokenUnique,omitempty"`

	// Time in seconds between when a refresh token is issued and when it expires. The value ranges from 1800 seconds to 31,536,000 seconds, in other words, from 30 minutes to one year. The validity of refresh tokens must be longer than the validity for access tokens. The system never issues refresh tokens if the validity is shorter. The default value is 604,800 seconds or 7 days. The value `-1` means that the token uses the default setting. Token policy settings apply to all service instances in the subaccount that haven't set a specific value in the application security descriptor (xs-security.json). For more information, see [Setting Token Policy](https://help.sap.com/docs/BTP/65de2977205c403bbc107264b8eccf4b/f117cab6b92d438cb2a0b5204713994b.html#setting-token-policy).
	RefreshTokenValidity int32 `json:"refreshTokenValidity,omitempty"`
}
