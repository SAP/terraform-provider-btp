package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

type subaccountEntitlementType struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	ServiceName  types.String `tfsdk:"service_name"`
	PlanName     types.String `tfsdk:"plan_name"`
	PlanId       types.String `tfsdk:"plan_id"`
	Amount       types.Int64  `tfsdk:"amount"`
	State        types.String `tfsdk:"state"`
	CreatedDate  types.String `tfsdk:"created_date"`
	LastModified types.String `tfsdk:"last_modified"`
}

func subaccountEntitlementValueFrom(ctx context.Context, value btpcli.UnfoldedEntitlement) (subaccountEntitlementType, diag.Diagnostics) {
	return subaccountEntitlementType{
		SubaccountId: types.StringValue(value.Assignment.EntityId),
		Id:           types.StringValue(value.Plan.UniqueIdentifier),
		ServiceName:  types.StringValue(value.Service.Name),
		PlanName:     types.StringValue(value.Plan.Name),
		PlanId:       types.StringValue(value.Plan.UniqueIdentifier),
		Amount:       types.Int64Value(int64(value.Assignment.Amount)),
		State:        types.StringValue(value.Assignment.EntityState),
		LastModified: timeToValue(value.Assignment.ModifiedDate.Time()),
		CreatedDate:  timeToValue(value.Assignment.CreatedDate.Time()),
	}, diag.Diagnostics{}
}
