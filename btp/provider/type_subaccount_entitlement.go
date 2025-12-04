package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

// Only the category SERVICE and QUOTA_BASED_APPLICATION have a numeric quota (amount)
const entitlementCategoryService = "SERVICE"
const entitlementCategoryQuotaBasedApplication = "QUOTA_BASED_APPLICATION"
const entitlementCategoryPlatformBasedApplication = "PLATFORM"

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

	if isTransferAmountRequired(subaccountEntitlement.Category.ValueString()) {
		subaccountEntitlement.Amount = types.Int64Value(int64(value.Assignment.Amount))
	}

	subaccountEntitlement.State = types.StringValue(value.Assignment.EntityState)
	subaccountEntitlement.LastModified = timeToValue(value.Assignment.ModifiedDate.Time())
	subaccountEntitlement.CreatedDate = timeToValue(value.Assignment.CreatedDate.Time())

	return subaccountEntitlement, diag.Diagnostics{}
}

func isTransferAmountRequired(category string) bool {
	// Check if Amount needs to be mapped - only true if the entitlement has a numeric quota
	return category == entitlementCategoryService || category == entitlementCategoryQuotaBasedApplication || category == entitlementCategoryPlatformBasedApplication
}
