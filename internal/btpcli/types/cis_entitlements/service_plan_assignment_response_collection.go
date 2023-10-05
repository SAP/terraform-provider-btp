/*
 * Entitlements Service
 *
 * Manual adoption of service endpoint leveraged by BTP CLI
 *
 */

package cis_entitlements

type ServicePlanAssignmentsResponseCollection struct {
	AssignedService []servicePlanAssignmentsResponseObject `json:"quotas,omitempty"`
}

type servicePlanAssignmentsResponseObject struct {
	ConsumedQuota      int64  `json:"consumedQuota,omitempty"`
	GlobalAccountGuid  string `json:"globalAccountGUID,omitempty"`
	globalAccountId    string `json:"globalAccountId,omitempty"`
	Plan               string `json:"plan,omitempty"`
	ProvisioningMethod string `json:"provisioningMethod,omitempty"`
	Quota              int32  `json:"quota,omitempty"`
	Service            string `json:"service,omitempty"`
	ServiceCattegaory  string `json:"serviceCategory,omitempty"`
	SubaccountGuid     string `json:"subaccountGUID,omitempty"`
	TenantId           string `json:"tenantId,omitempty"`
	UniqueIdentifier   string `json:"uniqueIdentifier,omitempty"`
	Unlimited          bool   `json:"unlimited,omitempty"`
	// resources	[...]

}
