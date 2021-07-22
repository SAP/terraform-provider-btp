package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
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
