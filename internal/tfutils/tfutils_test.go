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
			description: "error case - unsupported attribute type",
			uut: struct {
				AListField types.List `tfsdk:"a_list" btpcli:"aList"`
			}{},
			expects: expects{
				errorMessage: "the type 'basetypes.ListValue' assigned to 'aList' is not yet supported",
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
