package btpcli

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityRoleFacade_ListByGlobalAccount(t *testing.T) {
	command := "security/role"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionList, map[string]string{
				"globalAccount": "795b53bb-a3f0-4769-adf0-26173282a975",
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.ListByGlobalAccount(context.TODO())

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_ListBySubaccount(t *testing.T) {
	command := "security/role"

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

		_, res, err := uut.Security.Role.ListBySubaccount(context.TODO(), subaccountId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_ListByDirectory(t *testing.T) {
	command := "security/role"

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

		_, res, err := uut.Security.Role.ListByDirectory(context.TODO(), directoryId)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_GetByGlobalAccount(t *testing.T) {
	command := "security/role"

	globalAccountId := "795b53bb-a3f0-4769-adf0-26173282a975"
	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"globalAccount":    globalAccountId,
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.GetByGlobalAccount(context.TODO(), roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_GetBySubaccount(t *testing.T) {
	command := "security/role"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"subaccount":       subaccountId,
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.GetBySubaccount(context.TODO(), subaccountId, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_GetByDirectory(t *testing.T) {
	command := "security/role"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionGet, map[string]string{
				"directory":        directoryId,
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.GetByDirectory(context.TODO(), directoryId, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_DeleteByDirectory(t *testing.T) {
	command := "security/role"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"directory":        directoryId,
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.DeleteByDirectory(context.TODO(), directoryId, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_DeleteBySubaccount(t *testing.T) {
	command := "security/role"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"subaccount":       subaccountId,
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.DeleteBySubaccount(context.TODO(), subaccountId, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_DeleteByGlobalAccount(t *testing.T) {
	command := "security/role"

	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionDelete, map[string]string{
				"globalAccount":    "795b53bb-a3f0-4769-adf0-26173282a975",
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.DeleteByGlobalAccount(context.TODO(), roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_CreateByDirectory(t *testing.T) {
	command := "security/role"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"directory":        directoryId,
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.CreateByDirectory(context.TODO(), &DirectoryRoleCreateInput{
			RoleName:         roleName,
			AppId:            roleTemplateAppId,
			RoleTemplateName: roleTemplateName,
			DirectoryId:      directoryId,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_CreateBySubaccount(t *testing.T) {
	command := "security/role"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"subaccount":       "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f",
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.CreateBySubaccount(context.TODO(), &SubaccountRoleCreateInput{
			RoleName:         roleName,
			AppId:            roleTemplateAppId,
			RoleTemplateName: roleTemplateName,
			SubaccountId:     subaccountId,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_CreateByGlobalAccount(t *testing.T) {
	command := "security/role"

	roleName := "User and Role Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionCreate, map[string]string{
				"globalAccount":    "795b53bb-a3f0-4769-adf0-26173282a975",
				"appId":            roleTemplateAppId,
				"roleName":         roleName,
				"roleTemplateName": roleTemplateName,
			})
		}))
		defer srv.Close()

		_, res, err := uut.Security.Role.CreateByGlobalAccount(context.TODO(), &GlobalAccountRoleCreateInput{
			RoleName:         roleName,
			AppId:            roleTemplateAppId,
			RoleTemplateName: roleTemplateName,
		})

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_AddByGlobalAccount(t *testing.T) {
	command := "security/role"

	roleCollectionName := "my-role-collection"
	roleName := "XSUAA Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAdd, map[string]string{
				"globalAccount":      "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName": roleCollectionName,
				"roleName":           roleName,
				"roleTemplateAppID":  roleTemplateAppId,
				"roleTemplateName":   roleTemplateName,
			})
		}))
		defer srv.Close()

		res, err := uut.Security.Role.AddByGlobalAccount(context.TODO(), roleCollectionName, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_AddBySubaccount(t *testing.T) {
	command := "security/role"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my-role-collection"
	roleName := "XSUAA Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAdd, map[string]string{
				"subaccount":         subaccountId,
				"roleCollectionName": roleCollectionName,
				"roleName":           roleName,
				"roleTemplateAppID":  roleTemplateAppId,
				"roleTemplateName":   roleTemplateName,
			})
		}))
		defer srv.Close()

		res, err := uut.Security.Role.AddBySubaccount(context.TODO(), subaccountId, roleCollectionName, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_AddByDirectory(t *testing.T) {
	command := "security/role"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	roleCollectionName := "my-role-collection"
	roleName := "XSUAA Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionAdd, map[string]string{
				"directory":          directoryId,
				"roleCollectionName": roleCollectionName,
				"roleName":           roleName,
				"roleTemplateAppID":  roleTemplateAppId,
				"roleTemplateName":   roleTemplateName,
			})
		}))
		defer srv.Close()

		res, err := uut.Security.Role.AddByDirectory(context.TODO(), directoryId, roleCollectionName, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_RemoveByGlobalAccount(t *testing.T) {
	command := "security/role"

	roleCollectionName := "my-role-collection"
	roleName := "XSUAA Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionRemove, map[string]string{
				"globalAccount":      "795b53bb-a3f0-4769-adf0-26173282a975",
				"roleCollectionName": roleCollectionName,
				"roleName":           roleName,
				"roleTemplateAppID":  roleTemplateAppId,
				"roleTemplateName":   roleTemplateName,
			})
		}))
		defer srv.Close()

		res, err := uut.Security.Role.RemoveByGlobalAccount(context.TODO(), roleCollectionName, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_RemoveBySubaccount(t *testing.T) {
	command := "security/role"

	subaccountId := "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	roleCollectionName := "my-role-collection"
	roleName := "XSUAA Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionRemove, map[string]string{
				"subaccount":         subaccountId,
				"roleCollectionName": roleCollectionName,
				"roleName":           roleName,
				"roleTemplateAppID":  roleTemplateAppId,
				"roleTemplateName":   roleTemplateName,
			})
		}))
		defer srv.Close()

		res, err := uut.Security.Role.RemoveBySubaccount(context.TODO(), subaccountId, roleCollectionName, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}

func TestSecurityRoleFacade_RemoveByDirectory(t *testing.T) {
	command := "security/role"

	directoryId := "f6c7137d-c5a0-48c2-b2a4-fd64e6b35d3d"
	roleCollectionName := "my-role-collection"
	roleName := "XSUAA Auditor"
	roleTemplateAppId := "xsuaa!t1"
	roleTemplateName := "xsuaa_auditor"

	t.Run("constructs the CLI params correctly", func(t *testing.T) {
		var srvCalled bool

		uut, srv := prepareClientFacadeForTest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			srvCalled = true

			assertCall(t, r, command, ActionRemove, map[string]string{
				"directory":          directoryId,
				"roleCollectionName": roleCollectionName,
				"roleName":           roleName,
				"roleTemplateAppID":  roleTemplateAppId,
				"roleTemplateName":   roleTemplateName,
			})
		}))
		defer srv.Close()

		res, err := uut.Security.Role.RemoveByDirectory(context.TODO(), directoryId, roleCollectionName, roleName, roleTemplateAppId, roleTemplateName)

		if assert.True(t, srvCalled) && assert.NoError(t, err) {
			assert.Equal(t, 200, res.StatusCode)
		}
	})
}
