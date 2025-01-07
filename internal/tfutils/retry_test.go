package tfutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetryEnvInstance(t *testing.T) {
	type expects struct {
		retryPossible bool
	}

	tests := []struct {
		description string
		uut         error
		expects     expects
	}{

		{
			description: "Random Error",
			uut:         fmt.Errorf("Random Error"),
			expects: expects{
				retryPossible: false,
			},
		},
		{
			description: "Error 504 Gateway Timeout",
			uut:         fmt.Errorf("Some generic errror text. Command timed out. Please try again later."),
			expects: expects{
				retryPossible: true,
			},
		},
		{
			description: "No Error",
			uut:         nil,
			expects: expects{
				retryPossible: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := IsRetriableErrorForEnvInstance(test.uut)

			assert.Equal(t, test.expects.retryPossible, result)
		})
	}
}

func TestRetrySubaccountEntitlement(t *testing.T) {
	type expects struct {
		retryPossible bool
	}

	tests := []struct {
		description string
		uut         error
		expects     expects
	}{

		{
			description: "Random Error",
			uut:         fmt.Errorf("Random Error"),
			expects: expects{
				retryPossible: false,
			},
		},
		{
			description: "Rate Limit Error",
			uut:         fmt.Errorf("Some generic errror text. [Error: 30004/400]."),
			expects: expects{
				retryPossible: true,
			},
		},
		{
			description: "Locking Error",
			uut:         fmt.Errorf("Some generic errror text. [Error: 11006/429]. Some more generic error text."),
			expects: expects{
				retryPossible: true,
			},
		},

		{
			description: "No Error",
			uut:         nil,
			expects: expects{
				retryPossible: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := IsRetriableErrorForEntitlement(test.uut)

			assert.Equal(t, test.expects.retryPossible, result)
		})
	}
}
