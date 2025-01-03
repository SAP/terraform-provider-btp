// this file has been copied from terraform-plugin-testing (v1.2.0) as advised
// https://raw.githubusercontent.com/hashicorp/terraform-plugin-testing/v1.2.0/helper/resource/error.go

package tfutils

import (
	"fmt"
	"strings"
	"time"
)

// NotFoundError represents when a StateRefreshFunc returns a nil result
// during a StateChangeConf waiter method and that StateChangeConf is
// configured for specific targets.
type NotFoundError struct {
	LastError    error
	LastRequest  interface{}
	LastResponse interface{}
	Message      string
	Retries      int
}

// Error returns the Message string, if non-empty, or a string indicating
// the resource could not be found.
func (e *NotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	if e.Retries > 0 {
		return fmt.Sprintf("couldn't find resource (%d retries)", e.Retries)
	}

	return "couldn't find resource"
}

// Unwrap returns the LastError, compatible with errors.Unwrap.
func (e *NotFoundError) Unwrap() error {
	return e.LastError
}

// UnexpectedStateError is returned when Refresh returns a state that's neither in Target nor Pending
type UnexpectedStateError struct {
	LastError     error
	State         string
	ExpectedState []string
}

// Error returns a string with the unexpected state value, the desired target,
// and any last error.
func (e *UnexpectedStateError) Error() string {
	return fmt.Sprintf(
		"unexpected state '%s', wanted target '%s'. last error: %s",
		e.State,
		strings.Join(e.ExpectedState, ", "),
		e.LastError,
	)
}

// Unwrap returns the LastError, compatible with errors.Unwrap.
func (e *UnexpectedStateError) Unwrap() error {
	return e.LastError
}

// TimeoutError is returned when WaitForState times out
type TimeoutError struct {
	LastError     error
	LastState     string
	Timeout       time.Duration
	ExpectedState []string
}

// Error returns a string with any information available.
func (e *TimeoutError) Error() string {
	expectedState := "resource to be gone"
	if len(e.ExpectedState) > 0 {
		expectedState = fmt.Sprintf("state to become '%s'", strings.Join(e.ExpectedState, ", "))
	}

	extraInfo := make([]string, 0)
	if e.LastState != "" {
		extraInfo = append(extraInfo, fmt.Sprintf("last state: '%s'", e.LastState))
	}
	if e.Timeout > 0 {
		extraInfo = append(extraInfo, fmt.Sprintf("timeout: %s", e.Timeout.String()))
	}

	suffix := ""
	if len(extraInfo) > 0 {
		suffix = fmt.Sprintf(" (%s)", strings.Join(extraInfo, ", "))
	}

	if e.LastError != nil {
		return fmt.Sprintf("timeout while waiting for %s%s: %s",
			expectedState, suffix, e.LastError)
	}

	return fmt.Sprintf("timeout while waiting for %s%s\ntry increasing the timeout for this particular operation",
		expectedState, suffix)
}

// Unwrap returns the LastError, compatible with errors.Unwrap.
func (e *TimeoutError) Unwrap() error {
	return e.LastError
}
