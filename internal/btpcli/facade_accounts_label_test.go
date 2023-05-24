package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsLabelFacade_ListBySubaccount(t *testing.T) {
	command := "accounts/label"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": globalAccountId,
				"subaccountID":  subaccountId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Label.ListBySubaccount(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsLabelFacade_ListByDirectory(t *testing.T) {
	command := "accounts/label"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": globalAccountId,
				"directoryID":   directoryId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Label.ListByDirectory(context.TODO(), directoryId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
