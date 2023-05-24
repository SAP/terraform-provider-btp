package tfdefaults

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func DefaultStringValue(defaultValue string) planmodifier.String {
	return defaultStringModifier{
		defaultValue: types.StringValue(defaultValue),
	}
}

type defaultStringModifier struct {
	defaultValue        types.String
	description         string
	markdownDescription string
}

// Description returns a human-readable description of the plan modifier.
func (m defaultStringModifier) Description(_ context.Context) string {
	return m.description
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m defaultStringModifier) MarkdownDescription(_ context.Context) string {
	return m.markdownDescription
}

// PlanModifyBool implements the plan modification logic.
func (m defaultStringModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = m.defaultValue
	}
}

func DefaultBoolValue(defaultValue bool) planmodifier.Bool {
	return defaultBoolModifier{
		defaultValue: types.BoolValue(defaultValue),
	}
}

type defaultBoolModifier struct {
	defaultValue        types.Bool
	description         string
	markdownDescription string
}

// Description returns a human-readable description of the plan modifier.
func (m defaultBoolModifier) Description(_ context.Context) string {
	return m.description
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m defaultBoolModifier) MarkdownDescription(_ context.Context) string {
	return m.markdownDescription
}

// PlanModifyBool implements the plan modification logic.
func (m defaultBoolModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.PlanValue.IsNull() {
		resp.PlanValue = m.defaultValue
	}
}
