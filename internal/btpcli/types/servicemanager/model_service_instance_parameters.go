/*
 * Service Manager
 *
 * Manually created types for the handling of parameters of service instances.
 *
 */

package servicemanager

type ServiceInstanceParametersPlain struct {
	Parameters map[string]interface{} `json:"-"`
}

type ServiceInstanceParametersData struct {
	Parameters map[string]interface{} `json:"data,omitempty"`
}
