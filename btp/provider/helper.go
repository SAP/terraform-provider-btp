package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/saas_manager_service"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_trust"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// originMatches reports whether a BTP API origin matches a Terraform state origin.
// BTP uses "sap.default" in API responses for users assigned with the "ldap" origin alias.
func originMatches(apiOrigin, stateOrigin string) bool {
	if apiOrigin == stateOrigin {
		return true
	}
	return (stateOrigin == "ldap" && apiOrigin == "sap.default") ||
		(stateOrigin == "sap.default" && apiOrigin == "ldap")
}

// samlOriginMatches extends originMatches for SAML attribute assignments.
// For custom IAS tenants the idpDisplayName may differ from the origin key, so as a
// fallback we check whether the origin key is the first hostname label of the entity ID
// URL (e.g. origin "myidp" matches samlEntityId "https://myidp.accounts400.ondemand.com").
func samlOriginMatches(apiOrigin, samlEntityId, stateOrigin string) bool {
	if originMatches(apiOrigin, stateOrigin) {
		return true
	}
	u, err := url.Parse(samlEntityId)
	if err != nil {
		return false
	}
	return strings.HasPrefix(strings.ToLower(u.Hostname()), strings.ToLower(stateOrigin)+".")
}

// samlAssignmentOriginMatches reports whether a SAML assignment belongs to an
// origin. In addition to the assignment display name and origin-key prefix,
// compare the assignment entity ID with the entity ID from the matching trust
// configuration. The APIs may return one with a scheme and the other as a
// bare hostname.
func samlAssignmentOriginMatches(apiOrigin, samlEntityID, stateOrigin, trustEntityID string) bool {
	return samlOriginMatches(apiOrigin, samlEntityID, stateOrigin) ||
		(trustEntityID != "" && samlEntityIDsMatch(samlEntityID, trustEntityID))
}

func samlEntityIDsMatch(first, second string) bool {
	firstHost := identityProviderHost(first)
	return firstHost != "" && firstHost == identityProviderHost(second)
}

func identityProviderHost(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	parsedURL, err := url.Parse(value)
	if err == nil && parsedURL.Hostname() != "" {
		return strings.ToLower(strings.TrimSuffix(parsedURL.Hostname(), "."))
	}

	parsedURL, err = url.Parse("https://" + value)
	if err != nil {
		return ""
	}
	return strings.ToLower(strings.TrimSuffix(parsedURL.Hostname(), "."))
}

// trustEntityID returns the identity provider of a trust configuration lookup.
// A missing trust (404) or missing read permission (403) yields an empty value.
// The caller then falls back to the direct origin match. Any other error is
// returned, so that a failed lookup does not drop a SAML assignment from the state.
func trustEntityID(trustConfig xsuaa_trust.TrustConfigurationResponseObject, rawRes btpcli.CommandResponse, err error) (string, error) {
	if err != nil {
		if rawRes.StatusCode == http.StatusNotFound || rawRes.StatusCode == http.StatusForbidden {
			return "", nil
		}
		return "", fmt.Errorf("failed to read the trust configuration of the origin: %w", err)
	}
	return trustConfig.IdentityProvider, nil
}

func stringNullIfEmpty(val string) types.String {
	if len(val) == 0 {
		return types.StringNull()
	}
	return types.StringValue(val)
}

func timeToValue(t time.Time) types.String {
	if t.IsZero() {
		return types.StringNull()
	}

	return types.StringValue(t.Format(time.RFC3339))
}

func handleReadErrors(ctx context.Context, rawRes btpcli.CommandResponse, cliRes any, resp *resource.ReadResponse, err error, resLogName string) {
	// Treat HTTP 404 Not Found status as a signal to recreate resource see https://developer.hashicorp.com/terraform/plugin/framework/resources/read#recommendations
	if rawRes.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	// special case for subscriptions as a 404 is not returned in case of unsubscribed subscriptions
	if strings.Contains(resLogName, "Subscription") {

		if obj, ok := cliRes.(saas_manager_service.EntitledApplicationsResponseObject); ok {
			// if the state of the subscription is "NOT_SUBSCRIBED", the resource is removed from the state to trigger the recreation
			if obj.State == saas_manager_service.StateNotSubscribed {
				resp.State.RemoveResource(ctx)
				return
			}
		} else {
			resp.Diagnostics.AddError(
				"Invalid Response Object",
				"Expected object of type EntitledApplicationsResponseObject for subscriptions",
			)
			return
		}

	}

	resp.Diagnostics.AddError(fmt.Sprintf("API Error Reading %s", resLogName), fmt.Sprintf("%s", err))

}
