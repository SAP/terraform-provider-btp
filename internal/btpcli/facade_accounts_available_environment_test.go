package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsAvailableEnvironmentFacade_List(t *testing.T) {
	command := "accounts/available-environment"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"subaccount": "test",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.AvailableEnvironment.List(context.TODO(), "test")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
