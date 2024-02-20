package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountsDirectoryFacade_Get(t *testing.T) {
	command := "accounts/directory"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"directoryID":   "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0",
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Directory.Get(context.TODO(), "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsDirectoryFacade_Create(t *testing.T) {
	command := "accounts/directory"
	globalAccount := "795b53bb-a3f0-4769-adf0-26173282a975"

	displayName := "my-directory"
	description := "a description"
	subdomain := "my-directory-subdomain"
	parentId := "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0"

	t.Run("constructs the CLI params correctly - minimal", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"globalAccount": globalAccount,
				"displayName":   displayName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Directory.Create(context.TODO(), &DirectoryCreateInput{
			DisplayName: displayName,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
	t.Run("constructs the CLI params correctly - full", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"globalAccount": globalAccount,
				"displayName":   displayName,
				"description":   description,
				"subdomain":     subdomain,
				"parentID":      parentId,
				"labels":        "{}",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Directory.Create(context.TODO(), &DirectoryCreateInput{
			DisplayName: displayName,
			Description: &description,
			Subdomain:   &subdomain,
			ParentID:    &parentId,
			Labels:      map[string][]string{},
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsDirectoryFacade_Update(t *testing.T) {
	command := "accounts/directory"
	globalAccount := "795b53bb-a3f0-4769-adf0-26173282a975"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	displayName := "my-directory"
	description := "a description"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUpdate, map[string]string{
				"globalAccount": globalAccount,
				"directoryID":   directoryId,
				"displayName":   displayName,
				"description":   description,
				"labels":        "{}",
			})

		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Directory.Update(context.TODO(), &DirectoryUpdateInput{
			DirectoryId: directoryId,
			DisplayName: &displayName,
			Description: &description,
			Labels:      map[string][]string{},
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsDirectoryFacade_Enable(t *testing.T) {
	command := "accounts/directory"
	globalAccount := "795b53bb-a3f0-4769-adf0-26173282a975"
	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	subdomain := "my-directory-subdomain"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionEnable, map[string]string{
				"globalAccount":     globalAccount,
				"directoryID":       directoryId,
				"subdomain":         subdomain,
				"directoryFeatures": "E,A",
				"confirm":           "true",
			})

		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Directory.Enable(context.TODO(), &DirectoryEnableInput{
			DirectoryId: directoryId,
			Subdomain:   &subdomain,
			Features:    []string{"E", "A"},
			Confirm:     true,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestAccountsDirectoryFacade_Delete(t *testing.T) {
	command := "accounts/directory"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"directoryID":   "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0",
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
				"confirm":       "true",
				"forceDelete":   "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Accounts.Directory.Delete(context.TODO(), "dd005d8b-1fee-4e6b-b6ff-cb9a197b7fe0")

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
