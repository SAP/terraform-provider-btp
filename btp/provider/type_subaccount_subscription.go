package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/saas_manager_service"
)

type subaccountSubscriptionType struct {
	SubaccountId              types.String         `tfsdk:"subaccount_id"`
	Id                        types.String         `tfsdk:"id"`
	AppName                   types.String         `tfsdk:"app_name"`
	PlanName                  types.String         `tfsdk:"plan_name"`
	Parameters                jsontypes.Normalized `tfsdk:"parameters"`
	AdditionalPlanFeatures    types.Set            `tfsdk:"additional_plan_features"`
	AppId                     types.String         `tfsdk:"app_id"`
	AuthenticationProvider    types.String         `tfsdk:"authentication_provider"`
	Category                  types.String         `tfsdk:"category"`
	CommercialAppName         types.String         `tfsdk:"commercial_app_name"`
	CreatedDate               types.String         `tfsdk:"created_date"`
	CustomerDeveloped         types.Bool           `tfsdk:"customer_developed"`
	Description               types.String         `tfsdk:"description"`
	DisplayName               types.String         `tfsdk:"display_name"`
	FormationSolutionName     types.String         `tfsdk:"formation_solution_name"`
	GlobalAccountId           types.String         `tfsdk:"globalaccount_id"`
	Labels                    types.Map            `tfsdk:"labels"`
	LastModified              types.String         `tfsdk:"last_modified"`
	PlatformEntityId          types.String         `tfsdk:"platform_entity_id"`
	Quota                     types.Int64          `tfsdk:"quota"`
	State                     types.String         `tfsdk:"state"`
	SubscribedSubaccountId    types.String         `tfsdk:"subscribed_subaccount_id"`
	SubscribedTenantId        types.String         `tfsdk:"subscribed_tenant_id"`
	SubscriptionUrl           types.String         `tfsdk:"subscription_url"`
	SupportsParametersUpdates types.Bool           `tfsdk:"supports_parameters_updates"`
	SupportsPlanUpdates       types.Bool           `tfsdk:"supports_plan_updates"`
	TenantId                  types.String         `tfsdk:"tenant_id"`
	Timeouts                  timeouts.Value       `tfsdk:"timeouts"`
}

func subaccountSubscriptionValueFrom(ctx context.Context, value saas_manager_service.EntitledApplicationsResponseObject) (subaccountSubscriptionType, diag.Diagnostics) {
	subscription := subaccountSubscriptionType{
		SubaccountId:              types.StringValue(value.SubscribedSubaccountId),
		Id:                        types.StringValue(value.SubscriptionGUID),
		AppId:                     types.StringValue(value.AppId),
		AppName:                   types.StringValue(value.AppName),
		AuthenticationProvider:    types.StringValue(value.AuthenticationProvider),
		Category:                  types.StringValue(value.Category),
		CommercialAppName:         types.StringValue(value.CommercialAppName),
		CreatedDate:               timeToValue(value.CreatedDate.Time()),
		CustomerDeveloped:         types.BoolValue(value.CustomerDeveloped),
		Description:               types.StringValue(value.Description),
		DisplayName:               types.StringValue(value.DisplayName),
		FormationSolutionName:     types.StringValue(value.FormationSolutionName),
		GlobalAccountId:           types.StringValue(value.GlobalAccountId),
		LastModified:              timeToValue(value.ModifiedDate.Time()),
		Parameters:                jsontypes.NewNormalizedNull(),
		PlanName:                  types.StringValue(value.PlanName),
		PlatformEntityId:          types.StringValue(value.PlatformEntityId),
		Quota:                     types.Int64Value(int64(value.Quota)),
		State:                     types.StringValue(value.State),
		SubscribedSubaccountId:    types.StringValue(value.SubscribedSubaccountId),
		SubscribedTenantId:        types.StringValue(value.SubscribedTenantId),
		SubscriptionUrl:           types.StringValue(value.SubscriptionUrl),
		SupportsParametersUpdates: types.BoolValue(value.SupportsParametersUpdates),
		SupportsPlanUpdates:       types.BoolValue(value.SupportsPlanUpdates),
		TenantId:                  types.StringValue(value.TenantId),
	}

	var diags, diagnostics diag.Diagnostics

	subscription.AdditionalPlanFeatures, diags = types.SetValueFrom(ctx, types.StringType, value.AdditionalPlanFeatures)
	diagnostics.Append(diags...)

	subscription.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	return subscription, diagnostics
}

