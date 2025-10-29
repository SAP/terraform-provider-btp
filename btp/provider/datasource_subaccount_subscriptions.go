package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountSubscriptionsDataSource() datasource.DataSource {
	return &subaccountSubscriptionsDataSource{}
}

type subaccountSubscriptionsValue struct {
	AdditionalPlanFeatures    types.Set    `tfsdk:"additional_plan_features"`
	AppId                     types.String `tfsdk:"app_id"`
	AppName                   types.String `tfsdk:"app_name"`
	AuthenticationProvider    types.String `tfsdk:"authentication_provider"`
	AutomationState           types.String `tfsdk:"automation_state"`
	AutomationStateMessage    types.String `tfsdk:"automation_state_message"`
	Category                  types.String `tfsdk:"category"`
	CategoryDisplayName       types.String `tfsdk:"category_display_name"`
	CommercialAppName         types.String `tfsdk:"commercial_app_name"`
	CreatedDate               types.String `tfsdk:"created_date"`
	CustomerDeveloped         types.Bool   `tfsdk:"customer_developed"`
	Description               types.String `tfsdk:"description"`
	DisplayName               types.String `tfsdk:"display_name"`
	FormationSolutionName     types.String `tfsdk:"formation_solution_name"`
	GlobalAccountId           types.String `tfsdk:"globalaccount_id"`
	IncidentTrackingComponent types.String `tfsdk:"incident_tracking_component"`
	Labels                    types.Map    `tfsdk:"labels"`
	LastModified              types.String `tfsdk:"last_modified"`
	PlanDescription           types.String `tfsdk:"plan_description"`
	PlanName                  types.String `tfsdk:"plan_name"`
	PlatformEntityId          types.String `tfsdk:"platform_entity_id"`
	Quota                     types.Int64  `tfsdk:"quota"`
	ShortDescription          types.String `tfsdk:"short_description"`
	State                     types.String `tfsdk:"state"`
	SubscribedSubaccountId    types.String `tfsdk:"subscribed_subaccount_id"`
	SubscribedTenantId        types.String `tfsdk:"subscribed_tenant_id"`
	Id                        types.String `tfsdk:"id"`
	SubscriptionUrl           types.String `tfsdk:"subscription_url"`
	SupportsParametersUpdates types.Bool   `tfsdk:"supports_parameters_updates"`
	SupportsPlanUpdates       types.Bool   `tfsdk:"supports_plan_updates"`
	TenantId                  types.String `tfsdk:"tenant_id"`
}

type subaccountSubscriptionsDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	/* OUTPUT */
	Values []subaccountSubscriptionsValue `tfsdk:"values"`
}

type subaccountSubscriptionsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountSubscriptionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_subscriptions", req.ProviderTypeName)
}

