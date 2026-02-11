package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityIdentityProviderFacade_ListByGlobalAccount(t *testing.T) {
	command := "security/available-idp"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Idp.ListByGlobalAccount(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityIdentityProviderFacade_GetByGlobalAccount(t *testing.T) {
	command := "security/available-idp"
	host := "a2dynbhnd.accounts400.ondemand.com"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
				"iasTenantUrl":  host,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Idp.GetByGlobalAccount(context.TODO(), host)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityIdentityProviderFacade_ListBySubaccount(t *testing.T) {
	command := "security/available-idp"
	subaccountId := "77395f6a-a601-4c9e-8cd0-c1fcefc7f60f"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"subaccount": subaccountId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Idp.ListBySubaccount(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityIdentityProviderFacade_GetBySubaccount(t *testing.T) {
	command := "security/available-idp"
	subaccountId := "77395f6a-a601-4c9e-8cd0-c1fcefc7f60f"
	host := "a2dynbhnd.accounts400.ondemand.com"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount":   subaccountId,
				"iasTenantUrl": host,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Idp.GetBySubaccount(context.TODO(), subaccountId, host)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
