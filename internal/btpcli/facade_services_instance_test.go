package btpcli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServicesInstanceFacade_List(t *testing.T) {
	command := "services/instance"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"subaccount": subaccountId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Instance.List(context.TODO(), subaccountId, "", "")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
	t.Run("constructs the CLI params correctly - with fieldsFilter", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"subaccount":   subaccountId,
				"fieldsFilter": "ready eq 'true'",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Instance.List(context.TODO(), subaccountId, "ready eq 'true'", "")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
	t.Run("constructs the CLI params correctly - with labelsFilter", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"subaccount":   subaccountId,
				"labelsFilter": "label eq 'value'",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Instance.List(context.TODO(), subaccountId, "", "label eq 'value'")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
	t.Run("constructs the CLI params correctly - with labelsFilter", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"subaccount":   subaccountId,
				"fieldsFilter": "ready eq 'true'",
				"labelsFilter": "label eq 'value'",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Instance.List(context.TODO(), subaccountId, "ready eq 'true'", "label eq 'value'")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestServicesInstanceFacade_GetById(t *testing.T) {
	command := "services/instance"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	instanceId := "bc8a216f-1184-49dc-b4b4-17cfe2828965"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool
		var callCount int

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true
			callCount++

			// Due to fetching the parameters in separate calls, we need to consider the call parameters for the second and third call
			if callCount == 1 {
				assertCall(t, r, command, ActionGet, map[string]string{
					"subaccount": subaccountId,
					"id":         instanceId,
					"parameters": "false",
				})
			} else {
				assertCall(t, r, command, ActionGet, map[string]string{
					"subaccount": subaccountId,
					"id":         instanceId,
					"parameters": "true",
				})
			}
		}))
		defer srv.Close()

		_, res, err := uut.Services.Instance.GetById(context.TODO(), subaccountId, instanceId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})

}

