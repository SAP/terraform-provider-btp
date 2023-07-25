package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/saas_manager_service"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/SAP/terraform-provider-btp/internal/validation/jsonvalidator"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountSubscriptionResource() resource.Resource {
	return &subaccountSubscriptionResource{}
}

type subaccountSubscriptionResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountSubscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_subscription", req.ProviderTypeName)
}

func (rs *subaccountSubscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountSubscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Subscribes to a multitenant application from a subaccount.
Custom or partner-developed applications are currently not supported.

__Tip:__
You must be assigned to the subaccount admin role.`,
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
			},
			"plan_name": schema.StringAttribute{
				MarkdownDescription: "The plan name of the application to which the consumer has subscribed.",
				Required:            true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "The parameters of the subscription as a valid JSON object.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(`{}`),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					jsonvalidator.ValidJSON(),
				},
			},
			"additional_plan_features": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The list of features specific to this plan.",
				Computed:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The ID returned by XSUAA after the app provider has performed a bind of the multitenant application to an XSUAA service instance.",
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
				MarkdownDescription: "The description of the multitenant application for customer-facing UIs.",
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
				MarkdownDescription: "The ID of the subaccount, which is subscribed to the multitenant application.",
				Computed:            true,
			},
			"subscribed_tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant, which is subscribed to a multitenant application.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The technical ID generated by XSUAA for a multitenant application when a consumer subscribes to the application.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
	}
}

func (rs *subaccountSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountSubscriptionType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Accounts.Subscription.Get(ctx, state.SubaccountId.ValueString(), state.AppName.ValueString(), state.PlanName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	newState, diags := subaccountSubscriptionValueFrom(ctx, cliRes)

	if newState.Parameters.IsNull() && !state.Parameters.IsNull() {
		// The parameters are not returned by the API so we transfer the existing state to the read result if not existing
		newState.Parameters = state.Parameters
	} else if newState.Parameters.IsNull() && state.Parameters.IsNull() {
		// During the import of the resource both values might be empty, so we need to apply the default value form the schema if not existing
		newState.Parameters = types.StringValue("{}")
	}

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountSubscriptionType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Accounts.Subaccount.Subscribe(ctx, plan.SubaccountId.ValueString(), plan.AppName.ValueString(), plan.PlanName.ValueString(), plan.Parameters.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	createStateConf := &tfutils.StateChangeConf{
		Pending: []string{saas_manager_service.StateInProcess},
		Target:  []string{saas_manager_service.StateSubscribed, saas_manager_service.StateSubscribeFailed},
		Refresh: func() (interface{}, string, error) {
			subRes, _, err := rs.cli.Accounts.Subscription.Get(ctx, plan.SubaccountId.ValueString(), plan.AppName.ValueString(), plan.PlanName.ValueString())

			if err != nil {
				return subRes, "", err
			}

			return subRes, subRes.State, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	updatedRes, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
	}

	updatedPlan, diags := subaccountSubscriptionValueFrom(ctx, updatedRes.(saas_manager_service.EntitledApplicationsResponseObject))
	updatedPlan.Parameters = plan.Parameters
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedPlan)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountSubscriptionType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Subscription (Subaccount)", "This resource is not supposed to be updated")
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *subaccountSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountSubscriptionType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Accounts.Subaccount.Unsubscribe(ctx, state.SubaccountId.ValueString(), state.AppName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	deleteStateConf := &tfutils.StateChangeConf{
		Pending: []string{saas_manager_service.StateInProcess},
		Target:  []string{saas_manager_service.StateUnsubscribeFailed, saas_manager_service.StateNotSubscribed},
		Refresh: func() (interface{}, string, error) {
			subRes, _, err := rs.cli.Accounts.Subscription.Get(ctx, state.SubaccountId.ValueString(), state.AppName.ValueString(), state.PlanName.ValueString())

			if err != nil {
				return subRes, subRes.State, err
			}

			return subRes, subRes.State, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: subaccount,app_name,plan_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("plan_name"), idParts[2])...)
}
