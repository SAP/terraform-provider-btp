package provider

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/saas_manager_service"
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
