package cis

type SubaccountHierarchyResponseObject struct {
	// Whether the subaccount can use beta services and applications.
	BetaEnabled bool `json:"betaEnabled"`
	// Details of the user that created the subaccount.
	CreatedBy string `json:"createdBy,omitempty"`
	// The date the subaccount was created. Dates and times are in UTC format.
	CreatedDate Time `json:"createdDate"`
	// (Deprecated) Contains information about the labels assigned to a specified subaccount. This field supports only single values per key and is now replaced by the string array \"labels\", which supports multiple values per key. The \"customProperties\" field returns only the first value of any label key that has multiple values assigned to it.
	CustomProperties []PropertyResponseObject `json:"customProperties,omitempty"`
	// A description of the subaccount for customer-facing UIs.
	Description string `json:"description"`
	// A descriptive name of the subaccount for customer-facing UIs.
	DisplayName string `json:"displayName"`
	// The unique ID of the subaccount's global account.
	GlobalAccountGUID string `json:"globalAccountGUID"`
	// Unique ID of the subaccount.
	Guid string `json:"guid"`
	// Contains information about the labels assigned to a specified subaccount. Labels are represented in a JSON array of key-value pairs; each key has up to 10 corresponding values. This field replaces the deprecated \"customProperties\" field, which supports only single values per key.
	Labels map[string][]string `json:"labels,omitempty"`
	// The date the subaccount was last modified. Dates and times are in UTC format.
	ModifiedDate Time `json:"modifiedDate,omitempty"`
	// The features of parent entity of the subaccount.
	ParentFeatures []string `json:"parentFeatures"`
	// The GUID of the subaccount’s parent entity. If the subaccount is located directly in the global account (not in a directory), then this is the GUID of the global account.
	ParentGUID string `json:"parentGUID"`
	// The region in which the subaccount was created.
	Region string `json:"region"`
	// The current state of the subaccount.
	State string `json:"state"`
	// Information about the state of the subaccount.
	StateMessage string `json:"stateMessage,omitempty"`
	// The subdomain that becomes part of the path used to access the authorization tenant of the subaccount. Must be unique within the defined region. Use only letters (a-z), digits (0-9), and hyphens (not at the start or end). Maximum length is 63 characters. Cannot be changed after the subaccount has been created.
	Subdomain string `json:"subdomain"`
	// The technical name of the subaccount. Refers to: (1) the platform-based account name for Neo subaccounts, or (2) the account identifier (tenant ID) in XSUAA for multi-environment subaccounts.
	TechnicalName string `json:"technicalName"`
	// Whether the subaccount is used for production purposes. This flag can help your cloud operator to take appropriate action when handling incidents that are related to mission-critical accounts in production systems. Do not apply for subaccounts that are used for non-production purposes, such as development, testing, and demos. Applying this setting this does not modify the subaccount. * <b>UNSET:</b> Global account or subaccount admin has not set the production-relevancy flag. Default value. * <b>NOT_USED_FOR_PRODUCTION:</b> Subaccount is not used for production purposes. * <b>USED_FOR_PRODUCTION:</b> Subaccount is used for production purposes.
	UsedForProduction string `json:"usedForProduction"`
}