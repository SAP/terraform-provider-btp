package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/provisioning"
)

func TestGlobalaccountResourceProviderValueFrom(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var obj provisioning.ResourceProviderResponseObject
		err := json.Unmarshal([]byte(`
{
  "technicalName": "my_resource_provider",
  "displayName": "My AWS Resource Provider",
  "description": "My description",
  "resourceType": "IAAS_ACCOUNT",
  "resourceProvider": "AWS",
  "additionalInfo": {
    "access_key_id": "AWSACCESSKEY",
    "secret_access_key": "AWSSECRETKEY",
    "vpc_id": "vpc-test",
    "region": "us-east-1"
  }
}
		`), &obj)

		if assert.NoError(t, err) {
			uut, diags := globalaccountResourceProviderValueFrom(context.TODO(), obj)

			assert.False(t, diags.HasError())

			assert.Equal(t, "AWS", uut.ResourceProvider.ValueString())
			assert.Equal(t, "my_resource_provider", uut.Id.ValueString())
			assert.Equal(t, "My AWS Resource Provider", uut.DisplayName.ValueString())
			assert.Equal(t, "My description", uut.Description.ValueString())
			assert.Equal(t, "{\n    \"access_key_id\": \"AWSACCESSKEY\",\n    \"secret_access_key\": \"AWSSECRETKEY\",\n    \"vpc_id\": \"vpc-test\",\n    \"region\": \"us-east-1\"\n  }", uut.Parameters.ValueString())
		}
	})
}
