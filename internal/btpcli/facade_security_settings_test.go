package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecuritySettingsFacade_ListByGlobalAccount(t *testing.T) {
	command := "security/settings"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Settings.ListByGlobalAccount(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecuritySettingsFacade_ListBySubaccount(t *testing.T) {
	command := "security/settings"

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

		_, res, err := uut.Security.Settings.ListBySubaccount(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecuritySettingsFacade_UpdateByGlobalAccount(t *testing.T) {
	command := "security/settings"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUpdate, map[string]string{
				"globalAccount":                     globalAccountId,
				"iFrameDomain":                      "https://my-iframe-domain-1",
				"customEmailDomains":                "[\"customemaildomain1.com\",\"customemaildomain2.com\"]",
				"defaultIdp":                        "my-idp",
				"treatUsersWithSameEmailAsSameUser": "true",
				"homeRedirect":                      "my-new-redirect",
				"accessTokenValidity":               "3600",
				"refreshTokenValidity":              "3600",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Settings.UpdateByGlobalAccount(context.TODO(), SecuritySettingsUpdateInput{
			IFrame:                            "https://my-iframe-domain-1",
			CustomEmail:                       []string{"customemaildomain1.com", "customemaildomain2.com"},
			DefaultIDPForNonInteractiveLogon:  "my-idp",
			TreatUsersWithSameEmailAsSameUser: true,
			HomeRedirect:                      "my-new-redirect",
			AccessTokenValidity:               3600,
			RefreshTokenValidity:              3600,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecuritySettingsFacade_UpdateBySubaccount(t *testing.T) {
	command := "security/settings"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUpdate, map[string]string{
				"subaccount":                        subaccountId,
				"iFrameDomain":                      "https://my-iframe-domain-1",
				"customEmailDomains":                "[\"customemaildomain1.com\",\"customemaildomain2.com\"]",
				"defaultIdp":                        "my-idp",
				"treatUsersWithSameEmailAsSameUser": "true",
				"homeRedirect":                      "my-new-redirect",
				"accessTokenValidity":               "3600",
				"refreshTokenValidity":              "3600",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Settings.UpdateBySubaccount(context.TODO(), subaccountId, SecuritySettingsUpdateInput{
			IFrame:                            "https://my-iframe-domain-1",
			CustomEmail:                       []string{"customemaildomain1.com", "customemaildomain2.com"},
			DefaultIDPForNonInteractiveLogon:  "my-idp",
			TreatUsersWithSameEmailAsSameUser: true,
			HomeRedirect:                      "my-new-redirect",
			AccessTokenValidity:               3600,
			RefreshTokenValidity:              3600,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
