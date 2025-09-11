package labelvalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

	labels := request.ConfigValue.Elements()

	/*var labels map[string][]string
	diags := request.ConfigValue.ElementsAs(ctx, &labels, false)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}
	*/
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

		if values.IsUnknown() {
			continue
		}

		switch listValue := values.(type) {
		case basetypes.SetValue:
			if listValue.IsUnknown() {
				continue
			}
			setElements := listValue.Elements()
			for _, setElement := range setElements {
				if stringValue, ok := setElement.(basetypes.StringValue); ok {
					if stringValue.IsUnknown() {
						continue
					}
					if len(stringValue.ValueString()) > 63 {
						response.Diagnostics.AddAttributeError(
							request.Path,
							v.Description(ctx),
							fmt.Sprintf("The value for key '%s' exceeds the maximum length of 63 characters.", key),
						)
					}
				}
			}
		default:
			// Handle other types or add error
			response.Diagnostics.AddAttributeError(
				request.Path,
				v.Description(ctx),
				fmt.Sprintf("Unexpected element type for key '%s'", key),
			)
		}
	}
}

func ValidLabels() validator.Map {
	return labelsValidator{}
}
