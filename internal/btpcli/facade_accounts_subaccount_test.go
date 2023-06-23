package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsSubaccountFacade_List(t *testing.T) {
	command := "accounts/subaccount"

	t.Run("constructs the CLI params correctly - without filter", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Subaccount.List(context.TODO(), "")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
	t.Run("constructs the CLI params correctly - without filter", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
				"labelsFilter":  "name eq 'abc'",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Subaccount.List(context.TODO(), "name eq 'abc'")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsSubaccountFacade_Get(t *testing.T) {
	command := "accounts/subaccount"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
			})

		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Subaccount.Get(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsSubaccountFacade_Create(t *testing.T) {
	command := "accounts/subaccount"

	displayName := "my-account"
	subdomain := "my-account-sub"
	region := "eu30"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"displayName": displayName,
				"subdomain":   subdomain,
				"region":      region,
			})

		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Subaccount.Create(context.TODO(), &SubaccountCreateInput{
			DisplayName: displayName,
			Subdomain:   subdomain,
			Region:      region,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsSubaccountFacade_Update(t *testing.T) {
	command := "accounts/subaccount"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	displayName := "my-account"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUpdate, map[string]string{
				"subaccount":  subaccountId,
				"displayName": displayName,
			})

		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Subaccount.Update(context.TODO(), subaccountId, displayName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsSubaccountFacade_Delete(t *testing.T) {
	command := "accounts/subaccount"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"subaccount":  subaccountId,
				"confirm":     "true",
				"forceDelete": "true",
			})

		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Subaccount.Delete(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsSubaccountFacade_Subscribe(t *testing.T) {
	command := "accounts/subaccount"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	appName := "SAPLaunchpad"
	planName := "free"
	parameters := "{}"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionSubscribe, map[string]string{
				"subaccount":         subaccountId,
				"appName":            appName,
				"planName":           planName,
				"subscriptionParams": parameters,
			})

		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Subaccount.Subscribe(context.TODO(), subaccountId, appName, planName, parameters)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsSubaccountFacade_Unsubscribe(t *testing.T) {
	command := "accounts/subaccount"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	appName := "SAPLaunchpad"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUnsubscribe, map[string]string{
				"subaccount": subaccountId,
				"appName":    appName,
				"confirm":    "true",
			})

		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Subaccount.Unsubscribe(context.TODO(), subaccountId, appName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
