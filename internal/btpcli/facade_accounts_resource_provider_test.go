package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsResourceProviderFacade_List(t *testing.T) {
	command := "accounts/resource-provider"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.ResourceProvider.List(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsResourceProviderFacade_Get(t *testing.T) {
	command := "accounts/resource-provider"

	resourceProvider := "AWS"
	resourceTechnicalName := "my_id"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
				"provider":      "AWS",
				"technicalName": "my_id",
			})

		}))
		defer srv.Close()

		_, res, err := uut.Accounts.ResourceProvider.Get(context.TODO(), resourceProvider, resourceTechnicalName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsResourceProviderFacade_Create(t *testing.T) {
	command := "accounts/resource-provider"

	provider := "AWS"
	technicalName := "my_id"
	description := "my-description"
	displayName := "My display name"
	configuration := "{}"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"globalAccount":     "795b53bb-a3f0-4769-adf0-26173282a975",
				"provider":          provider,
				"technicalName":     technicalName,
				"description":       description,
				"displayName":       displayName,
				"configurationInfo": configuration,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.ResourceProvider.Create(context.TODO(), GlobalaccountResourceProviderCreateInput{
			Provider:      provider,
			TechnicalName: technicalName,
			Description:   description,
			DisplayName:   displayName,
			Configuration: configuration,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsResourceProviderFacade_Delete(t *testing.T) {
	command := "accounts/resource-provider"

	provider := "AWS"
	technicalName := "my_id"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
				"provider":      "AWS",
				"technicalName": "my_id",
				"confirm":       "true",
			})

		}))
		defer srv.Close()

		_, res, err := uut.Accounts.ResourceProvider.Delete(context.TODO(), provider, technicalName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
