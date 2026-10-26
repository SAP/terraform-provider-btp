package provider

import (
	"errors"
	"testing"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_trust"
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

func TestSamlAssignmentOriginMatches(t *testing.T) {
	tests := []struct {
		name          string
		apiOrigin     string
		samlEntityID  string
		stateOrigin   string
		trustEntityID string
		want          bool
	}{
		{
			name:          "trust entity hostname matches despite scheme",
			apiOrigin:     "iasprovidertestblr.accounts400.ondemand.com (platform users)",
			samlEntityID:  "https://iasprovidertestblr.accounts400.ondemand.com",
			stateOrigin:   "iasprovidertestblr-platform",
			trustEntityID: "iasprovidertestblr.accounts400.ondemand.com",
			want:          true,
		},
		{
			name:          "sap.custom matches trust entity hostname without identity provider",
			samlEntityID:  "https://myidp.accounts.ondemand.com",
			stateOrigin:   "sap.custom",
			trustEntityID: "myidp.accounts.ondemand.com",
			want:          true,
		},
		{
			name:          "different trust entity hostname does not match",
			samlEntityID:  "https://iasprovidertestblr.accounts400.ondemand.com",
			stateOrigin:   "iasprovidertestblr-platform",
			trustEntityID: "other.accounts400.ondemand.com",
			want:          false,
		},
		{
			name:          "empty trust entity ID does not match",
			samlEntityID:  "https://myidp.accounts.ondemand.com",
			stateOrigin:   "sap.custom",
			trustEntityID: "",
			want:          false,
		},
		{
			name:          "empty SAML entity ID does not match",
			samlEntityID:  "",
			stateOrigin:   "sap.custom",
			trustEntityID: "myidp.accounts.ondemand.com",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := samlAssignmentOriginMatches(tt.apiOrigin, tt.samlEntityID, tt.stateOrigin, tt.trustEntityID); got != tt.want {
				t.Errorf("samlAssignmentOriginMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrustEntityID(t *testing.T) {
	trustConfig := xsuaa_trust.TrustConfigurationResponseObject{IdentityProvider: "myidp.accounts.ondemand.com"}
	lookupErr := errors.New("lookup failed")

	tests := []struct {
		name        string
		trustConfig xsuaa_trust.TrustConfigurationResponseObject
		rawRes      btpcli.CommandResponse
		err         error
		want        string
		wantErr     bool
	}{
		{
			name:        "trust found",
			trustConfig: trustConfig,
			rawRes:      btpcli.CommandResponse{StatusCode: 200},
			want:        "myidp.accounts.ondemand.com",
		},
		{
			name:   "missing trust is not an error",
			rawRes: btpcli.CommandResponse{StatusCode: 404},
			err:    lookupErr,
		},
		{
			name:   "missing read permission is not an error",
			rawRes: btpcli.CommandResponse{StatusCode: 403},
			err:    lookupErr,
		},
		{
			name:    "server error is returned",
			rawRes:  btpcli.CommandResponse{StatusCode: 500},
			err:     lookupErr,
			wantErr: true,
		},
		{
			name:    "error without response is returned",
			err:     lookupErr,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := trustEntityID(tt.trustConfig, tt.rawRes, tt.err)
			if (err != nil) != tt.wantErr {
				t.Fatalf("trustEntityID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("trustEntityID() = %q, want %q", got, tt.want)
			}
		})
	}
}
