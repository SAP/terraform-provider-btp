package provider

import (
	"context"
	"encoding/json"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_trust"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGlobalaccountTrustConfigurationFromValue(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var obj xsuaa_trust.TrustConfigurationResponseObject
		err := json.Unmarshal([]byte(`
{
  "name": "sap.default",
  "originKey": "sap.default",
  "typeOfTrust": "Application",
  "status": "active",
  "description": null,
  "identityProvider": null,
  "domain": null,
  "linkTextForUserLogon": "Default Identity Provider",
  "availableForUserLogon": "true",
  "createShadowUsersDuringLogon": "true",
  "sapBtpCli": null,
  "protocol": "OpenID Connect",
  "readOnly": true
}
		`), &obj)

		if assert.NoError(t, err) {
			uut, diags := globalaccountTrustConfigurationFromValue(context.TODO(), obj)

			assert.False(t, diags.HasError())
			assert.Equal(t, "sap.default", uut.Origin.ValueString())
			assert.Equal(t, "sap.default", uut.Id.ValueString())
			assert.Equal(t, "sap.default", uut.Name.ValueString())
			assert.Equal(t, "", uut.Description.ValueString())
			assert.Equal(t, "Application", uut.Type.ValueString())
			assert.Equal(t, "", uut.IdentityProvider.ValueString())
			assert.Equal(t, "OpenID Connect", uut.Protocol.ValueString())
			assert.Equal(t, "active", uut.Status.ValueString())
			assert.True(t, uut.ReadOnly.ValueBool())
		}
	})
}
