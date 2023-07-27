package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/provisioning"
)

type subaccountEnvironmentInstanceType struct {
	SubaccountId    types.String   `tfsdk:"subaccount_id"`
	Id              types.String   `tfsdk:"id"`
	BrokerId        types.String   `tfsdk:"broker_id"`
	CreatedDate     types.String   `tfsdk:"created_date"`
	CustomLabels    types.Map      `tfsdk:"custom_labels"`
	DashboardUrl    types.String   `tfsdk:"dashboard_url"`
	Description     types.String   `tfsdk:"description"`
	EnvironmentType types.String   `tfsdk:"environment_type"`
	Labels          types.String   `tfsdk:"labels"`
	LandscapeLabel  types.String   `tfsdk:"landscape_label"`
	LastModified    types.String   `tfsdk:"last_modified"`
	Name            types.String   `tfsdk:"name"`
	Operation       types.String   `tfsdk:"operation"`
	Parameters      types.String   `tfsdk:"parameters"`
	PlanId          types.String   `tfsdk:"plan_id"`
	PlanName        types.String   `tfsdk:"plan_name"`
	PlatformId      types.String   `tfsdk:"platform_id"`
	ServiceId       types.String   `tfsdk:"service_id"`
	ServiceName     types.String   `tfsdk:"service_name"`
	State           types.String   `tfsdk:"state"`
	TenantId        types.String   `tfsdk:"tenant_id"`
	Type_           types.String   `tfsdk:"type"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}

func subaccountEnvironmentInstanceValueFrom(ctx context.Context, value provisioning.EnvironmentInstanceResponseObject) (subaccountEnvironmentInstanceType, diag.Diagnostics) {
	environmentInstance := subaccountEnvironmentInstanceType{
		Id:              types.StringValue(value.Id),
		BrokerId:        types.StringValue(value.BrokerId),
		CreatedDate:     timeToValue(value.CreatedDate.Time()),
		DashboardUrl:    types.StringValue(value.DashboardUrl),
		Description:     types.StringValue(value.Description),
		EnvironmentType: types.StringValue(value.EnvironmentType),
		Labels:          types.StringValue(value.Labels),
		LandscapeLabel:  types.StringValue(value.LandscapeLabel),
		LastModified:    timeToValue(value.ModifiedDate.Time()),
		Name:            types.StringValue(value.Name),
		Operation:       types.StringValue(value.Operation),
		Parameters:      types.StringValue(value.Parameters),
		PlanId:          types.StringValue(value.PlanId),
		PlanName:        types.StringValue(value.PlanName),
		PlatformId:      types.StringValue(value.PlatformId),
		ServiceId:       types.StringValue(value.ServiceId),
		ServiceName:     types.StringValue(value.ServiceName),
		SubaccountId:    types.StringValue(value.SubaccountGUID),
		State:           types.StringValue(value.State),
		TenantId:        types.StringValue(value.TenantId),
		Type_:           types.StringValue(value.Type_),
	}

	var diags, diagnostics diag.Diagnostics

	environmentInstance.CustomLabels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.CustomLabels)
	diagnostics.Append(diags...)

	return environmentInstance, diagnostics
}
