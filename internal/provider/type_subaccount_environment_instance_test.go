package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/provisioning"
)

func TestSubaccountEnvironmentInstanceFromValue(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var obj provisioning.EnvironmentInstanceResponseObject
		err := json.Unmarshal([]byte(`
{
  "id": "F03C6896-6BB5-44E4-8EFF-E0664150FB96",
  "name": "my-environment",
  "brokerId": "1749E930-2380-4725-BAC0-BE9EA15CDBBF",
  "globalAccountGUID": "795b53bb-a3f0-4769-adf0-26173282a975",
  "subaccountGUID": "917e3793-5ba7-4b20-860d-11102aec16ff",
  "tenantId": "917e3793-5ba7-4b20-860d-11102aec16ff",
  "serviceId": "fa31b750-375f-4268-bee1-604811a89fd9",
  "planId": "fc5abe63-2a7d-4848-babf-f63a5d316df1",
  "operation": "provision",
  "parameters": "{\"instance_name\":\"my-unique-org-name\",\"users\":[{\"email\":\"john.doe@mycompany.com\",\"id\":\"john.doe@mycompany.com\"}],\"status\":\"ACTIVE\"}",
  "labels": "{\"API Endpoint\":\"https://api.cf.eu30.hana.ondemand.com\",\"Org Name\":\"my-unique-org-name\",\"Org ID\":\"9c92a592-0544-46a5-9f3a-90c4b59b8c17\"}",
  "customLabels": {},
  "type": "Provision",
  "status": "Processed",
  "environmentType": "cloudfoundry",
  "landscapeLabel": "cf-eu30",
  "platformId": "9c92a592-0544-46a5-9f3a-90c4b59b8c17",
  "createdDate": 1676039738974,
  "modifiedDate": 1676039756332,
  "state": "OK",
  "serviceName": "cloudfoundry",
  "planName": "standard",
  "description": "my-description",
  "dashboardUrl": "http://my.dashboard/"
}
		`), &obj)

		if assert.NoError(t, err) {
			uut, diags := subaccountEnvironmentInstanceValueFrom(context.TODO(), obj)

			assert.False(t, diags.HasError())

			assert.Equal(t, "917e3793-5ba7-4b20-860d-11102aec16ff", uut.SubaccountId.ValueString())
			assert.Equal(t, "F03C6896-6BB5-44E4-8EFF-E0664150FB96", uut.Id.ValueString())
			assert.Equal(t, "1749E930-2380-4725-BAC0-BE9EA15CDBBF", uut.BrokerId.ValueString())
			assert.Equal(t, "2023-02-10T14:35:38Z", uut.CreatedDate.ValueString())
			assert.Equal(t, "{}", uut.CustomLabels.String())
			assert.Equal(t, "http://my.dashboard/", uut.DashboardUrl.ValueString())
			assert.Equal(t, "my-description", uut.Description.ValueString())
			assert.Equal(t, "cloudfoundry", uut.EnvironmentType.ValueString())
			assert.Equal(t, "{\"API Endpoint\":\"https://api.cf.eu30.hana.ondemand.com\",\"Org Name\":\"my-unique-org-name\",\"Org ID\":\"9c92a592-0544-46a5-9f3a-90c4b59b8c17\"}", uut.Labels.ValueString())
			assert.Equal(t, "cf-eu30", uut.LandscapeLabel.ValueString())
			assert.Equal(t, "2023-02-10T14:35:56Z", uut.LastModified.ValueString())
			assert.Equal(t, "my-environment", uut.Name.ValueString())
			assert.Equal(t, "provision", uut.Operation.ValueString())
			assert.Equal(t, "{\"instance_name\":\"my-unique-org-name\",\"users\":[{\"email\":\"john.doe@mycompany.com\",\"id\":\"john.doe@mycompany.com\"}],\"status\":\"ACTIVE\"}", uut.Parameters.ValueString())
			assert.Equal(t, "fc5abe63-2a7d-4848-babf-f63a5d316df1", uut.PlanId.ValueString())
			assert.Equal(t, "standard", uut.PlanName.ValueString())
			assert.Equal(t, "9c92a592-0544-46a5-9f3a-90c4b59b8c17", uut.PlatformId.ValueString())
			assert.Equal(t, "fa31b750-375f-4268-bee1-604811a89fd9", uut.ServiceId.ValueString())
			assert.Equal(t, "cloudfoundry", uut.ServiceName.ValueString())
			assert.Equal(t, "OK", uut.State.ValueString())
			assert.Equal(t, "917e3793-5ba7-4b20-860d-11102aec16ff", uut.TenantId.ValueString())
			assert.Equal(t, "Provision", uut.Type_.ValueString())
		}
	})
}
