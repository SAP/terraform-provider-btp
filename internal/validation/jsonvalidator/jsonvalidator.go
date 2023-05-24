package jsonvalidator

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type jsonValidator struct {
}

func (v jsonValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v jsonValidator) MarkdownDescription(_ context.Context) string {
	return "value must be valid json"
}

func (v jsonValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue

	if json.Valid([]byte(value.ValueString())) {
		return
	}

	response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
		request.Path,
		v.Description(ctx),
		value.String(),
	))
}

// ValidJSON checks that the String held in the attribute
// is a valid JSON string
func ValidJSON() validator.String {
	return jsonValidator{}
}
