package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
)

func TestSubaccountServiceInstanceValueFrom(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var obj servicemanager.ServiceInstanceResponseObject
		err := json.Unmarshal([]byte(`
{
  "id": "76765dca-6683-473a-8f42-809e33a2ea68",
  "ready": true,
  "last_operation": {
    "id": "b6d7d982-8ff5-4de2-baeb-990837d00750",
    "ready": true,
    "type": "create",
    "state": "succeeded",
    "resource_id": "76765dca-6683-473a-8f42-809e33a2ea68",
    "resource_type": "/v1/service_instances",
    "platform_id": "service-manager",
    "correlation_id": "48fe2fe9-3006-4d61-438a-225733e451ed",
    "reschedule": false,
    "reschedule_timestamp": "0001-01-01T00:00:00Z",
    "deletion_scheduled": "0001-01-01T00:00:00Z",
    "created_at": "2022-09-29T21:37:52.0201Z",
    "updated_at": "2022-09-29T21:37:53.535119Z"
  },
  "name": "default",
  "service_plan_id": "dc91ce2b-55c5-42eb-ad92-70aee582d827",
  "platform_id": "service-manager",
  "context": {
    "global_account_id": "795b53bb-a3f0-4769-adf0-26173282a975",
    "subaccount_id": "98cbd1c8-49e2-42d5-8266-980e3e8728a4",
    "crm_customer_id": "",
    "license_type": "DEVELOPER",
    "subdomain": "my-subdomain",
    "region": "cf-eu10",
    "platform": "sapcp",
    "origin": "sapcp",
    "zone_id": "98cbd1c8-49e2-42d5-8266-980e3e8728a4",
    "instance_name": "default"
  },
  "usable": true,
  "subaccount_id": "98cbd1c8-49e2-42d5-8266-980e3e8728a4",
  "created_at": "2022-09-29T21:37:52.020097Z",
  "updated_at": "2022-09-29T21:37:53.531107Z",
  "labels": "subaccount_id = 98cbd1c8-49e2-42d5-8266-980e3e8728a4"
}
		`), &obj)

		if assert.NoError(t, err) {
			uut, diags := subaccountServiceInstanceValueFrom(context.TODO(), obj)

			assert.False(t, diags.HasError())

			assert.Equal(t, "76765dca-6683-473a-8f42-809e33a2ea68", uut.Id.ValueString())
			assert.Equal(t, "default", uut.Name.ValueString())
			assert.Equal(t, "98cbd1c8-49e2-42d5-8266-980e3e8728a4", uut.SubaccountId.ValueString())
			assert.Equal(t, "", uut.Parameters.ValueString())
			assert.True(t, uut.Ready.ValueBool())
			assert.Equal(t, "dc91ce2b-55c5-42eb-ad92-70aee582d827", uut.ServicePlanId.ValueString())
			assert.Equal(t, "service-manager", uut.PlatformId.ValueString())
			assert.Equal(t, "", uut.ReferencedInstanceId.ValueString())
			assert.False(t, uut.Shared.ValueBool())
			assert.Equal(t, "{\"crm_customer_id\":\"\",\"global_account_id\":\"795b53bb-a3f0-4769-adf0-26173282a975\",\"instance_name\":\"default\",\"license_type\":\"DEVELOPER\",\"origin\":\"sapcp\",\"platform\":\"sapcp\",\"region\":\"cf-eu10\",\"subaccount_id\":\"98cbd1c8-49e2-42d5-8266-980e3e8728a4\",\"subdomain\":\"my-subdomain\",\"zone_id\":\"98cbd1c8-49e2-42d5-8266-980e3e8728a4\"}", uut.Context.String())
			assert.True(t, uut.Usable.ValueBool())
			assert.Equal(t, "succeeded", uut.State.ValueString())
			assert.Equal(t, "2022-09-29T21:37:52Z", uut.CreatedDate.ValueString())
			assert.Equal(t, "2022-09-29T21:37:53Z", uut.LastModified.ValueString())
			assert.Equal(t, "{\"subaccount_id\":[\"98cbd1c8-49e2-42d5-8266-980e3e8728a4\\\"\"]}", uut.Labels.String())
		}
	})
}
