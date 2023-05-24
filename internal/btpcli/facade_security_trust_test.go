package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityTrustFacade_ListByGlobalAccount(t *testing.T) {
	command := "security/trust"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Trust.ListByGlobalAccount(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityTrustFacade_ListBySubaccount(t *testing.T) {
	command := "security/trust"

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

		_, res, err := uut.Security.Trust.ListBySubaccount(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityTrustFacade_GetByGlobalAccount(t *testing.T) {
	command := "security/trust"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"globalAccount": globalAccountId,
				"origin":        origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Trust.GetByGlobalAccount(context.TODO(), origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityTrustFacade_GetBySubaccount(t *testing.T) {
	command := "security/trust"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
				"origin":     origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Trust.GetBySubaccount(context.TODO(), subaccountId, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityTrustFacade_CreateByGlobalAccount(t *testing.T) {
	command := "security/trust"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	idp := "my-ias-tentant.local"
	name := "my-ias"
	description := "this is a description for the ias tenant"
	origin := "custom-origin-platform"

	t.Run("constructs the CLI params correctly - minimal", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"globalAccount": globalAccountId,
				"iasTenantUrl":  idp,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Trust.CreateByGlobalAccount(context.TODO(), TrustConfigurationInput{
			IdentityProvider: idp,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
	t.Run("constructs the CLI params correctly - fully customized", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"globalAccount": globalAccountId,
				"iasTenantUrl":  idp,
				"name":          name,
				"description":   description,
				"origin":        origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Trust.CreateByGlobalAccount(context.TODO(), TrustConfigurationInput{
			IdentityProvider: idp,
			Name:             &name,
			Description:      &description,
			Origin:           &origin,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityTrustFacade_CreateBySubaccount(t *testing.T) {
	command := "security/trust"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	idp := "my-ias-tenant.local"
	name := "my-ias"
	description := "this is a description for the ias tenant"
	origin := "custom-origin-platform"

	t.Run("constructs the CLI params correctly - minimal", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"subaccount":   subaccountId,
				"iasTenantUrl": idp,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Trust.CreateBySubaccount(context.TODO(), subaccountId, TrustConfigurationInput{
			IdentityProvider: idp,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
	t.Run("constructs the CLI params correctly - fully customized", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"subaccount":   subaccountId,
				"iasTenantUrl": idp,
				"name":         name,
				"description":  description,
				"origin":       origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Trust.CreateBySubaccount(context.TODO(), subaccountId, TrustConfigurationInput{
			IdentityProvider: idp,
			Name:             &name,
			Description:      &description,
			Origin:           &origin,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityTrustFacade_DeleteByGlobalAccount(t *testing.T) {
	command := "security/trust"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	originKey := "my-idp-platform"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"globalAccount": globalAccountId,
				"originKey":     originKey,
				"confirm":       "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Trust.DeleteByGlobalAccount(context.TODO(), originKey)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityTrustFacade_DeleteBySubaccount(t *testing.T) {
	command := "security/trust"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	originKey := "my-idp-platform"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"subaccount": subaccountId,
				"originKey":  originKey,
				"confirm":    "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Trust.DeleteBySubaccount(context.TODO(), subaccountId, originKey)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
