package tfutils

import (
	"testing"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/stretchr/testify/assert"
)

func TestLabelFilter(t *testing.T) {
	type expects struct {
		output servicemanager.ServiceManagerLabels
	}

	tests := []struct {
		description string
		uut         servicemanager.ServiceManagerLabels
		expects     expects
	}{

		{
			description: "Filter subaccount_id label",
			uut: servicemanager.ServiceManagerLabels{
				"foo":           {"bar"},
				"subaccount_id": {"123"},
			},
			expects: expects{
				output: servicemanager.ServiceManagerLabels{
					"foo": {"bar"},
				},
			},
		},
		{
			description: "Filter no subaccount_id in labels",
			uut: servicemanager.ServiceManagerLabels{
				"foo": {"bar"},
			},
			expects: expects{
				output: servicemanager.ServiceManagerLabels{
					"foo": {"bar"},
				},
			},
		},
		{
			description: "Filter no labels",
			uut: servicemanager.ServiceManagerLabels{
				"foo": {"bar"},
			},
			expects: expects{
				output: servicemanager.ServiceManagerLabels{
					"foo": {"bar"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := RemoveComputedlabels(test.uut)

			assert.Equal(t, test.expects.output, result)
		})
	}
}
