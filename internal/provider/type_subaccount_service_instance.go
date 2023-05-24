package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
)

type subaccountServiceInstanceType struct {
	SubaccountId         types.String `tfsdk:"subaccount_id"`
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Parameters           types.String `tfsdk:"parameters"`
	Ready                types.Bool   `tfsdk:"ready"`
	ServicePlanId        types.String `tfsdk:"serviceplan_id"`
	PlatformId           types.String `tfsdk:"platform_id"`
	ReferencedInstanceId types.String `tfsdk:"referenced_instance_id"`
	Shared               types.Bool   `tfsdk:"shared"`
	Context              types.Map    `tfsdk:"context"`
	Usable               types.Bool   `tfsdk:"usable"`
	State                types.String `tfsdk:"state"`
	CreatedDate          types.String `tfsdk:"created_date"`
	LastModified         types.String `tfsdk:"last_modified"`
	Labels               types.Map    `tfsdk:"labels"`
}

func subaccountServiceInstanceValueFrom(ctx context.Context, value servicemanager.ServiceInstanceResponseObject) (subaccountServiceInstanceType, diag.Diagnostics) {
	serviceInstance := subaccountServiceInstanceType{
		SubaccountId:         types.StringValue(value.SubaccountId),
		Id:                   types.StringValue(value.Id),
		Ready:                types.BoolValue(value.Ready),
		Name:                 types.StringValue(value.Name),
		ServicePlanId:        types.StringValue(value.ServicePlanId),
		PlatformId:           types.StringValue(value.PlatformId),
		ReferencedInstanceId: types.StringValue(value.ReferencedInstanceId),
		Shared:               types.BoolValue(value.Shared),
		Usable:               types.BoolValue(value.Usable),
		State:                types.StringValue(value.LastOperation.State),
		CreatedDate:          timeToValue(value.CreatedAt),
		LastModified:         timeToValue(value.UpdatedAt),
	}

	var diags, diagnostics diag.Diagnostics

	serviceInstance.Context, diags = types.MapValueFrom(ctx, types.StringType, value.Context)
	diagnostics.Append(diags...)

	serviceInstance.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	return serviceInstance, diagnostics
}
