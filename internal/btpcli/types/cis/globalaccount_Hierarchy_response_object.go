package cis

type GlobalAccountHierarchyResponseObject struct {
	// Specifies if global account is backward-compliant for EU access.
	BackwardCompliantEU bool `json:"backwardCompliantEU,omitempty"`
	// The list of directories associated with the specified global account.
	Children []DirectoryHierarchyResponseObject `json:"children,omitempty"`
	// The type of the commercial contract that was signed.
	CommercialModel string `json:"commercialModel"`
	// Whether the customer of the global account pays only for services that they actually use (consumption-based) or pay for subscribed services at a fixed cost irrespective of consumption (subscription-based). * <b>TRUE:</b> Consumption-based commercial model. * <b>FALSE:</b> Subscription-based commercial model.
	ConsumptionBased bool `json:"consumptionBased"`
	// The status of the customer contract and its associated root global account. * <b>ACTIVE:</b> The customer contract and its associated global account is currently active. * <b>PENDING_TERMINATION:</b> A termination process has been triggered for a customer contract (the customer contract has expired, or a customer has given notification that they wish to terminate their contract), and the global account is currently in the validation period. The customer can still access their global account until the end of the validation period. * <b>SUSPENDED:</b> For enterprise accounts, specifies that the customer's global account is currently in the grace period of the termination process. Access to the global account by the customer is blocked. No data is deleted until the deletion date is reached at the end of the grace period. For trial accounts, specifies that the account is suspended, and the account owner has not yet extended the trial period.
	ContractStatus string `json:"contractStatus,omitempty"`
	// The number of the cost center that is charged for the creation and usage of the global account. This is a duplicate property used for backward compatibility; the cost center is also stored in costObjectId. This property must be null if the global account is tied to an internal order or Work Breakdown Structure element.
	CostCenter string `json:"costCenter,omitempty"`
	// The number or code of the cost center, internal order, or Work Breakdown Structure element that is charged for the creation and usage of the global account. The type of the cost object must be configured in costObjectType.
	CostObjectId string `json:"costObjectId,omitempty"`
	// The type of accounting assignment object that is associated with the global account owner and used to charge for the creation and usage of the global account. Support types: COST_CENTER, INTERNAL_ORDER, WBS_ELEMENT. The number or code of the specified cost object is defined in costObjectId. For a cost object of type 'cost center', the value is also configured in costCenter for backward compatibility purposes.
	CostObjectType string `json:"costObjectType,omitempty"`
	// The date the global account was created. Dates and times are in UTC format.
	CreatedDate Time `json:"createdDate"`
	// The ID of the customer as registered in the CRM system.
	CrmCustomerId string `json:"crmCustomerId,omitempty"`
	// The ID of the customer tenant as registered in the CRM system.
	CrmTenantId string `json:"crmTenantId,omitempty"`
	// (Deprecated) Contains information about the labels assigned to a specified directory. This field supports only single values per key and is now replaced by the string array \"labels\", which supports multiple values per key. The \"customProperties\" field returns only the first value of any label key that has multiple values assigned to it.
	CustomProperties []PropertyResponseObject `json:"customProperties,omitempty"`
	// A description of the global account.
	Description string `json:"description"`
	// The display name of the global account.
	DisplayName string `json:"displayName"`
	// The current state of the global account. * <b>STARTED:</b> CRUD operation on an entity has started. * <b>CREATING:</b> Creating entity operation is in progress. * <b>UPDATING:</b> Updating entity operation is in progress. * <b>MOVING:</b> Moving entity operation is in progress. * <b>PROCESSING:</b> A series of operations related to the entity is in progress. * <b>DELETING:</b> Deleting entity operation is in progress. * <b>OK:</b> The CRUD operation or series of operations completed successfully. * <b>PENDING REVIEW:</b> The processing operation has been stopped for reviewing and can be restarted by the operator. * <b>CANCELLED:</b> The operation or processing was canceled by the operator. * <b>CREATION_FAILED:</b> The creation operation failed, and the entity was not created or was created but cannot be used. * <b>UPDATE_FAILED:</b> The update operation failed, and the entity was not updated. * <b>PROCESSING_FAILED:</b> The processing operations failed. * <b>DELETION_FAILED:</b> The delete operation failed, and the entity was not deleted. * <b>MOVE_FAILED:</b> Entity could not be moved to a different location. * <b>MIGRATING:</b> Migrating entity from NEO to CF.
	EntityState string `json:"entityState,omitempty"`
	// The planned date that the global account expires. This is the same date as the Contract End Date, unless a manual adjustment has been made to the actual expiration date of the global account. Typically, this property is automatically populated only when a formal termination order is received from the CRM system. From a customer perspective, this date marks the start of the grace period, which is typically 30 days before the actual deletion of the account.
	ExpiryDate Time `json:"expiryDate,omitempty"`
	// The geographic locations from where the global account can be accessed. * <b>STANDARD:</b> The global account can be accessed from any geographic location. * <b>EU_ACCESS:</b> The global account can be accessed only within locations in the EU.
	GeoAccess string `json:"geoAccess"`
	// The GUID of the directory's global account entity.
	GlobalAccountGUID string `json:"globalAccountGUID"`
	// The unique ID of the global account.
	Guid string `json:"guid"`
	// Contains information about the labels assigned to a specified global account. Labels are represented in a JSON array of key-value pairs; each key has up to 10 corresponding values. This field replaces the deprecated \"customProperties\" field, which supports only single values per key.
	Labels map[string][]string `json:"labels,omitempty"`
	// The type of license for the global account. The license type affects the scope of functions of the account. * <b>DEVELOPER:</b> For internal developer global accounts on Staging or Canary landscapes. * <b>CUSTOMER:</b> For customer global accounts. * <b>PARTNER:</b> For partner global accounts. * <b>INTERNAL_DEV:</b> For internal global accounts on the Dev landscape. * <b>INTERNAL_PROD:</b> For internal global accounts on the Live landscape. * <b>TRIAL:</b> For customer trial accounts.
	LicenseType string `json:"licenseType"`
	// The date the global account was last modified. Dates and times are in UTC format.
	ModifiedDate Time `json:"modifiedDate,omitempty"`
	// The origin of the account. * <b>ORDER:</b> Created by the Order Processing API or Submit Order wizard. * <b>OPERATOR:</b> Created by the Global Account wizard. * <b>REGION_SETUP:</b> Created automatically as part of the region setup.
	Origin string `json:"origin,omitempty"`
	// The GUID of the global account's parent entity. Typically this is the global account.
	ParentGUID string `json:"parentGUID"`
	// The Type of the global account's parent entity.
	ParentType string `json:"parentType"`
	// The date that an expired contract was renewed. Dates and times are in UTC format.
	RenewalDate Time `json:"renewalDate,omitempty"`
	// For internal accounts, the service for which the global account was created.
	ServiceId string `json:"serviceId,omitempty"`
	// Information about the state.
	StateMessage string `json:"stateMessage,omitempty"`
	// The subaccounts contained in the global account.
	Subaccounts []SubaccountHierarchyResponseObject `json:"subaccounts,omitempty"`
	// Relevant only for entities that require authorization (e.g. global account). The subdomain that becomes part of the path used to access the authorization tenant of the global account. Unique within the defined region.
	Subdomain string `json:"subdomain,omitempty"`
	// Specifies the current stage of the termination notifications sequence. * <b>PENDING_FIRST_NOTIFICATION:</b> A notification has not yet been sent to the global account owner informing them of the expired contract or termination request. * <b>FIRST_NOTIFICATION_PROCESSED:</b> A first notification has been sent to the global account owner informing them of the expired contract, and the termination date when the global account will be closed. * <b>SECOND_NOTIFICATION_PROCESSED:</b> A follow-up notification has been sent to the global account owner.  Your mail server must be configured so that termination notifications can be sent by the Core Commercialization Foundation service.
	TerminationNotificationStatus string `json:"terminationNotificationStatus,omitempty"`
	// For internal accounts, the intended purpose of the global account. Possible purposes: * <b>Development:</b> For development of a service. * <b>Testing:</b> For testing development. * <b>Demo:</b> For creating demos. * <b>Production:</b> For delivering a service in a production landscape.
	UseFor string `json:"useFor,omitempty"`
}
