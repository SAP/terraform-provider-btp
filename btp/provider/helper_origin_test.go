package provider

import (
	"testing"
)

func TestOriginMatches(t *testing.T) {
	tests := []struct {
		apiOrigin   string
		stateOrigin string
		want        bool
	}{
		{"sap.default", "sap.default", true},
		{"ldap", "ldap", true},
		{"sap.default", "ldap", true},
		{"ldap", "sap.default", true},
		{"sap.custom", "sap.custom", true},
		{"sap.custom", "ldap", false},
		{"sap.custom", "sap.default", false},
		{"", "", true},
		{"", "ldap", false},
	}

	for _, tt := range tests {
		got := originMatches(tt.apiOrigin, tt.stateOrigin)
		if got != tt.want {
			t.Errorf("originMatches(%q, %q) = %v, want %v", tt.apiOrigin, tt.stateOrigin, got, tt.want)
		}
	}
}

func TestSamlOriginMatches(t *testing.T) {
	tests := []struct {
		name         string
		apiOrigin    string
		samlEntityId string
		stateOrigin  string
		want         bool
	}{
		// Direct idpDisplayName match — same as originMatches
		{
			name:         "exact match via idpDisplayName",
			apiOrigin:    "sap.custom",
			samlEntityId: "",
			stateOrigin:  "sap.custom",
			want:         true,
		},
		{
			name:         "ldap/sap.default alias match",
			apiOrigin:    "sap.default",
			samlEntityId: "",
			stateOrigin:  "ldap",
			want:         true,
		},
		// Custom IAS tenant: idpDisplayName is the hostname, not the origin key
		{
			name:         "custom IAS tenant matched via samlEntityId hostname prefix",
			apiOrigin:    "myidp.accounts400.ondemand.com",
			samlEntityId: "https://myidp.accounts400.ondemand.com",
			stateOrigin:  "myidp",
			want:         true,
		},
		{
			name:         "custom IAS tenant matched via samlEntityId hostname prefix — case insensitive",
			apiOrigin:    "MYIDP.accounts400.ondemand.com",
			samlEntityId: "https://MYIDP.accounts400.ondemand.com",
			stateOrigin:  "myidp",
			want:         true,
		},
		// Reproduces the exact scenario from issue #1670
		{
			name:         "issue 1670 — iasproviderdevblr custom origin matched via samlEntityId",
			apiOrigin:    "iasproviderdevblr.accounts400.ondemand.com",
			samlEntityId: "https://iasproviderdevblr.accounts400.ondemand.com",
			stateOrigin:  "iasproviderdevblr",
			want:         true,
		},
		// No match cases
		{
			name:         "different origin, no samlEntityId",
			apiOrigin:    "other.accounts400.ondemand.com",
			samlEntityId: "",
			stateOrigin:  "myidp",
			want:         false,
		},
		{
			name:         "samlEntityId hostname belongs to a different origin",
			apiOrigin:    "other.accounts400.ondemand.com",
			samlEntityId: "https://other.accounts400.ondemand.com",
			stateOrigin:  "myidp",
			want:         false,
		},
		{
			name:         "origin is prefix of a different origin — no partial match",
			apiOrigin:    "myidpextra.accounts400.ondemand.com",
			samlEntityId: "https://myidpextra.accounts400.ondemand.com",
			stateOrigin:  "myidp",
			want:         false,
		},
		// Edge cases
		{
			name:         "empty samlEntityId falls back gracefully",
			apiOrigin:    "something-else",
			samlEntityId: "",
			stateOrigin:  "myidp",
			want:         false,
		},
		{
			name:         "invalid samlEntityId URL falls back gracefully",
			apiOrigin:    "something-else",
			samlEntityId: "://not-a-url",
			stateOrigin:  "myidp",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := samlOriginMatches(tt.apiOrigin, tt.samlEntityId, tt.stateOrigin)
			if got != tt.want {
				t.Errorf("samlOriginMatches(%q, %q, %q) = %v, want %v",
					tt.apiOrigin, tt.samlEntityId, tt.stateOrigin, got, tt.want)
			}
		})
	}
}
