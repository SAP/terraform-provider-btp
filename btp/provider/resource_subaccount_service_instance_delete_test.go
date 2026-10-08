package provider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func TestServiceInstanceDeletionState(t *testing.T) {
	transportError := errors.New("synthetic transport failure")
	failed := &servicemanager.OperationResponseObject{
		State:  servicemanager.StateFailed,
		Errors: json.RawMessage(`{"broker_error":{"description":"synthetic deletion failure"}}`),
	}
	tests := []struct {
		name        string
		response    servicemanager.ServiceInstanceResponseObject
		status      int
		readError   error
		wantState   string
		wantError   error
		wantMessage string
	}{
		{name: "deadline without metadata", readError: context.DeadlineExceeded, wantError: context.DeadlineExceeded},
		{name: "transport failure without metadata", readError: transportError, wantError: transportError},
		{name: "failed read with operation metadata", response: servicemanager.ServiceInstanceResponseObject{LastOperation: failed}, readError: transportError, wantError: transportError},
		{name: "not found confirms deletion", status: http.StatusNotFound, readError: errors.New("not found"), wantState: "DELETED"},
		{name: "rate limiting remains pending", status: http.StatusTooManyRequests, readError: errors.New("rate limited"), wantState: servicemanager.StateInProgress},
		{name: "missing operation remains pending", status: http.StatusOK, wantState: servicemanager.StateInProgress},
		{name: "deletion in progress", response: servicemanager.ServiceInstanceResponseObject{LastOperation: &servicemanager.OperationResponseObject{State: servicemanager.StateInProgress}}, wantState: servicemanager.StateInProgress},
		{name: "failed deletion keeps broker error", response: servicemanager.ServiceInstanceResponseObject{LastOperation: failed}, wantState: servicemanager.StateFailed, wantMessage: "API error during service instance deletion - synthetic deletion failure"},
		{name: "succeeded operation is not inferred deleted", response: servicemanager.ServiceInstanceResponseObject{LastOperation: &servicemanager.OperationResponseObject{State: servicemanager.StateSucceeded}}, wantState: servicemanager.StateSucceeded},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, state, err := serviceInstanceDeletionState(test.response, btpcli.CommandResponse{StatusCode: test.status}, test.readError)
			if result == nil {
				t.Fatal("refresh must retain a result for explicit polling states")
			}
			if state != test.wantState {
				t.Fatalf("state = %q, want %q", state, test.wantState)
			}
			if test.wantMessage != "" {
				if err == nil || err.Error() != test.wantMessage {
					t.Fatalf("error = %v, want %q", err, test.wantMessage)
				}
			} else if !errors.Is(err, test.wantError) {
				t.Fatalf("error = %v, want %v", err, test.wantError)
			}
		})
	}
}

func TestServiceInstanceDeletionPolling(t *testing.T) {
	t.Run("missing metadata continues until not found", func(t *testing.T) {
		calls := 0
		conf := tfutils.StateChangeConf{
			Pending:      []string{servicemanager.StateInProgress},
			Target:       []string{"DELETED"},
			Timeout:      time.Second,
			PollInterval: time.Millisecond,
			Refresh: func() (any, string, error) {
				calls++
				status := http.StatusOK
				if calls == 3 {
					status = http.StatusNotFound
				}
				return serviceInstanceDeletionState(servicemanager.ServiceInstanceResponseObject{}, btpcli.CommandResponse{StatusCode: status}, nil)
			},
		}
		if _, err := conf.WaitForStateContext(context.Background()); err != nil {
			t.Fatal(err)
		}
		if calls != 3 {
			t.Fatalf("read count = %d, want 3", calls)
		}
	})
	t.Run("read failure stops polling with original error", func(t *testing.T) {
		conf := tfutils.StateChangeConf{
			Pending: []string{servicemanager.StateInProgress},
			Target:  []string{"DELETED"},
			Timeout: time.Second,
			Refresh: func() (any, string, error) {
				return serviceInstanceDeletionState(servicemanager.ServiceInstanceResponseObject{}, btpcli.CommandResponse{}, context.DeadlineExceeded)
			},
		}
		if _, err := conf.WaitForStateContext(context.Background()); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("error = %v, want deadline exceeded", err)
		}
	})
}
