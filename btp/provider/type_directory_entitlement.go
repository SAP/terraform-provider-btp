package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis_entitlements"
)

type directoryEntitlementType struct {
	DirectoryId          types.String `tfsdk:"directory_id"`
	Id                   types.String `tfsdk:"id"`
	ServiceName          types.String `tfsdk:"service_name"`
	PlanName             types.String `tfsdk:"plan_name"`
	PlanUniqueIdentifier types.String `tfsdk:"plan_unique_identifier"`
	Amount               types.Int64  `tfsdk:"amount"`
	AutoAssign           types.Bool   `tfsdk:"auto_assign"`
	AutoDistributeAmount types.Int64  `tfsdk:"auto_distribute_amount"`
	Distribute           types.Bool   `tfsdk:"distribute"`
	Category             types.String `tfsdk:"category"`
	PlanId               types.String `tfsdk:"plan_id"`
}

func directoryEntitlementValueFrom(ctx context.Context, value btpcli.UnfoldedEntitlement, directoryId string, distribute bool) (directoryEntitlementType, diag.Diagnostics) {
	var directoryEntitlement directoryEntitlementType

	directoryEntitlement.DirectoryId = types.StringValue(directoryId)
	directoryEntitlement.Id = types.StringValue(value.Plan.UniqueIdentifier)
	directoryEntitlement.ServiceName = types.StringValue(value.Service.Name)
	directoryEntitlement.PlanName = types.StringValue(value.Plan.Name)
	directoryEntitlement.Category = types.StringValue(value.Plan.Category)
	directoryEntitlement.PlanId = types.StringValue(value.Plan.UniqueIdentifier)
	directoryEntitlement.PlanUniqueIdentifier = types.StringValue(value.Plan.UniqueIdentifier)

	if isTransferAmountRequired(directoryEntitlement.Category.ValueString()) {
		directoryEntitlement.Amount = types.Int64Value(int64(value.Plan.Amount))
	}

	directoryEntitlement.AutoAssign = types.BoolValue(value.Plan.AutoAssign)
	directoryEntitlement.AutoDistributeAmount = types.Int64Value(int64(value.Plan.AutoDistributeAmount))
	directoryEntitlement.Distribute = types.BoolValue(distribute)

	return directoryEntitlement, diag.Diagnostics{}
}

func directoryEntitlementListValueFrom(service cis_entitlements.EntitledServicesResponseObject, servicePlan cis_entitlements.ServicePlanResponseObject, directoryId string) *directoryEntitlementType {
	resDm := &directoryEntitlementType{
		DirectoryId:          types.StringValue(directoryId),
		Id:                   types.StringValue(servicePlan.UniqueIdentifier),
		ServiceName:          types.StringValue(service.Name),
		PlanName:             types.StringValue(servicePlan.Name),
		PlanUniqueIdentifier: types.StringValue(servicePlan.UniqueIdentifier),
		AutoAssign:           types.BoolValue(servicePlan.AutoAssign),
		AutoDistributeAmount: types.Int64Value(int64(servicePlan.AutoDistributeAmount)),
		Category:             types.StringValue(servicePlan.Category),
		PlanId:               types.StringValue(servicePlan.UniqueIdentifier),
	}

	if isTransferAmountRequired(resDm.Category.ValueString()) {
		resDm.Amount = types.Int64Value(int64(servicePlan.Amount))
	}

	return resDm
}
