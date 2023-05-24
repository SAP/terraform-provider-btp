package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServicesPlatformFacade_List(t *testing.T) {
	command := "services/platform"

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

		_, res, err := uut.Services.Platform.List(context.TODO(), subaccountId, "", "")

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

		_, res, err := uut.Services.Platform.List(context.TODO(), subaccountId, "ready eq 'true'", "")

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

		_, res, err := uut.Services.Platform.List(context.TODO(), subaccountId, "", "label eq 'value'")

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

		_, res, err := uut.Services.Platform.List(context.TODO(), subaccountId, "ready eq 'true'", "label eq 'value'")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestServicesPlatformFacade_GetById(t *testing.T) {
	command := "services/platform"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	platformId := "76765dca-6683-473a-8f42-809e33a2ea68"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
				"id":         platformId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Platform.GetById(context.TODO(), subaccountId, platformId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestServicesPlatformFacade_GetByName(t *testing.T) {
	command := "services/platform"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	platformName := "my-platform"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
				"name":       platformName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Services.Platform.GetByName(context.TODO(), subaccountId, platformName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
