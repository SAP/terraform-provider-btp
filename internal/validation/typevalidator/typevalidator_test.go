package typevalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var testSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"type":             schema.StringAttribute{},
		"field_under_test": schema.StringAttribute{},
	},
}

func TestTypeValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		destType  string
		expErrors int
	}

	testCases := map[string]testCase{
		"type-http": {
			destType:  "HTTP",
			expErrors: 0,
		},
		"type-non-http": {
			destType:  "RFC",
			expErrors: 1,
		},
	}

	for name, test := range testCases {

		t.Run(name, func(t *testing.T) {

			config := tfsdk.Config{
				Schema: testSchema,
				Raw: tftypes.NewValue(
					tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"type":             tftypes.String,
							"field_under_test": tftypes.String,
						},
					},
					map[string]tftypes.Value{
						"type":             tftypes.NewValue(tftypes.String, test.destType),
						"field_under_test": tftypes.NewValue(tftypes.String, "test"),
					},
				),
			}

			v := typeValidator{
				typeExpr: path.MatchRoot("type"),
			}

			req := validator.StringRequest{
				Config:      config,
				ConfigValue: types.StringValue("test"),
			}

			res := validator.StringResponse{}

			v.ValidateString(
				context.TODO(),
				req,
				&res,
			)

			if !res.Diagnostics.HasError() && test.expErrors > 0 {
				t.Fatalf("expected %d error(s), got none", test.expErrors)
			}

			if res.Diagnostics.HasError() && test.expErrors == 0 {
				t.Fatalf("expected no error(s), got %d: %v", res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}

			if res.Diagnostics.ErrorsCount() != test.expErrors {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}
		})
	}
}
