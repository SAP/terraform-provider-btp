package xsuaa_settings

// LinksSettings LinksSettings
//
// swagger:model LinksSettings
type LinksSettings struct {

	// Overrides the home page of the service and issues a redirect to this URL when the browser requests `/` or `/home`.
	HomeRedirect string `json:"homeRedirect,omitempty"`
}
