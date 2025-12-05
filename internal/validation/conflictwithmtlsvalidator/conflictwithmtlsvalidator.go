package conflictwithmtlsvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// custom validator to ensure cert and key are not provided when mtls is set
type ConflictsWithMTLSValidator struct{}

func (v ConflictsWithMTLSValidator) Description(_ context.Context) string {
	return "Cannot be provided when mtls is true."
}

func (v ConflictsWithMTLSValidator) MarkdownDescription(_ context.Context) string {
	return "Cannot be provided when `mtls` is `true`."
}

func (v ConflictsWithMTLSValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Check mtls attribute
	var mtlsVal types.Bool
	diag := req.Config.GetAttribute(ctx, path.Root("mtls"), &mtlsVal)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	if mtlsVal.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Combination",
			"When `mtls` is true, `cert` and `key` must NOT be provided.",
		)
	}
}
