package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
)

func TestSubaccountServiceBindingValueFrom(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var obj servicemanager.ServiceBindingResponseObject
		err := json.Unmarshal([]byte(`
{
  "id": "72b34886-603e-4616-a868-05254f2cacc4",
  "ready": true,
  "last_operation": {
    "id": "c9ff2cb9-031b-422f-9475-af43f155b779",
    "ready": true,
    "type": "create",
    "state": "succeeded",
    "resource_id": "72b34886-603e-4616-a868-05254f2cacc4",
    "resource_type": "/v1/service_bindings",
    "platform_id": "service-manager",
    "correlation_id": "0f03446e-8ea2-45d9-846f-aa73051fc9b4",
    "reschedule": false,
    "reschedule_timestamp": "0001-01-01T00:00:00Z",
    "deletion_scheduled": "0001-01-01T00:00:00Z",
    "created_at": "2023-02-20T10:45:48.064104Z",
    "updated_at": "2023-02-20T10:45:48.172656Z"
  },
  "name": "test",
  "service_instance_id": "8911491d-0e1d-425d-a233-785512602d6f",
  "context": {
    "crm_customer_id": "",
    "global_account_id": "795b53bb-a3f0-4769-adf0-26173282a975",
    "instance_name": "malware-scanner",
    "license_type": "DEVELOPER",
    "origin": "sapcp",
    "platform": "sapcp",
    "region": "cf-eu30",
    "service_instance_id": "8911491d-0e1d-425d-a233-785512602d6f",
    "subaccount_id": "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f",
    "subdomain": "malware-scan-1qlhbviw",
    "zone_id": "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
  },
  "credentials": {
    "password": "a-password",
    "sync_scan_url": "https://malware-scanner.cf.eu30.hana.ondemand.com",
    "uri": "malware-scanner.cf.eu30.hana.ondemand.com",
    "url": "https://malware-scanner.cf.eu30.hana.ondemand.com",
    "username": "a-username"
  },
  "subaccount_id": "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f",
  "created_at": "2023-02-20T10:45:48.064099Z",
  "updated_at": "2023-02-20T10:45:48.165746Z",
  "labels": "subaccount_id = 6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
}
		`), &obj)

		if assert.NoError(t, err) {
			uut, diags := subaccountServiceBindingValueFrom(context.TODO(), obj)

			assert.False(t, diags.HasError())
			assert.Equal(t, "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f", uut.SubaccountId.ValueString())
			assert.Equal(t, "8911491d-0e1d-425d-a233-785512602d6f", uut.ServiceInstanceId.ValueString())
			assert.Equal(t, "test", uut.Name.ValueString())
			assert.Equal(t, "", uut.Parameters.ValueString())
			assert.Equal(t, "72b34886-603e-4616-a868-05254f2cacc4", uut.Id.ValueString())
			assert.True(t, uut.Ready.ValueBool())
			assert.Equal(t, "{\"crm_customer_id\":\"\",\"global_account_id\":\"795b53bb-a3f0-4769-adf0-26173282a975\",\"instance_name\":\"malware-scanner\",\"license_type\":\"DEVELOPER\",\"origin\":\"sapcp\",\"platform\":\"sapcp\",\"region\":\"cf-eu30\",\"service_instance_id\":\"8911491d-0e1d-425d-a233-785512602d6f\",\"subaccount_id\":\"6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f\",\"subdomain\":\"malware-scan-1qlhbviw\",\"zone_id\":\"6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f\"}", uut.Context.String())
			assert.Equal(t, "<null>", uut.BindResource.String())
			assert.Equal(t, "{\n    \"password\": \"a-password\",\n    \"sync_scan_url\": \"https://malware-scanner.cf.eu30.hana.ondemand.com\",\n    \"uri\": \"malware-scanner.cf.eu30.hana.ondemand.com\",\n    \"url\": \"https://malware-scanner.cf.eu30.hana.ondemand.com\",\n    \"username\": \"a-username\"\n  }", uut.Credentials.ValueString())
			assert.Equal(t, "succeeded", uut.State.ValueString())
			assert.Equal(t, "2023-02-20T10:45:48Z", uut.LastModified.ValueString())
			assert.Equal(t, "2023-02-20T10:45:48Z", uut.CreatedDate.ValueString())
			assert.Equal(t, "{\"subaccount_id\":[\"6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f\\\"\"]}", uut.Labels.String())
		}
	})
}
