package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectivityDestinationTrustFacade_GetBySubaccount(t *testing.T) {
	command := "connectivity/destination-trust"

	subaccountId := "12345678-aaaa-bbbb-cccc-123456789000"

	tests := []struct {
		name        string
		trustType   bool
		wantPassive string
	}{
		{
			name:        "trustType=active → passive=false",
			trustType:   true,
			wantPassive: "false",
		},
		{
			name:        "trustType=passive → passive=true",
			trustType:   false,
			wantPassive: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var srvCalled bool

			uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				srvCalled = true

				assertCall(t, r, command, ActionGet, map[string]string{
					"subaccount": subaccountId,
					"passive":    tt.wantPassive,
				})
			}))
			defer srv.Close()

			_, res, err := uut.Connectivity.DestinationTrust.GetBySubaccount(context.TODO(), subaccountId, tt.trustType)

			if assert.True(t, srvCalled) && assert.NoError(t, err) {
				assert.Equal(t, 200, res.StatusCode)
			}
		})
	}
}
