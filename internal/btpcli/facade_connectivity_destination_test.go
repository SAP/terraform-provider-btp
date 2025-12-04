package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectivityDestinationFacade_Get(t *testing.T) {
	command := "connectivity/destination"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	name := "destTest"

	t.Run("constructs params correctly without serviceInstance", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"name":       name,
				"subaccount": subaccountId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Connectivity.Destination.GetBySubaccount(context.TODO(), subaccountId, name, "")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})

	t.Run("constructs params correctly with serviceInstance", func(t *testing.T) {
		var srvCalled bool

		serviceInstance := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"name":            name,
				"subaccount":      subaccountId,
				"serviceInstance": serviceInstance,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Connectivity.Destination.GetBySubaccount(context.TODO(), subaccountId, name, serviceInstance)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
