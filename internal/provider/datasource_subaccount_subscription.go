package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountSubscriptionDataSource() datasource.DataSource {
	return &subaccountSubscriptionDataSource{}
}

type subaccountSubscriptionDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountSubscriptionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_subscription", req.ProviderTypeName)
}

func (ds *subaccountSubscriptionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountSubscriptionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets the details of a specific multitenant application to which a subaccount is entitled to subscribe. If this application is in a different global account than the current one, you need to specify its plan with '--plan'.

__Tip:__
You must be assigned to the subaccount admin or viewer role.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"app_name": schema.StringAttribute{
				MarkdownDescription: "The unique registration name of the deployed multitenant application as defined by the app developer.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"plan_name": schema.StringAttribute{
				MarkdownDescription: "The plan name of the application to which the consumer has subscribed.",
				Required:            true, // TODO optional
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "The parameters of the subscription as a valid JSON object.",
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
			"category": schema.StringAttribute{
				MarkdownDescription: "The technical name of the category defined by the app developer to which the multitenant application is grouped in customer-facing UIs.",
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
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
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
				MarkdownDescription: "Set of words or phrases assigned to the multitenant application subscription.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountSubscriptionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountSubscriptionType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.Subscription.Get(ctx, data.SubaccountId.ValueString(), data.AppName.ValueString(), data.PlanName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data, diags = subaccountSubscriptionValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
