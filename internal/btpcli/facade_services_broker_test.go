package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServicesBrokerFacade_List(t *testing.T) {
	command := "services/broker"

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

		_, res, err := uut.Services.Broker.List(context.TODO(), subaccountId, "", "")

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

		_, res, err := uut.Services.Broker.List(context.TODO(), subaccountId, "ready eq 'true'", "")

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

		_, res, err := uut.Services.Broker.List(context.TODO(), subaccountId, "", "label eq 'value'")

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

		_, res, err := uut.Services.Broker.List(context.TODO(), subaccountId, "ready eq 'true'", "label eq 'value'")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestServicesBrokerFacade_GetById(t *testing.T) {
	command := "services/broker"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	brokerId := "9ff44f1b-b2a8-43ae-9072-32bd1dce60e4"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
				"id":         brokerId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Broker.GetById(context.TODO(), subaccountId, brokerId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestServicesBrokerFacade_GetByName(t *testing.T) {
	command := "services/broker"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	brokerName := "my-broker"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
				"name":       brokerName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Broker.GetByName(context.TODO(), subaccountId, brokerName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
