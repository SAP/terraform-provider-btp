package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

type directoryEntitlementType struct {
	DirectoryId          types.String `tfsdk:"directory_id"`
	Id                   types.String `tfsdk:"id"`
	ServiceName          types.String `tfsdk:"service_name"`
	PlanName             types.String `tfsdk:"plan_name"`
	Amount               types.Int64  `tfsdk:"amount"`
	AutoAssign           types.Bool   `tfsdk:"auto_assign"`
	AutoDistributeAmount types.Int64  `tfsdk:"auto_distribute_amount"`
	Distribute           types.Bool   `tfsdk:"distribute"`
	Category             types.String `tfsdk:"category"`
	PlanId               types.String `tfsdk:"plan_id"`
	State                types.String `tfsdk:"state"`
	CreatedDate          types.String `tfsdk:"created_date"`
	LastModified         types.String `tfsdk:"last_modified"`
}

func directoryEntitlementValueFrom(ctx context.Context, value btpcli.UnfoldedEntitlement) (directoryEntitlementType, diag.Diagnostics) {
	return directoryEntitlementType{
		DirectoryId:          types.StringValue(value.Assignment.EntityId),
		Id:                   types.StringValue(value.Plan.UniqueIdentifier),
		ServiceName:          types.StringValue(value.Service.Name),
		PlanName:             types.StringValue(value.Plan.Name),
		Category:             types.StringValue(value.Plan.Category),
		PlanId:               types.StringValue(value.Plan.UniqueIdentifier),
		Amount:               types.Int64Value(int64(value.Assignment.Amount)),
		AutoAssign:           types.BoolValue(value.Assignment.AutoAssign),
		AutoDistributeAmount: types.Int64Value(int64(value.Assignment.AutoDistributeAmount)),
		State:                types.StringValue(value.Assignment.EntityState),
		LastModified:         timeToValue(value.Assignment.ModifiedDate.Time()),
		CreatedDate:          timeToValue(value.Assignment.CreatedDate.Time()),
	}, diag.Diagnostics{}
}
