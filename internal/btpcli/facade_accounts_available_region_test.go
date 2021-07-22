package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsAvailableRegionFacade_List(t *testing.T) {
	command := "accounts/available-region"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.AvailableRegion.List(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