func (ds *subaccountSubscriptionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountSubscriptionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists all multitenant applications to which the subaccount is entitled to subscribe, including their subscription details.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` attribute instead",
				MarkdownDescription: "The ID of the subaccount.",
				Computed:            true,
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"app_name": schema.StringAttribute{
							MarkdownDescription: "The unique registration name of the deployed multitenant application as defined by the app developer.",
							Computed:            true,
						},
						"plan_name": schema.StringAttribute{
							MarkdownDescription: "The plan name of the application to which the consumer has subscribed.",
							Computed:            true,
						},
						"additional_plan_features": schema.SetAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "The list of features specific to this plan.",
							Computed:            true,
						},
						"app_id": schema.StringAttribute{
							MarkdownDescription: "The ID returned by XSUAA after the app provider has performed a bind of the multitenant application to a XSUAA service instance.",
							Computed:            true,
						},
						"authentication_provider": schema.StringAttribute{
							MarkdownDescription: "The authentication provider of the multitenant application. * XSUAA is the SAP Authorization and Trust Management service that defines scopes and permissions for users as tenants at the global account level. * IAS is Identity Authentication Service that defines scopes and permissions for users in zones (common data isolation systems across systems, SaaS tenants, and services).",
							Computed:            true,
						},
						"automation_state": schema.StringAttribute{
							MarkdownDescription: "The state of the automation solution.",
							Computed:            true,
						},
						"automation_state_message": schema.StringAttribute{
							MarkdownDescription: "The message that describes the automation solution state.",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							MarkdownDescription: "The technical name of the category defined by the app developer to which the multitenant application is grouped in customer-facing UIs.",
							Computed:            true,
						},
						"category_display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the category for customer-facing UIs.",
							Computed:            true,
						},
						"commercial_app_name": schema.StringAttribute{
							MarkdownDescription: "The commercial name of the deployed multitenant application as defined by the app developer.",
							Computed:            true,
						},
						"created_date": schema.StringAttribute{
							MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
							Computed:            true,
						},
						"customer_developed": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the application was developed by a customer. If not, then the application is developed by the cloud operator, such as SAP.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the multitenant application.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the application for customer-facing UIs.",
							Computed:            true,
						},
						"formation_solution_name": schema.StringAttribute{
							MarkdownDescription: "The name of the formations solution associated with the multitenant application.",
							Computed:            true,
						},
						"globalaccount_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the associated global account.",
							Computed:            true,
						},
						"incident_tracking_component": schema.StringAttribute{
							MarkdownDescription: "The application's incident-tracking component provided in metadata for customer-facing UIs.",
							Computed:            true,
						},
						"last_modified": schema.StringAttribute{
							MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
							Computed:            true,
						},
						"plan_description": schema.StringAttribute{
							MarkdownDescription: "The description of the plan for customer-facing UIs.",
							Computed:            true,
						},
						"platform_entity_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the landscape-specific environment.",
							Computed:            true,
						},
						"quota": schema.Int64Attribute{
							MarkdownDescription: "The total amount the subscribed subaccount is entitled to consume.",
							Computed:            true,
						},
						"short_description": schema.StringAttribute{
							MarkdownDescription: "The short description of the multitenant application for customer-facing UIs.",
							Computed:            true,
						},
						"state": schema.StringAttribute{
							MarkdownDescription: "The subscription state of the subaccount regarding the multitenant application.",
							Computed:            true,
						},
						"subscribed_subaccount_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the subaccount which is subscribed to the multitenant application.",
							Computed:            true,
						},
						"subscribed_tenant_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the tenant which is subscribed to a multitenant application.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "The technical ID generated by XSUAA for a multitenant application when a consumer subscribes to the application.",
							Computed:            true,
						},
						"subscription_url": schema.StringAttribute{
							MarkdownDescription: "The URL for app users to launch the subscribed application.",
							Computed:            true,
						},
						"supports_parameters_updates": schema.BoolAttribute{
							MarkdownDescription: "Specifies whether a consumer, whose subaccount is subscribed to the application, can change its subscriptions parameters.",
							Computed:            true,
						},
						"supports_plan_updates": schema.BoolAttribute{
							MarkdownDescription: "Specifies whether a consumer, whose subaccount is subscribed to the application, can change the subscription to a different plan that is available for this application and subaccount.",
							Computed:            true,
						},
						"tenant_id": schema.StringAttribute{
							MarkdownDescription: "The tenant ID of the application provider.",
							Computed:            true,
						},
						"labels": schema.MapAttribute{
							ElementType: types.SetType{
								ElemType: types.StringType,
							},
							MarkdownDescription: "The set of words or phrases assigned to the multitenant application subscription.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountSubscriptionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountSubscriptionsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.Subscription.List(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Subscriptions (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.SubaccountId
	data.Values = []subaccountSubscriptionsValue{}

	for _, subscription := range cliRes {
		value := subaccountSubscriptionsValue{
			AppId:                     types.StringValue(subscription.AppId),
			AppName:                   types.StringValue(subscription.AppName),
			AuthenticationProvider:    types.StringValue(subscription.AuthenticationProvider),
			Category:                  types.StringValue(subscription.Category),
			CommercialAppName:         types.StringValue(subscription.CommercialAppName),
			CreatedDate:               timeToValue(subscription.CreatedDate.Time()),
			CustomerDeveloped:         types.BoolValue(subscription.CustomerDeveloped),
			Description:               types.StringValue(subscription.Description),
			DisplayName:               types.StringValue(subscription.DisplayName),
			FormationSolutionName:     types.StringValue(subscription.FormationSolutionName),
			GlobalAccountId:           types.StringValue(subscription.GlobalAccountId),
			LastModified:              timeToValue(subscription.ModifiedDate.Time()),
			PlanName:                  types.StringValue(subscription.PlanName),
			PlatformEntityId:          types.StringValue(subscription.PlatformEntityId),
			Quota:                     types.Int64Value(int64(subscription.Quota)),
			State:                     types.StringValue(subscription.State),
			SubscribedSubaccountId:    types.StringValue(subscription.SubscribedSubaccountId),
			SubscribedTenantId:        types.StringValue(subscription.SubscribedTenantId),
			Id:                        types.StringValue(subscription.SubscriptionGUID),
			SubscriptionUrl:           types.StringValue(subscription.SubscriptionUrl),
			SupportsParametersUpdates: types.BoolValue(subscription.SupportsParametersUpdates),
			SupportsPlanUpdates:       types.BoolValue(subscription.SupportsPlanUpdates),
			TenantId:                  types.StringValue(subscription.TenantId),
		}

		value.AdditionalPlanFeatures, diags = types.SetValueFrom(ctx, types.StringType, subscription.AdditionalPlanFeatures)
		resp.Diagnostics.Append(diags...)

		value.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, subscription.Labels)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, value)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
