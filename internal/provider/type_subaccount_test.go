package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
)

func TestSubaccountValueFrom(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var obj cis.SubaccountResponseObject
		err := json.Unmarshal([]byte(`
{
  "guid": "64eb181c-d4d9-4b28-a9bf-063374eaebeb",
  "technicalName": "64eb181c-d4d9-4b28-a9bf-063374eaebeb",
  "displayName": "my-subaccount",
  "globalAccountGUID": "795b53bb-a3f0-4769-adf0-26173282a975",
  "parentGUID": "795b53bb-a3f0-4769-adf0-26173282a975",
  "parentType": "ROOT",
  "region": "eu30",
  "subdomain": "my-subaccount-k3j4dbrq",
  "betaEnabled": false,
  "usedForProduction": "NOT_USED_FOR_PRODUCTION",
  "description": null,
  "state": "OK",
  "contentAutomationState": null,
  "contentAutomationStateDetails": null,
  "createdDate": 1674480698934,
  "createdBy": "john.doe@mycompany.com",
  "modifiedDate": 1674480712350,
  "labels": {"a": ["b"]},
  "parentFeatures": ["DEFAULT"]
}
		`), &obj)

		if assert.NoError(t, err) {
			uut, diags := subaccountValueFrom(context.TODO(), obj)

			assert.False(t, diags.HasError())
			assert.Equal(t, "64eb181c-d4d9-4b28-a9bf-063374eaebeb", uut.ID.ValueString())
			assert.Equal(t, false, uut.BetaEnabled.ValueBool())
			assert.Equal(t, "john.doe@mycompany.com", uut.CreatedBy.ValueString())
			assert.Equal(t, "2023-01-23T13:31:38Z", uut.CreatedDate.ValueString())
			assert.Equal(t, "", uut.Description.ValueString())
			assert.Equal(t, "{\"a\":[\"b\"]}", uut.Labels.String())
			assert.Equal(t, "2023-01-23T13:31:52Z", uut.LastModified.ValueString())
			assert.Equal(t, "my-subaccount", uut.Name.ValueString())
			assert.Equal(t, "795b53bb-a3f0-4769-adf0-26173282a975", uut.ParentID.ValueString())
			assert.Equal(t, "[\"DEFAULT\"]", uut.ParentFeatures.String())
			assert.Equal(t, "eu30", uut.Region.ValueString())
			assert.Equal(t, "OK", uut.State.ValueString())
			assert.Equal(t, "my-subaccount-k3j4dbrq", uut.Subdomain.ValueString())
			assert.Equal(t, "NOT_USED_FOR_PRODUCTION", uut.Usage.ValueString())
		}
	})
}
