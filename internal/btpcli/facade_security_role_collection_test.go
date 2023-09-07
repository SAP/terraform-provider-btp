package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityRoleCollectionFacade_ListByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.ListByGlobalAccount(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_ListBySubaccount(t *testing.T) {
	command := "security/role-collection"

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

		_, res, err := uut.Security.RoleCollection.ListBySubaccount(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_ListByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"directory": directoryId,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.ListByDirectory(context.TODO(), directoryId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_GetByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	roleCollectionName := "Global Account Administrator"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"globalAccount":      globalAccountId,
				"roleCollectionName": roleCollectionName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.GetByGlobalAccount(context.TODO(), roleCollectionName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_GetBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "Subaccount Administrator"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount":         subaccountId,
				"roleCollectionName": roleCollectionName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.GetBySubaccount(context.TODO(), subaccountId, roleCollectionName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_GetByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	roleCollectionName := "Directory Administrator"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"directory":          directoryId,
				"roleCollectionName": roleCollectionName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.GetByDirectory(context.TODO(), directoryId, roleCollectionName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_CreateByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	roleCollectionName := "my own rolecollection"
	description := "This is the description of my own rolecollection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"globalAccount":      "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName": roleCollectionName,
				"description":        description,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.CreateByGlobalAccount(context.TODO(), roleCollectionName, description)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_CreateBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	description := "This is the description of my own rolecollection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"subaccount":         subaccountId,
				"roleCollectionName": roleCollectionName,
				"description":        description,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.CreateBySubaccount(context.TODO(), subaccountId, roleCollectionName, description)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_CreateByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	description := "This is the description of my own rolecollection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"directory":          directoryId,
				"roleCollectionName": roleCollectionName,
				"description":        description,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.CreateByDirectory(context.TODO(), directoryId, roleCollectionName, description)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UpdateByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	roleCollectionName := "my own rolecollection"
	description := "This is the updated description of my own rolecollection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUpdate, map[string]string{
				"globalAccount":      "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName": roleCollectionName,
				"description":        description,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UpdateByGlobalAccount(context.TODO(), roleCollectionName, description)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UpdateBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	description := "This is the updated description of my own rolecollection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUpdate, map[string]string{
				"subaccount":         subaccountId,
				"roleCollectionName": roleCollectionName,
				"description":        description,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UpdateBySubaccount(context.TODO(), subaccountId, roleCollectionName, description)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UpdateByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	description := "This is the updated description of my own rolecollection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUpdate, map[string]string{
				"directory":          directoryId,
				"roleCollectionName": roleCollectionName,
				"description":        description,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UpdateByDirectory(context.TODO(), directoryId, roleCollectionName, description)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_DeleteByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	roleCollectionName := "my own rolecollection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"globalAccount":      "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName": roleCollectionName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.DeleteByGlobalAccount(context.TODO(), roleCollectionName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_DeleteBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"subaccount":         subaccountId,
				"roleCollectionName": roleCollectionName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.DeleteBySubaccount(context.TODO(), subaccountId, roleCollectionName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_DeleteByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"directory":          directoryId,
				"roleCollectionName": roleCollectionName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.DeleteByDirectory(context.TODO(), directoryId, roleCollectionName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_AssignUserByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	roleCollectionName := "my own rolecollection"
	userName := "john.doe@test.com"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"globalAccount":       "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName":  roleCollectionName,
				"userName":            userName,
				"origin":              origin,
				"createUserIfMissing": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.AssignUserByGlobalaccount(context.TODO(), roleCollectionName, userName, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_AssignUserBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	userName := "john.doe@test.com"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"subaccount":          subaccountId,
				"roleCollectionName":  roleCollectionName,
				"userName":            userName,
				"origin":              origin,
				"createUserIfMissing": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.AssignUserBySubaccount(context.TODO(), subaccountId, roleCollectionName, userName, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_AssignUserByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	userName := "john.doe@test.com"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"directory":           directoryId,
				"roleCollectionName":  roleCollectionName,
				"userName":            userName,
				"origin":              origin,
				"createUserIfMissing": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.AssignUserByDirectory(context.TODO(), directoryId, roleCollectionName, userName, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_AssignGroupByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	roleCollectionName := "my own rolecollection"
	group := "my/group"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"globalAccount":       "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName":  roleCollectionName,
				"group":               group,
				"origin":              origin,
				"createUserIfMissing": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.AssignGroupByGlobalaccount(context.TODO(), roleCollectionName, group, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_AssignGroupBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	group := "my/group"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"subaccount":          subaccountId,
				"roleCollectionName":  roleCollectionName,
				"group":               group,
				"origin":              origin,
				"createUserIfMissing": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.AssignGroupBySubaccount(context.TODO(), subaccountId, roleCollectionName, group, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_AssignGroupByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	group := "my/group"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"directory":           directoryId,
				"roleCollectionName":  roleCollectionName,
				"group":               group,
				"origin":              origin,
				"createUserIfMissing": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.AssignGroupByDirectory(context.TODO(), directoryId, roleCollectionName, group, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UnassignUserByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	roleCollectionName := "my own rolecollection"
	userName := "john.doe@test.com"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUnassign, map[string]string{
				"globalAccount":      "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName": roleCollectionName,
				"userName":           userName,
				"origin":             origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UnassignUserByGlobalaccount(context.TODO(), roleCollectionName, userName, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UnassignUserBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	userName := "john.doe@test.com"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUnassign, map[string]string{
				"subaccount":         subaccountId,
				"roleCollectionName": roleCollectionName,
				"userName":           userName,
				"origin":             origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UnassignUserBySubaccount(context.TODO(), subaccountId, roleCollectionName, userName, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UnassignUserByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	userName := "john.doe@test.com"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUnassign, map[string]string{
				"directory":          directoryId,
				"roleCollectionName": roleCollectionName,
				"userName":           userName,
				"origin":             origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UnassignUserByDirectory(context.TODO(), directoryId, roleCollectionName, userName, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UnassignGroupByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	roleCollectionName := "my own rolecollection"
	group := "my/group"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUnassign, map[string]string{
				"globalAccount":      "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName": roleCollectionName,
				"group":              group,
				"origin":             origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UnassignGroupByGlobalaccount(context.TODO(), roleCollectionName, group, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UnassignGroupBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	group := "my/group"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUnassign, map[string]string{
				"subaccount":         subaccountId,
				"roleCollectionName": roleCollectionName,
				"group":              group,
				"origin":             origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UnassignGroupBySubaccount(context.TODO(), subaccountId, roleCollectionName, group, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UnassignGroupByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	group := "my/group"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUnassign, map[string]string{
				"directory":          directoryId,
				"roleCollectionName": roleCollectionName,
				"group":              group,
				"origin":             origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UnassignGroupByDirectory(context.TODO(), directoryId, roleCollectionName, group, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_AssignAttributeByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	roleCollectionName := "my own rolecollection"
	attributeName := "my/attributename"
	attributeValue := "my/attributevalue"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"globalAccount":       "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName":  roleCollectionName,
				"attributeName":       attributeName,
				"attributeValue":      attributeValue,
				"origin":              origin,
				"createUserIfMissing": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.AssignAttributeByGlobalaccount(context.TODO(), roleCollectionName, attributeName, attributeValue, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_AssignAttributeBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	attributeName := "my/attributename"
	attributeValue := "my/attributevalue"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"subaccount":          subaccountId,
				"roleCollectionName":  roleCollectionName,
				"attributeName":       attributeName,
				"attributeValue":      attributeValue,
				"origin":              origin,
				"createUserIfMissing": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.AssignAttributeBySubaccount(context.TODO(), subaccountId, roleCollectionName, attributeName, attributeValue, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_AssignAttributeByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	attributeName := "my/attributename"
	attributeValue := "my/attributevalue"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAssign, map[string]string{
				"directory":           directoryId,
				"roleCollectionName":  roleCollectionName,
				"attributeName":       attributeName,
				"attributeValue":      attributeValue,
				"origin":              origin,
				"createUserIfMissing": "true",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.AssignAttributeByDirectory(context.TODO(), directoryId, roleCollectionName, attributeName, attributeValue, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UnassignAttributeByGlobalAccount(t *testing.T) {
	command := "security/role-collection"

	roleCollectionName := "my own rolecollection"
	attributeName := "my/attributename"
	attributeValue := "my/attributevalue"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUnassign, map[string]string{
				"globalAccount":      "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName": roleCollectionName,
				"attributeName":      attributeName,
				"attributeValue":     attributeValue,
				"origin":             origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UnassignAttributeByGlobalaccount(context.TODO(), roleCollectionName, attributeName, attributeValue, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UnassignAttributeBySubaccount(t *testing.T) {
	command := "security/role-collection"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	attributeName := "my/attributename"
	attributeValue := "my/attributevalue"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUnassign, map[string]string{
				"subaccount":         subaccountId,
				"roleCollectionName": roleCollectionName,
				"attributeName":      attributeName,
				"attributeValue":     attributeValue,
				"origin":             origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UnassignAttributeBySubaccount(context.TODO(), subaccountId, roleCollectionName, attributeName, attributeValue, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleCollectionFacade_UnassignAttributeByDirectory(t *testing.T) {
	command := "security/role-collection"

	directoryId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my own rolecollection"
	attributeName := "my/attributename"
	attributeValue := "my/attributevalue"
	origin := "ldap"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionUnassign, map[string]string{
				"directory":          directoryId,
				"roleCollectionName": roleCollectionName,
				"attributeName":      attributeName,
				"attributeValue":     attributeValue,
				"origin":             origin,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.RoleCollection.UnassignAttributeByDirectory(context.TODO(), directoryId, roleCollectionName, attributeName, attributeValue, origin)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
