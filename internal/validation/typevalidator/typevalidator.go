package typevalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type typeValidator struct {
	typeExpr path.Expression
}

func (v typeValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v typeValidator) MarkdownDescription(ctx context.Context) string {
	return "field can only be configured when destination is of type \"HTTP\""
}

func (v typeValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {

	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	// get the path for attribute type from the expression
	typePath, diags := request.Config.PathMatches(ctx, v.typeExpr)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	// get the value of the attribute type from the path
	var typeVal attr.Value
	diags = request.Config.GetAttribute(ctx, typePath[0], &typeVal)
	if diags.HasError() {
		response.Diagnostics.Append(diags...)
		return
	}

	val, ok := typeVal.(types.String)
	if !ok {
		return
	}
	rawVal := val.ValueString() // safely extract the raw string value

	if rawVal != "HTTP" {
		response.Diagnostics.AddAttributeError(
			request.Path,
			v.Description(ctx),
			fmt.Sprintf("Please configure field \"%s\" in the attribute additional_configuration.\nRefer to the examples documented for your specific destination type.", request.Path.String()),
		)
	}
}

func ValidateType(typeExpr path.Expression) validator.String {
	return typeValidator{
		typeExpr: typeExpr,
	}
}
