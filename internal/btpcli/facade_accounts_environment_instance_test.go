package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsEnvironmentInstanceFacade_List(t *testing.T) {
	command := "accounts/environment-instance"

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

		_, res, err := uut.Accounts.EnvironmentInstance.List(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEnvironmentInstanceFacade_Get(t *testing.T) {
	command := "accounts/environment-instance"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	environmentId := "6D079379-6442-464A-90EB-65FAC05B176F"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount":    subaccountId,
				"environmentID": environmentId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.EnvironmentInstance.Get(context.TODO(), subaccountId, environmentId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEnvironmentInstanceFacade_Create(t *testing.T) {
	command := "accounts/environment-instance"

	displayName := "my-instance-name"
	environmentType := "cloudfoundry"
	service := "cloudfoundry"
	landscape := "a-landscape"
	parameters := "{}"
	plan := "free"
	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"subaccount":      subaccountId,
				"displayName":     displayName,
				"environmentType": environmentType,
				"landscapeLabel":  landscape,
				"plan":            plan,
				"service":         service,
				"parameters":      parameters,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.EnvironmentInstance.Create(context.TODO(), &SubaccountEnvironmentInstanceCreateInput{
			DisplayName:     displayName,
			EnvironmentType: environmentType,
			Landscape:       landscape,
			Parameters:      parameters,
			Plan:            plan,
			Service:         service,
			SubaccountID:    subaccountId,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsEnvironmentInstanceFacade_Delete(t *testing.T) {
	command := "accounts/environment-instance"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	environmentId := "6D079379-6442-464A-90EB-65FAC05B176F"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"subaccount":    subaccountId,
				"environmentID": environmentId,
				"confirm":       "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.EnvironmentInstance.Delete(context.TODO(), subaccountId, environmentId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
