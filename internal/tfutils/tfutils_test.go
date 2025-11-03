package tfutils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestToBTPCLIParamsMap(t *testing.T) {
	type expects struct {
		output       map[string]string
		errorMessage string
	}

	expectsNOP := expects{
		output: map[string]string{},
	}

	tests := []struct {
		description string
		uut         any
		tag         string
		expects     expects
	}{
		{
			description: "NOP",
			uut:         struct{}{},
			expects:     expectsNOP,
		},
		{
			description: "NOP - no btpcli annotations",
			uut: struct {
				AStringField types.String `tfsdk:"a_string_field"`
			}{
				AStringField: types.StringValue("a string value"),
			},
			expects: expectsNOP,
		},
		{
			description: "NOP - everything correctly setup, but field value is unknown",
			uut: struct {
				AStringField types.String `tfsdk:"a_string_field" btpcli:"aStringField"`
			}{
				AStringField: types.StringUnknown(),
			},
			expects: expectsNOP,
		},
		{
			description: "happy path - simple case",
			uut: struct {
				AStringField types.String `tfsdk:"a_string_field" btpcli:"aStringField"`
			}{
				AStringField: types.StringValue("some value set"),
			},
			expects: expects{
				output: map[string]string{
					"aStringField": "some value set",
				},
			},
		},
		{
			description: "happy path - different types",
			uut: struct {
				AStringField       types.String `tfsdk:"a_string_field" btpcli:"aStringField"`
				AnotherStringField types.String `tfsdk:"another_string_field" btpcli:"anotherStringField"`
				ABoolean           types.Bool   `tfsdk:"a_bool" btpcli:"aBoolField"`
			}{
				AStringField:       types.StringValue("a value"),
				AnotherStringField: types.StringValue("another value"),
				ABoolean:           types.BoolValue(true),
			},
			expects: expects{
				output: map[string]string{
					"aStringField":       "a value",
					"anotherStringField": "another value",
					"aBoolField":         "true",
				},
			},
		},
		{
			description: "happy path - unknown and null values get skipped",
			uut: struct {
				Id         types.String `tfsdk:"id" btpcli:"id"`
				BindingId  types.String `tfsdk:"binding_id" btpcli:"bindingID"`
				Name       types.String `tfsdk:"name" btpcli:"name"`
				Subaccount types.String `tfsdk:"subaccount" btpcli:"subaccount"`
			}{
				Id:         types.StringValue("c2d02852-1678-4c1e-b546-74d5274f1522"),
				Subaccount: types.StringValue("6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"),
			},
			expects: expects{
				output: map[string]string{
					"id":         "c2d02852-1678-4c1e-b546-74d5274f1522",
					"subaccount": "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f",
				},
			},
		},
		{
			description: "happy path - string slices",
			uut: struct {
				Features []string `tfsdk:"features" btpcli:"directoryFeatures"`
			}{
				Features: []string{"DEFAULT", "AUTHORIZATIONS", "ENTITLEMENTS"},
			},
			expects: expects{
				output: map[string]string{
					"directoryFeatures": "DEFAULT,AUTHORIZATIONS,ENTITLEMENTS",
				},
			},
		},
		{
			description: "happy path - slice as json",
			uut: struct {
				Features []string `btpcli:"directoryFeatures,json"`
			}{
				Features: []string{"DEFAULT", "AUTHORIZATIONS", "ENTITLEMENTS"},
			},
			expects: expects{
				output: map[string]string{
					"directoryFeatures": "[\"DEFAULT\",\"AUTHORIZATIONS\",\"ENTITLEMENTS\"]",
				},
			},
		},
		{
			description: "happy path - map as json",
			uut: struct {
				Labels map[string]string `btpcli:"labels,json"`
			}{
				Labels: map[string]string{
					"a": "b",
				},
			},
			expects: expects{
				output: map[string]string{
					"labels": "{\"a\":\"b\"}",
				},
			},
		},
		{
			description: "error case - unsupported attribute type",
			uut: struct {
				AListField types.List `tfsdk:"a_list" btpcli:"aList"`
			}{},
			expects: expects{
				errorMessage: "unable to encode 'aList': unsupported type 'basetypes.ListValue'",
			},
		},

		// TODO check that strings get properly escaped
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			output, err := ToBTPCLIParamsMap(&test.uut)

			if len(test.expects.errorMessage) > 0 {
				assert.EqualError(t, err, test.expects.errorMessage)
				assert.Empty(t, output)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expects.output, output)
			}
		})
	}
}

func TestNormalizeJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Valid JSON with reordered keys",
			input:       `{"id":1,"email":"test@sap.com"}`,
			expected:    `{"email":"test@sap.com","id":1}`,
			expectError: false,
		},
		{
			name:        "Valid JSON with nested structure",
			input:       `{"instance_name":"test","cf_users":[{"email":"test@sap.com","id":3}]}`,
			expected:    `{"cf_users":[{"email":"test@sap.com","id":3}],"instance_name":"test"}`,
			expectError: false,
		},
		{
			name:        "Empty JSON object",
			input:       `{}`,
			expected:    `{}`,
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			input:       `{"a":1,}`,
			expected:    "",
			expectError: true,
		},
		{
			name:        "Valid JSON array",
			input:       `[{"b":2,"a":1},{"d":4,"c":3}]`,
			expected:    `[{"a":1,"b":2},{"c":3,"d":4}]`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeJSON(tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, tt.expected, result)
			}
		})
	}
}