type subaccountSubscriptionDataSourceType struct {
	SubaccountId              types.String `tfsdk:"subaccount_id"`
	Id                        types.String `tfsdk:"id"`
	AppName                   types.String `tfsdk:"app_name"`
	PlanName                  types.String `tfsdk:"plan_name"`
	Parameters                types.String `tfsdk:"parameters"`
	AdditionalPlanFeatures    types.Set    `tfsdk:"additional_plan_features"`
	AppId                     types.String `tfsdk:"app_id"`
	AuthenticationProvider    types.String `tfsdk:"authentication_provider"`
	Category                  types.String `tfsdk:"category"`
	CommercialAppName         types.String `tfsdk:"commercial_app_name"`
	CreatedDate               types.String `tfsdk:"created_date"`
	CustomerDeveloped         types.Bool   `tfsdk:"customer_developed"`
	Description               types.String `tfsdk:"description"`
	DisplayName               types.String `tfsdk:"display_name"`
	FormationSolutionName     types.String `tfsdk:"formation_solution_name"`
	GlobalAccountId           types.String `tfsdk:"globalaccount_id"`
	Labels                    types.Map    `tfsdk:"labels"`
	LastModified              types.String `tfsdk:"last_modified"`
	PlatformEntityId          types.String `tfsdk:"platform_entity_id"`
	Quota                     types.Int64  `tfsdk:"quota"`
	State                     types.String `tfsdk:"state"`
	SubscribedSubaccountId    types.String `tfsdk:"subscribed_subaccount_id"`
	SubscribedTenantId        types.String `tfsdk:"subscribed_tenant_id"`
	SubscriptionUrl           types.String `tfsdk:"subscription_url"`
	SupportsParametersUpdates types.Bool   `tfsdk:"supports_parameters_updates"`
	SupportsPlanUpdates       types.Bool   `tfsdk:"supports_plan_updates"`
	TenantId                  types.String `tfsdk:"tenant_id"`
}

func subaccountSubscriptionDataSourceValueFrom(ctx context.Context, value saas_manager_service.EntitledApplicationsResponseObject) (subaccountSubscriptionDataSourceType, diag.Diagnostics) {
	subscription := subaccountSubscriptionDataSourceType{
		SubaccountId:              types.StringValue(value.SubscribedSubaccountId),
		Id:                        types.StringValue(value.SubscriptionGUID),
		AppId:                     types.StringValue(value.AppId),
		AppName:                   types.StringValue(value.AppName),
		AuthenticationProvider:    types.StringValue(value.AuthenticationProvider),
		Category:                  types.StringValue(value.Category),
		CommercialAppName:         types.StringValue(value.CommercialAppName),
		CreatedDate:               timeToValue(value.CreatedDate.Time()),
		CustomerDeveloped:         types.BoolValue(value.CustomerDeveloped),
		Description:               types.StringValue(value.Description),
		DisplayName:               types.StringValue(value.DisplayName),
		FormationSolutionName:     types.StringValue(value.FormationSolutionName),
		GlobalAccountId:           types.StringValue(value.GlobalAccountId),
		LastModified:              timeToValue(value.ModifiedDate.Time()),
		Parameters:                types.StringNull(),
		PlanName:                  types.StringValue(value.PlanName),
		PlatformEntityId:          types.StringValue(value.PlatformEntityId),
		Quota:                     types.Int64Value(int64(value.Quota)),
		State:                     types.StringValue(value.State),
		SubscribedSubaccountId:    types.StringValue(value.SubscribedSubaccountId),
		SubscribedTenantId:        types.StringValue(value.SubscribedTenantId),
		SubscriptionUrl:           types.StringValue(value.SubscriptionUrl),
		SupportsParametersUpdates: types.BoolValue(value.SupportsParametersUpdates),
		SupportsPlanUpdates:       types.BoolValue(value.SupportsPlanUpdates),
		TenantId:                  types.StringValue(value.TenantId),
	}

	var diags, diagnostics diag.Diagnostics

	subscription.AdditionalPlanFeatures, diags = types.SetValueFrom(ctx, types.StringType, value.AdditionalPlanFeatures)
	diagnostics.Append(diags...)

	subscription.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, value.Labels)
	diagnostics.Append(diags...)

	return subscription, diagnostics
}
