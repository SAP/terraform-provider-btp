package tfdefaults

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/stretchr/testify/assert"
)

func TestDefaulStringtValue(t *testing.T) {
	tests := []struct {
		description   string
		plannedValue  types.String
		currentValue  types.String
		defaultValue  string
		expectedValue attr.Value
	}{
		{
			description:   "value is set, nothing to do",
			plannedValue:  types.StringValue("gamma"),
			currentValue:  types.StringValue("beta"),
			defaultValue:  "alpha",
			expectedValue: types.StringValue("gamma"),
		},
		{
			description:   "value is planned, nothing to do",
			plannedValue:  types.StringValue("gamma"),
			currentValue:  types.StringNull(),
			defaultValue:  "alpha",
			expectedValue: types.StringValue("gamma"),
		},
		{
			description:   "default string",
			plannedValue:  types.StringUnknown(),
			currentValue:  types.StringUnknown(),
			defaultValue:  "alpha",
			expectedValue: types.StringValue("alpha"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctx := context.Background()
			request := planmodifier.StringRequest{
				Path:       path.Root("test"),
				PlanValue:  test.plannedValue,
				StateValue: test.currentValue,
			}
			response := planmodifier.StringResponse{
				PlanValue: request.PlanValue,
			}
			DefaultStringValue(test.defaultValue).PlanModifyString(ctx, request, &response)

			assert.False(t, response.Diagnostics.HasError())
			assert.Equal(t, test.expectedValue.String(), response.PlanValue.String())
		})
	}
}

func TestDefaultBoolValue(t *testing.T) {
	tests := []struct {
		description   string
		plannedValue  types.Bool
		currentValue  types.Bool
		defaultValue  bool
		expectedValue attr.Value
	}{
		/*{
			description:   "default string",
			plannedValue:  types.BoolNull(),
			currentValue:  types.BoolNull(),
			defaultValue:  types.BoolValue("alpha"),
			expectedValue: types.BoolValue("alpha"),
		},*/ // TODO
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			ctx := context.Background()
			request := planmodifier.BoolRequest{
				Path:       path.Root("test"),
				PlanValue:  test.plannedValue,
				StateValue: test.currentValue,
			}
			response := planmodifier.BoolResponse{
				PlanValue: request.PlanValue,
			}
			DefaultBoolValue(test.defaultValue).PlanModifyBool(ctx, request, &response)

			assert.False(t, response.Diagnostics.HasError())
			assert.Equal(t, test.expectedValue.String(), response.PlanValue.String())
		})
	}
}
