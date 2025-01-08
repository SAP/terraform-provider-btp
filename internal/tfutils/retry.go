package tfutils

import "strings"

func IsRetriableErrorForEntitlement(err error) bool {

	if err == nil {
		//If no error was raised a retry is not necessary
		return false
	}

	if strings.Contains(err.Error(), "[Error: 30004/400]") {
		// Error code for a locking scenario - API call must be retried
		return true
	} else if strings.Contains(err.Error(), "[Error: 11006/429]") {
		// Error code when hitting a rate limit - API call must be retried
		return true
	} else {
		// No retry possible as error code does not indicate a valid retry scenario
		return false
	}

}

func IsRetriableErrorForEnvInstance(err error) bool {

	if err == nil {
		//If no error was raised a retry is not necessary
		return false
	}

	if strings.Contains(err.Error(), "Command timed out. Please try again later.") {
		// Error 504 Gateway Timeout - API call must be retried
		return true
	} else {
		// No retry possible as error code does not indicate a valid retry scenario
		return false
	}

}
