package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func handleReadErrors(ctx context.Context, rawRes btpcli.CommandResponse, resp *resource.ReadResponse, err error, resLogName string) {
	// Treat HTTP 404 Not Found status as a signal to recreate resource see https://developer.hashicorp.com/terraform/plugin/framework/resources/read#recommendations
	if rawRes.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
	} else {
		resp.Diagnostics.AddError(fmt.Sprintf("API Error Reading %s", resLogName), fmt.Sprintf("%s", err))
	}

}