func TestServicesInstanceFacade_GetByName(t *testing.T) {
	command := "services/instance"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	instanceName := "my-instance"

	t.Run("Multiple calls for parameter retrieval retries", func(t *testing.T) {
		var srvCalled bool

		var callCount int

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true
			callCount++

			if callCount == 1 {
				assertCall(t, r, command, ActionGet, map[string]string{
					"subaccount": subaccountId,
					"name":       instanceName,
					"parameters": "false",
				})

			} else {
				assertCall(t, r, command, ActionGet, map[string]string{
					"subaccount": subaccountId,
					"name":       instanceName,
					"parameters": "true",
				})
			}
		}))
		defer srv.Close()

		_, res, err := uut.Services.Instance.GetByName(context.TODO(), subaccountId, instanceName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestServicesInstanceFacade_Create(t *testing.T) {
	command := "services/instance"

	subaccountId := "59cd458e-e66e-4b60-b6d8-8f219379f9a5"
	instanceName := "my-instance"
	servicePlanId := "b50d1b0b-2059-4f21-a014-2ea87752eb48"
	parameters := "{}"
	labels := map[string][]string{
		"a": {"b"},
	}

	t.Run("constructs the CLI params correctly - with parameters set", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"subaccount": subaccountId,
				"name":       instanceName,
				"parameters": parameters,
				"plan":       servicePlanId,
				"labels":     `{"a":["b"]}`,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Instance.Create(context.TODO(), &ServiceInstanceCreateInput{
			Name:          instanceName,
			Subaccount:    subaccountId,
			ServicePlanId: servicePlanId,
			Parameters:    &parameters,
			Labels:        labels,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
	t.Run("constructs the CLI params correctly - no parameters set", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"subaccount": subaccountId,
				"name":       instanceName,
				"plan":       servicePlanId,
				"labels":     `{"a":["b"]}`,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Instance.Create(context.TODO(), &ServiceInstanceCreateInput{
			Name:          instanceName,
			Subaccount:    subaccountId,
			ServicePlanId: servicePlanId,
			Labels:        labels,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestServicesInstanceFacade_CreateVisibility(t *testing.T) {
	const instanceID = "created-instance"
	const subaccountID = "subaccount"
	tests := []struct {
		name            string
		createStatus    int
		createBody      string
		getStatuses     []int
		cancelOnMissing bool
		blockGet        bool
		parentTimeout   time.Duration
		wantStatus      int
		wantError       error
		wantMessage     string
		wantGetTimes    []time.Duration
		wantElapsed     time.Duration
	}{
		{
			name:        "immediately readable",
			getStatuses: []int{200}, wantStatus: 200,
			wantGetTimes: []time.Duration{0},
		},
		{
			name:        "not found twice then readable",
			getStatuses: []int{404, 404, 200}, wantStatus: 200,
			wantGetTimes: []time.Duration{0, time.Second, 3 * time.Second},
			wantElapsed:  3 * time.Second,
		},
		{
			name:        "persistent not found reaches visibility deadline",
			getStatuses: []int{404}, wantStatus: 404,
			wantError: context.DeadlineExceeded, wantMessage: "last GET error: NotFound",
			wantGetTimes: []time.Duration{
				0, time.Second, 3 * time.Second, 7 * time.Second,
				12 * time.Second,
			},
			wantElapsed: 15 * time.Second,
		},
		{
			name:        "cancellation interrupts backoff",
			getStatuses: []int{404}, cancelOnMissing: true, wantStatus: 404,
			wantError: context.Canceled, wantMessage: "last GET error: NotFound",
			wantGetTimes: []time.Duration{0}, wantElapsed: 500 * time.Millisecond,
		},
		{
			name:        "earlier caller deadline interrupts backoff",
			getStatuses: []int{404}, parentTimeout: 500 * time.Millisecond, wantStatus: 404,
			wantError: context.DeadlineExceeded, wantMessage: "last GET error: NotFound",
			wantGetTimes: []time.Duration{0}, wantElapsed: 500 * time.Millisecond,
		},
		{
			name:     "cancellation interrupts GET",
			blockGet: true, cancelOnMissing: true,
			wantError:    context.Canceled,
			wantGetTimes: []time.Duration{0}, wantElapsed: 500 * time.Millisecond,
		},
		{
			name:     "visibility deadline interrupts GET",
			blockGet: true, wantError: context.DeadlineExceeded,
			wantGetTimes: []time.Duration{0}, wantElapsed: 15 * time.Second,
		},
		{
			name:        "other backend errors are returned immediately",
			getStatuses: []int{403}, wantStatus: 403, wantMessage: "Forbidden",
			wantGetTimes: []time.Duration{0},
		},
		{
			name:       "accepted response without ID",
			createBody: "{}", wantStatus: 202, wantMessage: "did not contain an instance ID",
		},
		{
			name:         "synchronous creation does not fetch",
			createStatus: 201, wantStatus: 201,
		},
		{
			name:         "failed creation does not fetch",
			createStatus: 400, createBody: `{"error":"BadRequest"}`,
			wantStatus: 400, wantMessage: "BadRequest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				if tt.parentTimeout != 0 {
					var cancelDeadline context.CancelFunc
					ctx, cancelDeadline = context.WithTimeout(ctx, tt.parentTimeout)
					defer cancelDeadline()
				}
				start := time.Now()
				var createCalls, parameterCalls int
				var getTimes []time.Duration
				transport := serviceInstanceVisibilityTransport(func(r *http.Request) (*http.Response, error) {
					status := http.StatusOK
					body := `{"id":"created-instance"}`
					switch Action(r.URL.RawQuery) {
					case ActionCreate:
						createCalls++
						assertCall(t, r, "services/instance", ActionCreate, map[string]string{
							"name": "instance", "subaccount": subaccountID, "plan": "plan",
						})
						status = tt.createStatus
						if status == 0 {
							status = http.StatusAccepted
						}
						if tt.createBody != "" {
							body = tt.createBody
						}
					case ActionGet:
						var payload struct {
							ParamValues map[string]string
						}
						if !assert.NoError(t, json.NewDecoder(r.Body).Decode(&payload)) {
							t.FailNow()
						}
						assert.Equal(t, instanceID, payload.ParamValues["id"])
						assert.Equal(t, subaccountID, payload.ParamValues["subaccount"])
						if payload.ParamValues["parameters"] == "true" {
							parameterCalls++
							body = `{"data":{"key":"value"}}`
							break
						}
						getTimes = append(getTimes, time.Since(start))
						if tt.cancelOnMissing {
							go func() {
								time.Sleep(500 * time.Millisecond)
								cancel()
							}()
						}
						if tt.blockGet {
							<-r.Context().Done()
							return nil, r.Context().Err()
						}
						if !assert.NotEmpty(t, tt.getStatuses, "unexpected GET") {
							t.FailNow()
						}
						status = tt.getStatuses[min(len(getTimes)-1, len(tt.getStatuses)-1)]
						if status == http.StatusNotFound {
							body = `{"error":"NotFound"}`
						} else if status == http.StatusForbidden {
							body = `{"error":"Forbidden"}`
						}
					default:
						t.Fatalf("unexpected action: %s", r.URL.RawQuery)
					}
					w := httptest.NewRecorder()
					w.Header().Set(HeaderCLIBackendStatus, strconv.Itoa(status))
					w.WriteHeader(http.StatusOK)
					_, err := w.WriteString(body)
					assert.NoError(t, err)
					return w.Result(), nil
				})
				serverURL, err := url.Parse("https://cli.example.test")
				if !assert.NoError(t, err) {
					t.FailNow()
				}
				client := NewV2ClientWithHttpClient(&http.Client{Transport: transport}, serverURL, &RetryConfig{Enabled: false})
				facade := newServicesInstanceFacade(client)
				instance, cr, err := facade.Create(ctx, &ServiceInstanceCreateInput{
					Name: "instance", Subaccount: subaccountID, ServicePlanId: "plan",
				})
				if tt.wantError != nil {
					assert.ErrorIs(t, err, tt.wantError)
					assert.ErrorContains(t, err, instanceID)
					assert.ErrorContains(t, err, "was accepted")
				}
				if tt.wantMessage != "" {
					assert.ErrorContains(t, err, tt.wantMessage)
				}
				if tt.wantError == nil && tt.wantMessage == "" {
					assert.NoError(t, err)
					assert.Equal(t, instanceID, instance.Id)
				}
				assert.Equal(t, tt.wantStatus, cr.StatusCode)
				assert.Equal(t, 1, createCalls)
				assert.Equal(t, tt.wantGetTimes, getTimes)
				assert.Equal(t, tt.wantElapsed, time.Since(start))
				if err == nil && len(getTimes) > 0 {
					assert.Equal(t, 1, parameterCalls)
					assert.JSONEq(t, `{"key":"value"}`, instance.Parameters)
				} else {
					assert.Zero(t, parameterCalls)
				}
			})
		})
	}
}

// serviceInstanceVisibilityTransport keeps requests in memory so synctest can
// advance the retry clock without network I/O or recorded fixtures.
type serviceInstanceVisibilityTransport func(*http.Request) (*http.Response, error)

func (f serviceInstanceVisibilityTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestServicesInstanceFacade_Delete(t *testing.T) {
	command := "services/instance"

	subaccountId := "59cd458e-e66e-4b60-b6d8-8f219379f9a5"
	instanceId := "bc8a216f-1184-49dc-b4b4-17cfe2828965"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"subaccount": subaccountId,
				"id":         instanceId,
				"confirm":    "true",
			})
		}))
		defer srv.Close()

		res, err := uut.Services.Instance.Delete(context.TODO(), subaccountId, instanceId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestServicesInstanceFacade_ComputeLabelDiff(t *testing.T) {

	tests := []struct {
		description string
		labelsOld   map[string][]string
		labelsNew   map[string][]string
		expects     string
	}{

		{
			description: "No changes",
			labelsOld: map[string][]string{
				"foo": {"bar"},
			},
			labelsNew: map[string][]string{
				"foo": {"bar"},
			},
			expects: "",
		},
		{
			description: "Change a label",
			labelsOld: map[string][]string{
				"foo": {"bar"},
			},
			labelsNew: map[string][]string{
				"foo": {"BAR"},
			},
			expects: "[{\"op\":\"remove\",\"key\":\"foo\",\"values\":[\"bar\"]},{\"op\":\"add\",\"key\":\"foo\",\"values\":[\"BAR\"]}]",
		},
		{
			description: "Remove a label",
			labelsOld: map[string][]string{
				"foo": {"bar"},
			},
			labelsNew: map[string][]string{},
			expects:   "[{\"op\":\"remove\",\"key\":\"foo\",\"values\":[\"bar\"]}]",
		},
		{
			description: "Add a label",
			labelsOld: map[string][]string{
				"foo": {"bar"},
			},
			labelsNew: map[string][]string{
				"foo": {"bar"},
				"bar": {"foo"},
			},
			expects: "[{\"op\":\"add\",\"key\":\"bar\",\"values\":[\"foo\"]}]",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := computeLabelParam(test.labelsNew, test.labelsOld)

			assert.Equal(t, test.expects, result)
		})
	}
}
