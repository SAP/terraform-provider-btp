/*
 * Service Manager
 *
 * Manually created types for the handling of parameters of service instances.
 *
 */

package servicemanager

import "encoding/json"

// ServiceInstanceParametersPlain models the parameter response of a service
// instance that returns its parameters as a flat, top-level JSON object, e.g.
//
//	{"backend":{"api_enabled":true},"ingest_otlp":{"enabled":true}}
//
// The whole body IS the parameter map, so a struct tag cannot address it — the
// map is populated via a custom UnmarshalJSON. A plain `json:"-"` tag here left
// the field permanently empty, silently dropping instance parameters (#278).
type ServiceInstanceParametersPlain struct {
	Parameters map[string]any `json:"-"`
}

func (p *ServiceInstanceParametersPlain) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Parameters)
}

// ServiceInstanceParametersData models the response of instances that wrap
// their parameters in a top-level "data" object.
type ServiceInstanceParametersData struct {
	Parameters map[string]any `json:"data,omitempty"`
}
