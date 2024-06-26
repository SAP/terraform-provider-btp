/*
 * Authorization
 *
 * Provides functions to administrate the Authorization and Trust Management service (XSUAA) of SAP BTP, Cloud Foundry environment. You can manage service instances of the Authorization and Trust Management service. You can also manage roles, role templates, and role collections of your subaccount.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package xsuaa_authz

type RoleAttribute struct {
	// The name has a maximum length of 64 characters. Only the following characters are allowed: alphanumeric characters (aA-zZ) and (0-9) and underscore (_).
	AttributeName        string   `json:"attributeName,omitempty"`
	AttributeValueOrigin string   `json:"attributeValueOrigin,omitempty"`
	AttributeValues      []string `json:"attributeValues,omitempty"`
	Description          string   `json:"description,omitempty"`
	ValueRequired        bool     `json:"valueRequired,omitempty"`
}
