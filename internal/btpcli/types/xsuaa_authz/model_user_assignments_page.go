package xsuaa_authz

type UserAssignmentsPage struct {
	Count      int             `json:"count,omitempty"`
	TotalPages int             `json:"totalPages,omitempty"`
	Items      []UserReference `json:"items,omitempty"`
}
