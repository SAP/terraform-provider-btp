package provider

import (
	"context"
	"encoding/json"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGlobalaccountRoleFromValue(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var obj xsuaa_authz.Role
		err := json.Unmarshal([]byte(`
{
  "roleTemplateName": "xsuaa_admin",
  "roleTemplateAppId": "xsuaa!t1",
  "name": "My Role",
  "identityZone": "795b53bb-a3f0-4769-adf0-26173282a975",
  "attributeList": [],
  "creationType": "ADMIN",
  "roleTemplate": {
    "appId": "xsuaa!t1",
    "name": "xsuaa_admin",
    "description": "Manage authorizations, trusted identity providers, and users.",
    "default-role-name": "User and Role Administrator",
    "scope-references": [
      "xs_account.access",
      "xs_authorization.read",
      "xs_authorization.write",
      "xs_idp.read",
      "xs_idp.write",
      "xs_user.read",
      "xs_user.write"
    ],
    "attribute-references": [],
    "attributeReferences": [],
    "defaultRoleName": "User and Role Administrator",
    "scopeReferences": [
      "xs_account.access",
      "xs_authorization.read",
      "xs_authorization.write",
      "xs_idp.read",
      "xs_idp.write",
      "xs_user.read",
      "xs_user.write"
    ]
  }
}
		`), &obj)

		if assert.NoError(t, err) {
			uut, diags := globalaccountRoleFromValue(context.TODO(), obj)

			assert.False(t, diags.HasError())
			assert.Equal(t, "xsuaa_admin", uut.RoleTemplateName.ValueString())
			assert.Equal(t, "xsuaa!t1", uut.RoleTemplateAppId.ValueString())
			assert.Equal(t, "My Role", uut.Name.ValueString())
			assert.False(t, uut.IsReadOnly.ValueBool())
			assert.Equal(t, "", uut.Description.ValueString())
		}
	})
}
