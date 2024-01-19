package btpcli

import (
	"context"
	"net/http"
	"testing"

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

	t.Run("looks up instance with parameters", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
				"id":         instanceId,
				"parameters": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Instance.GetById(context.TODO(), subaccountId, instanceId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
	t.Run("if parameters can't be resolved, it falls back to resolution without parameters", func(t *testing.T) {
		var srvCalled bool

		var callCount int

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true
			callCount++

			if callCount == 1 {
				assertCall(t, r, command, ActionGet, map[string]string{
					"subaccount": subaccountId,
					"id":         instanceId,
					"parameters": "true",
				})

				w.WriteHeader(http.StatusNotFound)
			} else {
				assertCall(t, r, command, ActionGet, map[string]string{
					"subaccount": subaccountId,
					"id":         instanceId,
					"parameters": "false",
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

	t.Run("looks up instance with parameters", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
				"name":       instanceName,
				"parameters": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Instance.GetByName(context.TODO(), subaccountId, instanceName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
	t.Run("if parameters can't be resolved, it falls back to resolution without parameters", func(t *testing.T) {
		var srvCalled bool

		var callCount int

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true
			callCount++

			if callCount == 1 {
				assertCall(t, r, command, ActionGet, map[string]string{
					"subaccount": subaccountId,
					"name":       instanceName,
					"parameters": "true",
				})

				w.WriteHeader(http.StatusNotFound)
			} else {
				assertCall(t, r, command, ActionGet, map[string]string{
					"subaccount": subaccountId,
					"name":       instanceName,
					"parameters": "false",
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

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
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

func TestServicesInstanceFacade_Delete(t *testing.T) {
	command := "services/instance"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
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
