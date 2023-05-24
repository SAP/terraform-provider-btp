package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis_entitlements"
)

func TestSubaccountEntitlementFromValue(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var obj cis_entitlements.EntitledAndAssignedServicesResponseObject
		err := json.Unmarshal([]byte(`
{
  "assignedServices": [
    {
      "name": "alert-notification",
      "displayName": "Alert Notification",
      "businessCategory": {
        "id": "FOUNDATION_CROSS_SERVICES",
        "displayName": "Foundation / Cross Services"
      },
      "servicePlans": [
        {
          "name": "free",
          "displayName": "free",
          "uniqueIdentifier": "alert-notification-free",
          "category": "ELASTIC_SERVICE",
          "beta": false,
          "maxAllowedSubaccountQuota": null,
          "unlimited": false,
          "assignmentInfo": [
            {
              "entityId": "d162b191-594c-4184-bce5-a24d3fbc0818",
              "entityType": "SUBACCOUNT",
              "amount": 10,
              "requestedAmount": null,
              "entityState": "OK",
              "autoAssign": false,
              "autoDistributeAmount": null,
              "createdDate": 1676299127572,
              "modifiedDate": 1676585526250,
              "resources": [],
              "unlimitedAmountAssigned": false,
              "parentId": "795b53bb-a3f0-4769-adf0-26173282a975",
              "parentType": "GLOBAL_ACCOUNT",
              "parentRemainingAmount": 10,
              "parentAmount": 10,
              "autoAssigned": false
            }
          ]
        }
      ]
    }
  ]
}
		`), &obj)

		if assert.NoError(t, err) {
			uut, diags := subaccountEntitlementValueFrom(context.TODO(), btpcli.UnfoldedEntitlement{
				Service:    obj.AssignedServices[0],
				Plan:       obj.AssignedServices[0].ServicePlans[0],
				Assignment: obj.AssignedServices[0].ServicePlans[0].AssignmentInfo[0],
			})

			assert.False(t, diags.HasError())

			assert.Equal(t, "d162b191-594c-4184-bce5-a24d3fbc0818", uut.SubaccountId.ValueString())
			assert.Equal(t, "alert-notification-free", uut.Id.ValueString())
			assert.Equal(t, "alert-notification", uut.ServiceName.ValueString())
			assert.Equal(t, "free", uut.PlanName.ValueString())
			assert.Equal(t, "alert-notification-free", uut.PlanId.ValueString())
			assert.Equal(t, int64(10), uut.Amount.ValueInt64())
			assert.Equal(t, "OK", uut.State.ValueString())
			assert.Equal(t, "2023-02-16T22:12:06Z", uut.LastModified.ValueString())
			assert.Equal(t, "2023-02-13T14:38:47Z", uut.CreatedDate.ValueString())
		}
	})
}
