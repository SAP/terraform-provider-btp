package provider

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToGetFormattedValue(t *testing.T) {
	tests := []struct {
		value            string
		descriptionValue string
		description      string
		expects          string
	}{
		{
			value:            "`DefaultValue`",
			descriptionValue: "This is the default value.",
			description:      "happy path - formats the value and the description as a markdown row of a table",
			expects:          "\n  | `DefaultValue` | This is the default value. | ",
		},
		{
			value:            "",
			descriptionValue: "This is the default value.",
			description:      "happy path - empty value, returns an empty cell for the value",
			expects:          "\n  |  | This is the default value. | ",
		},
		{
			value:            "SomeValue",
			descriptionValue: "",
			description:      "happy path - empty description, returns an empty cell for the description",
			expects:          "\n  | SomeValue |  | ",
		},
		{
			value:            "SomeValue|",
			descriptionValue: "",
			description:      "happy path - value contains pipe, the pipe is escaped",
			expects:          "\n  | SomeValue\\| |  | ",
		},
		{
			value:            "",
			descriptionValue: "Description of the value|\"",
			description:      "happy path - description contains pipe, the pipe is escaped",
			expects:          "\n  |  | Description of the value|\" | ",
		},
		{
			value:            "---",
			descriptionValue: "This is the default value.",
			description:      "happy path - value contains markdown header delimiter, returns hyphens as a value",
			expects:          "\n  | --- | This is the default value. | ",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			outputValue := getFormattedValue(test.value, test.descriptionValue)

			assert.Equal(t, test.expects, outputValue)
		})
	}
}
