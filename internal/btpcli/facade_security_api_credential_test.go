package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityApiCredential_CreateBySubaccount(t *testing.T){
	
	command := "security/api-credential"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	name := "Subaccount Api-Credentials"
	certificate := "-----BEGIN CERTIFICATE-----\nMock-PEM-Certificate\n-----END CERTIFICATE-----"


	t.Run("constructs the CLI params correctly - client secret", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"subaccount": subaccountId,
				"name": name,
				"readOnly": "false",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.CreateBySubaccount(context.TODO(), &ApiCredentialInput{
			Subaccount: subaccountId,
			Name: name,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})

	t.Run("constructs the CLI params correctly - certificate", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"subaccount": subaccountId,
				"name": name,
				"readOnly": "true",
				"certificate": certificate,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.CreateBySubaccount(context.TODO(), &ApiCredentialInput{
			Subaccount: subaccountId,
			Name: name,
			Certificate: certificate,
			ReadOnly: true,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityApiCredential_DeleteBySubaccount(t *testing.T){
	
	command := "security/api-credential"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	name := "Subaccount Api-Credentials"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"subaccount": subaccountId,
				"name": name,
				"readOnly" : "false",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.DeleteBySubaccount(context.TODO(), &ApiCredentialInput{
			Subaccount: subaccountId,
			Name: name,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityApiCredential_GetBySubaccount(t *testing.T){
	
	command := "security/api-credential"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	name := "Subaccount Api-Credentials"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount": subaccountId,
				"name": name,
				"readOnly" : "false",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.GetBySubaccount(context.TODO(), &ApiCredentialInput{
			Subaccount: subaccountId,
			Name: name,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityApiCredential_CreateByDirectory(t *testing.T){
	
	command := "security/api-credential"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	name := "Directory Api-Credentials"
	certificate := "-----BEGIN CERTIFICATE-----\nMock-PEM-Certificate\n-----END CERTIFICATE-----"

	t.Run("constructs the CLI params correctly - client secret", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"directory": directoryId,
				"name": name,
				"readOnly" : "false",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.CreateByDirectory(context.TODO(), &ApiCredentialInput{
			Directory: directoryId,
			Name: name,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})

	t.Run("constructs the CLI params correctly - certificate", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"directory": directoryId,
				"name": name,
				"readOnly" : "false",
				"certificate" : certificate,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.CreateByDirectory(context.TODO(), &ApiCredentialInput{
			Directory: directoryId,
			Name: name,
			Certificate: certificate,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityApiCredential_DeleteByDirectory(t *testing.T){
	
	command := "security/api-credential"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	name := "Directory Api-Credentials"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"directory": directoryId,
				"name": name,
				"readOnly" : "false",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.DeleteByDirectory(context.TODO(), &ApiCredentialInput{
			Directory: directoryId,
			Name: name,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityApiCredential_GetByDirectory(t *testing.T){
	
	command := "security/api-credential"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	name := "Directory Api-Credentials"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"directory": directoryId,
				"name": name,
				"readOnly" : "false",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.GetByDirectory(context.TODO(), &ApiCredentialInput{
			Directory: directoryId,
			Name: name,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityApiCredential_CreateByGlobalAccount(t *testing.T){
	
	command := "security/api-credential"

	name := "Global Account Api-Credentials"
	certificate := "-----BEGIN CERTIFICATE-----\nMock-PEM-Certificate\n-----END CERTIFICATE-----"


	t.Run("constructs the CLI params correctly - client secret", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
				"name": name,
				"readOnly": "false",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.CreateByGlobalAccount(context.TODO(), &ApiCredentialInput{
			Name: name,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})

	t.Run("constructs the CLI params correctly - certificate", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
				"name": name,
				"readOnly": "true",
				"certificate": certificate,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.CreateByGlobalAccount(context.TODO(), &ApiCredentialInput{
			Name: name,
			Certificate: certificate,
			ReadOnly: true,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityApiCredential_DeleteByGlobalAccount(t *testing.T){
	
	command := "security/api-credential"

	name := "Global Account Api-Credentials"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
				"name": name,
				"readOnly" : "false",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.DeleteByGlobalAccount(context.TODO(), &ApiCredentialInput{
			Name: name,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityApiCredential_GetByGlobalAccount(t *testing.T){
	
	command := "security/api-credential"

	name := "Global Account Api-Credentials"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
				"name": name,
				"readOnly" : "false",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.ApiCredential.GetByGlobalAccount(context.TODO(), &ApiCredentialInput{
			Name: name,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}