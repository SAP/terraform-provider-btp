package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
)

func TestDirectoryValueFrom(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var obj cis.DirectoryResponseObject
		err := json.Unmarshal([]byte(`
{
  "guid": "04fb4993-350b-4ae3-b2b1-d252bc2cb646",
  "parentType": "ROOT",
  "globalAccountGUID": "795b53bb-a3f0-4769-adf0-26173282a975",
  "displayName": "my-directory",
  "description": "my description",
  "createdDate": 1676030306299,
  "createdBy": "john.doe@sap.com",
  "modifiedDate": 1676030306299,
  "entityState": "OK",
  "directoryType": "FOLDER",
  "directoryFeatures": [
    "DEFAULT"
  ],
  "labels": {"a": ["b"]},
  "subdomain": "my-subdomain",
  "parentGuid": "795b53bb-a3f0-4769-adf0-26173282a975",
  "parentGUID": "795b53bb-a3f0-4769-adf0-26173282a975"
}
		`), &obj)

		if assert.NoError(t, err) {
			uut, diags := directoryValueFrom(context.TODO(), obj)

			assert.False(t, diags.HasError())
			assert.Equal(t, "04fb4993-350b-4ae3-b2b1-d252bc2cb646", uut.ID.ValueString())
			assert.Equal(t, "john.doe@sap.com", uut.CreatedBy.ValueString())
			assert.Equal(t, "2023-02-10T11:58:26Z", uut.CreatedDate.ValueString())
			assert.Equal(t, "my description", uut.Description.ValueString())
			assert.Equal(t, "2023-02-10T11:58:26Z", uut.LastModified.ValueString())
			assert.Equal(t, "my-directory", uut.Name.ValueString())
			assert.Equal(t, "795b53bb-a3f0-4769-adf0-26173282a975", uut.ParentID.ValueString())
			assert.Equal(t, "OK", uut.State.ValueString())
			assert.Equal(t, "my-subdomain", uut.Subdomain.ValueString())
			assert.Equal(t, "[\"DEFAULT\"]", uut.Features.String())
			assert.Equal(t, "{\"a\":[\"b\"]}", uut.Labels.String())
		}
	})
}
