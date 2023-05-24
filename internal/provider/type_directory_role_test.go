package provider

import (
	"context"
	"encoding/json"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectoryRoleFromValue(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var obj xsuaa_authz.Role
		err := json.Unmarshal([]byte(`
{
  "roleTemplateName": "Directory_Usage_Reporting_Viewer",
  "roleTemplateAppId": "uas!b36585",
  "name": "DirUsageRepViewTest",
  "identityZone": "ddfc2206-5f11-48ed-a1ec-29010af70050",
  "attributeList": [],
  "creationType": "ADMIN",
  "roleTemplate": {
    "appId": "uas!b36585",
    "name": "Directory_Usage_Reporting_Viewer",
    "description": "Role for directory members with read-only authorizations for core commercialization operations, such as viewing directory usage information.",
    "default-role-name": "Directory Usage Reporting Viewer",
    "scope-references": [
      "uas!b36585.UAS.reporting.directory.read",
      "xs_account.access"
    ],
    "attribute-references": [],
    "defaultRoleName": "Directory Usage Reporting Viewer",
    "scopeReferences": [
      "uas!b36585.UAS.reporting.directory.read",
      "xs_account.access"
    ],
    "attributeReferences": []
  }
}
		`), &obj)

		if assert.NoError(t, err) {
			uut, diags := directoryRoleFromValue(context.TODO(), obj)

			assert.False(t, diags.HasError())
			assert.Equal(t, "Directory_Usage_Reporting_Viewer", uut.RoleTemplateName.ValueString())
			assert.Equal(t, "uas!b36585", uut.RoleTemplateAppId.ValueString())
			assert.Equal(t, "DirUsageRepViewTest", uut.Name.ValueString())
			assert.False(t, uut.IsReadOnly.ValueBool())
			assert.Equal(t, "", uut.Description.ValueString())
		}
	})
}
