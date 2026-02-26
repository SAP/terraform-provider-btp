package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisasterRecoverySubaccountPairFacade_Get(t *testing.T) {
	command := "disaster-recovery/subaccount-pair"

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

		_, res, err := uut.DisasterRecovery.SubaccountPair.Get(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestDisasterRecoverySubaccountPairFacade_Create(t *testing.T) {
	command := "disaster-recovery/subaccount-pair"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	withSubaccountId := "6D079379-6442-464A-90EB-65FAC05B176F"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"subaccount":      subaccountId,
				"with-subaccount": withSubaccountId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.DisasterRecovery.SubaccountPair.Create(context.TODO(), &SubaccountPairCreateInput{
			SubaccountId:     subaccountId,
			WithSubaccountId: withSubaccountId,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestDisasterRecoverySubaccountPairFacade_Delete(t *testing.T) {
	command := "disaster-recovery/subaccount-pair"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"subaccount": subaccountId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.DisasterRecovery.SubaccountPair.Delete(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
