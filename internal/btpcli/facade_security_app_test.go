package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityAppFacade_ListByGlobalAccount(t *testing.T) {
	command := "security/app"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.App.ListByGlobalAccount(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityAppFacade_ListBySubaccount(t *testing.T) {
	command := "security/app"

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

		_, res, err := uut.Security.App.ListBySubaccount(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityAppFacade_ListByDirectory(t *testing.T) {
	command := "security/app"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"directory": directoryId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.App.ListByDirectory(context.TODO(), directoryId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityAppFacade_GetByGlobalAccount(t *testing.T) {
	command := "security/app"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	appId := "Global Account Administrator"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"globalAccount": globalAccountId,
				"appId":         appId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.App.GetByGlobalAccount(context.TODO(), appId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityAppFacade_GetBySubaccount(t *testing.T) {
	command := "security/app"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	appId := "Subaccount Administrator"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
				"appId":      appId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.App.GetBySubaccount(context.TODO(), subaccountId, appId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityAppFacade_GetByDirectory(t *testing.T) {
	command := "security/app"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	appId := "Directory Administrator"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"directory": directoryId,
				"appId":     appId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.App.GetByDirectory(context.TODO(), directoryId, appId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
