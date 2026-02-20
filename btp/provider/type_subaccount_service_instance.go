package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

type subaccountServiceInstanceType struct {
	SubaccountId         types.String         `tfsdk:"subaccount_id"`
	Id                   types.String         `tfsdk:"id"`
	Name                 types.String         `tfsdk:"name"`
	Parameters           jsontypes.Normalized `tfsdk:"parameters"`
	Ready                types.Bool           `tfsdk:"ready"`
	ServicePlanId        types.String         `tfsdk:"serviceplan_id"`
	PlatformId           types.String         `tfsdk:"platform_id"`
	ReferencedInstanceId types.String         `tfsdk:"referenced_instance_id"`
	Shared               types.Bool           `tfsdk:"shared"`
	Context              types.String         `tfsdk:"context"`
	Usable               types.Bool           `tfsdk:"usable"`
	State                types.String         `tfsdk:"state"`
	CreatedDate          types.String         `tfsdk:"created_date"`
	LastModified         types.String         `tfsdk:"last_modified"`
	Labels               types.Map            `tfsdk:"labels"`
	Timeouts             timeouts.Value       `tfsdk:"timeouts"`
	DashboardUrl         types.String         `tfsdk:"dashboard_url"`
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
		Context:              types.StringValue(string(value.Context)),
		Usable:               types.BoolValue(value.Usable),
		State:                types.StringValue(value.LastOperation.State),
		CreatedDate:          timeToValue(value.CreatedAt),
		LastModified:         timeToValue(value.UpdatedAt),
		DashboardUrl:         types.StringValue(value.DashboardUrl),
	}

	var diags, diagnostics diag.Diagnostics

	//Remove computed labels to avoid state inconsistencies
	value.Labels = tfutils.RemoveComputedlabels(value.Labels)

	serviceInstance.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	return serviceInstance, diagnostics
}

func subaccountServiceInstanceListValueFrom(ctx context.Context, value servicemanager.ServiceInstanceResponseObject) (subaccountServiceInstanceType, diag.Diagnostics) {
	timeoutAttrTypes := map[string]attr.Type{
		"create": types.StringType,
		"delete": types.StringType,
		"update": types.StringType,
	}

	serviceInstance := subaccountServiceInstanceType{
		SubaccountId:         types.StringValue(value.SubaccountId),
		Id:                   types.StringValue(value.Id),
		Ready:                types.BoolValue(value.Ready),
		Name:                 types.StringValue(value.Name),
		ServicePlanId:        types.StringValue(value.ServicePlanId),
		PlatformId:           types.StringValue(value.PlatformId),
		ReferencedInstanceId: types.StringValue(value.ReferencedInstanceId),
		Shared:               types.BoolValue(value.Shared),
		Context:              types.StringValue(string(value.Context)),
		Usable:               types.BoolValue(value.Usable),
		CreatedDate:          timeToValue(value.CreatedAt),
		LastModified:         timeToValue(value.UpdatedAt),
		DashboardUrl:         types.StringValue(value.DashboardUrl),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(timeoutAttrTypes),
		},
	}

	var diags, diagnostics diag.Diagnostics

	//Remove computed labels to avoid state inconsistencies
	value.Labels = tfutils.RemoveComputedlabels(value.Labels)

	serviceInstance.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	return serviceInstance, diagnostics
}

type subaccountServiceInstanceDataSourceType struct {
	SubaccountId         types.String `tfsdk:"subaccount_id"`
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Parameters           types.String `tfsdk:"parameters"`
	Ready                types.Bool   `tfsdk:"ready"`
	ServicePlanId        types.String `tfsdk:"serviceplan_id"`
	PlatformId           types.String `tfsdk:"platform_id"`
	ReferencedInstanceId types.String `tfsdk:"referenced_instance_id"`
	Shared               types.Bool   `tfsdk:"shared"`
	Context              types.String `tfsdk:"context"`
	Usable               types.Bool   `tfsdk:"usable"`
	State                types.String `tfsdk:"state"`
	CreatedDate          types.String `tfsdk:"created_date"`
	LastModified         types.String `tfsdk:"last_modified"`
	Labels               types.Map    `tfsdk:"labels"`
	DashboardUrl         types.String `tfsdk:"dashboard_url"`
}

func subaccountServiceInstanceDataSourceValueFrom(ctx context.Context, value servicemanager.ServiceInstanceResponseObject) (subaccountServiceInstanceDataSourceType, diag.Diagnostics) {
	serviceInstance := subaccountServiceInstanceDataSourceType{
		SubaccountId:         types.StringValue(value.SubaccountId),
		Id:                   types.StringValue(value.Id),
		Ready:                types.BoolValue(value.Ready),
		Name:                 types.StringValue(value.Name),
		ServicePlanId:        types.StringValue(value.ServicePlanId),
		PlatformId:           types.StringValue(value.PlatformId),
		ReferencedInstanceId: types.StringValue(value.ReferencedInstanceId),
		Shared:               types.BoolValue(value.Shared),
		Usable:               types.BoolValue(value.Usable),
		Parameters:           types.StringValue(value.Parameters),
		State:                types.StringValue(value.LastOperation.State),
		Context:              types.StringValue(string(value.Context)),
		CreatedDate:          timeToValue(value.CreatedAt),
		LastModified:         timeToValue(value.UpdatedAt),
		DashboardUrl:         types.StringValue(value.DashboardUrl),
	}

	var diags, diagnostics diag.Diagnostics

	//Remove computed labels to avoid state inconsistencies
	value.Labels = tfutils.RemoveComputedlabels(value.Labels)

	serviceInstance.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	return serviceInstance, diagnostics
}
