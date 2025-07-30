package labelvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type labelsValidator struct{}

func (v labelsValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v labelsValidator) MarkdownDescription(ctx context.Context) string {
	return "labels must have at most 10 keys and each value must be at most 63 characters."
}

func (v labelsValidator) ValidateMap(ctx context.Context, request validator.MapRequest, response *validator.MapResponse) {
	// Validation for API constraints given by:
	// https://api.sap.com/api/APIAccountsService/path/createSubaccountLabels
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	var labels map[string][]string
	diags := request.ConfigValue.ElementsAs(ctx, &labels, false)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	// Check the number of keys. Must not exceed 10
	if len(labels) > 10 {
		response.Diagnostics.AddAttributeError(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("The labels exceeds the number of 10 keys. Currently has %d keys.", len(labels)),
		)
	}

	// Check the length of each value string, must not exceed 63 characters
	for key, values := range labels {
		for _, value := range values {
			if len(value) > 63 {
				response.Diagnostics.AddAttributeError(
					request.Path,
					v.Description(ctx),
					fmt.Sprintf("The value for key '%s' exceeds the maximum length of 63 characters.", key),
				)
			}
		}
	}
}

func ValidLabels() validator.Map {
	return labelsValidator{}
}
