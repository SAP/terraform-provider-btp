package jsonvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestJSONValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.String
		expErrors int
	}

	testCases := map[string]testCase{
		"simple-match-object": {
			in:        types.StringValue("{}"),
			expErrors: 0,
		},
		"simple-match-array": {
			in:        types.StringValue("[]"),
			expErrors: 0,
		},
		"complex-match-object": {
			in:        types.StringValue("{\"string\": \"value\", \"bool\": true, \"number\": 3.14}"),
			expErrors: 0,
		},
		"simple-mismatch": {
			in:        types.StringValue("foz"),
			expErrors: 1,
		},
		"skip-validation-on-null": {
			in:        types.StringNull(),
			expErrors: 0,
		},
		"skip-validation-on-unknown": {
			in:        types.StringUnknown(),
			expErrors: 0,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: test.in,
			}
			res := validator.StringResponse{}
			ValidJSON().ValidateString(context.TODO(), req, &res)

			if test.expErrors > 0 && !res.Diagnostics.HasError() {
				t.Fatalf("expected %d error(s), got none", test.expErrors)
			}

			if test.expErrors > 0 && test.expErrors != res.Diagnostics.ErrorsCount() {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}

			if test.expErrors == 0 && res.Diagnostics.HasError() {
				t.Fatalf("expected no error(s), got %d: %v", res.Diagnostics.ErrorsCount(), res.Diagnostics)
			}
		})
	}
}
