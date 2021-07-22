package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityUserFacade_ListByGlobalAccount(t *testing.T) {
	command := "security/user"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": globalAccountId,
				"origin":        origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.User.ListByGlobalAccount(context.TODO(), origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityUserFacade_ListBySubaccount(t *testing.T) {
	command := "security/user"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"subaccount": subaccountId,
				"origin":     origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.User.ListBySubaccount(context.TODO(), subaccountId, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityUserFacade_ListByDirectory(t *testing.T) {
	command := "security/user"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"directory": directoryId,
				"origin":    origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.User.ListByDirectory(context.TODO(), directoryId, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityUserFacade_GetByGlobalAccount(t *testing.T) {
	command := "security/user"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	userName := "john.doe@mycompany.com"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"globalAccount": globalAccountId,
				"userName":      userName,
				"origin":        origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.User.GetByGlobalAccount(context.TODO(), userName, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityUserFacade_GetBySubaccount(t *testing.T) {
	command := "security/user"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	userName := "john.doe@mycompany.com"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
				"userName":   userName,
				"origin":     origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.User.GetBySubaccount(context.TODO(), subaccountId, userName, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityUserFacade_GetByDirectory(t *testing.T) {
	command := "security/user"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	userName := "john.doe@mycompany.com"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"directory": directoryId,
				"userName":  userName,
				"origin":    origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.User.GetByDirectory(context.TODO(), directoryId, userName, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
