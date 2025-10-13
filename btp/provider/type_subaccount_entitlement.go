package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

const subaccountEntitlementCategoryElasticService = "ELASTIC_SERVICE"
const subaccountEntitlementCategoryApplication = "APPLICATION"

type subaccountEntitlementType struct {
	SubaccountId         types.String `tfsdk:"subaccount_id"`
	Id                   types.String `tfsdk:"id"`
	ServiceName          types.String `tfsdk:"service_name"`
	PlanName             types.String `tfsdk:"plan_name"`
	Category             types.String `tfsdk:"category"`
	PlanId               types.String `tfsdk:"plan_id"`
	Amount               types.Int64  `tfsdk:"amount"`
	PlanUniqueIdentifier types.String `tfsdk:"plan_unique_identifier"`
	State                types.String `tfsdk:"state"`
	CreatedDate          types.String `tfsdk:"created_date"`
	LastModified         types.String `tfsdk:"last_modified"`
}

func subaccountEntitlementValueFrom(ctx context.Context, value btpcli.UnfoldedAssignment) (subaccountEntitlementType, diag.Diagnostics) {
	var subaccountEntitlement subaccountEntitlementType

	subaccountEntitlement.SubaccountId = types.StringValue(value.Assignment.EntityId)
	subaccountEntitlement.Id = types.StringValue(value.Plan.UniqueIdentifier)
	subaccountEntitlement.ServiceName = types.StringValue(value.Service.Name)
	subaccountEntitlement.PlanName = types.StringValue(value.Plan.Name)
	subaccountEntitlement.Category = types.StringValue(value.Plan.Category)
	subaccountEntitlement.PlanId = types.StringValue(value.Plan.UniqueIdentifier)
	subaccountEntitlement.PlanUniqueIdentifier = types.StringValue(value.Plan.UniqueIdentifier)

	if subaccountEntitlement.Category != types.StringValue(subaccountEntitlementCategoryElasticService) && subaccountEntitlement.Category != types.StringValue(subaccountEntitlementCategoryApplication) {
		// Transfer Amount only if the entitlement has a numeric quota
		subaccountEntitlement.Amount = types.Int64Value(int64(value.Assignment.Amount))
	}
	subaccountEntitlement.State = types.StringValue(value.Assignment.EntityState)
	subaccountEntitlement.LastModified = timeToValue(value.Assignment.ModifiedDate.Time())
	subaccountEntitlement.CreatedDate = timeToValue(value.Assignment.CreatedDate.Time())

	return subaccountEntitlement, diag.Diagnostics{}
}
