package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsEntitlementFacade_ListByGlobalAccount(t *testing.T) {
	command := "accounts/entitlement"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Entitlement.ListByGlobalAccount(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEntitlementFacade_ListBySubaccount(t *testing.T) {
	command := "accounts/entitlement"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"subaccountFilter": "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Entitlement.ListBySubaccount(context.TODO(), "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEntitlementFacade_ListByDirectory(t *testing.T) {
	command := "accounts/entitlement"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"directory": "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Entitlement.ListByDirectory(context.TODO(), "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEntitlementFacade_AssignToSubaccount(t *testing.T) {
	command := "accounts/entitlement"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	serviceName := "alert-notification"
	planName := "free"
	amount := 10

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"subaccount":      subaccountId,
				"serviceName":     serviceName,
				"servicePlanName": planName,
				"amount":          "10",
			})
		}))
		defer srv.Close()

		res, err := uut.Accounts.Entitlement.AssignToSubaccount(context.TODO(), subaccountId, serviceName, planName, amount)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEntitlementFacade_EnableInSubaccount(t *testing.T) {
	command := "accounts/entitlement"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	serviceName := "alert-notification"
	planName := "free"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"subaccount":      subaccountId,
				"serviceName":     serviceName,
				"servicePlanName": planName,
				"enable":          "true",
			})
		}))
		defer srv.Close()

		res, err := uut.Accounts.Entitlement.EnableInSubaccount(context.TODO(), subaccountId, serviceName, planName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEntitlementFacade_DisableInSubaccount(t *testing.T) {
	command := "accounts/entitlement"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	serviceName := "alert-notification"
	planName := "free"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"subaccount":      subaccountId,
				"serviceName":     serviceName,
				"servicePlanName": planName,
				"enable":          "false",
			})
		}))
		defer srv.Close()

		res, err := uut.Accounts.Entitlement.DisableInSubaccount(context.TODO(), subaccountId, serviceName, planName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEntitlementFacade_AssignToDirectory(t *testing.T) {
	command := "accounts/entitlement"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	serviceName := "alert-notification"
	planName := "free"
	amount := 10

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"directory":       directoryId,
				"serviceName":     serviceName,
				"servicePlanName": planName,
				"amount":          "10",
			})
		}))
		defer srv.Close()

		res, err := uut.Accounts.Entitlement.AssignToDirectory(context.TODO(), directoryId, serviceName, planName, amount)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEntitlementFacade_EnableInDirectory(t *testing.T) {
	command := "accounts/entitlement"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	serviceName := "alert-notification"
	planName := "free"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"directory":       directoryId,
				"serviceName":     serviceName,
				"servicePlanName": planName,
				"enable":          "true",
			})
		}))
		defer srv.Close()

		res, err := uut.Accounts.Entitlement.EnableInDirectory(context.TODO(), directoryId, serviceName, planName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEntitlementFacade_DisableInDirectory(t *testing.T) {
	command := "accounts/entitlement"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	serviceName := "alert-notification"
	planName := "free"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"directory":       directoryId,
				"serviceName":     serviceName,
				"servicePlanName": planName,
				"enable":          "false",
			})
		}))
		defer srv.Close()

		res, err := uut.Accounts.Entitlement.DisableInDirectory(context.TODO(), directoryId, serviceName, planName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
