package uuidvalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUUIDValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.String
		expErrors int
	}

	testCases := map[string]testCase{
		"simple-match-lowercase": {
			in:        types.StringValue("dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"),
			expErrors: 0,
		},
		"simple-match-uppercase": {
			in:        types.StringValue("6AA64C2F-38C1-49A9-B2E8-CF9FEA769B7F"),
			expErrors: 0,
		},
		"simple-mismatch": {
			in:        types.StringValue("sth-which-is-not-a-uuid"),
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
			ValidUUID().ValidateString(context.TODO(), req, &res)

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
