package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
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
)

const updateSubscriptionResource = "UpdateResource"
const updateTimeoutOnly = "UpdateTimeoutOnly"
const invalidUpdateRequest = "InvalidUpdateRequest"

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

func (rs *subaccountSubscriptionResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Subscribes a subaccount to a multitenant application.
Custom or partner-developed applications are currently not supported.

__Tip:__
You must be assigned to the admin role of the subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"app_name": schema.StringAttribute{
				MarkdownDescription: "The unique registration name of the deployed multitenant application as defined by the app developer.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					jsonvalidator.ValidJSON(),
				},
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create:            true,
				CreateDescription: "Timeout for creating the subscription.",
				Update:            true,
				UpdateDescription: "Timeout for updating the subscription.",
				Delete:            true,
				DeleteDescription: "Timeout for deleting the subscription.",
			}),
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

type SubaccountSubscriptionResourceIdentityModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	AppName      types.String `tfsdk:"app_name"`
	PlanName     types.String `tfsdk:"plan_name"`
}

func (rs *subaccountSubscriptionResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"subaccount_id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"app_name": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"plan_name": identityschema.StringAttribute{
				RequiredForImport: true,
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

	timeoutsLocal := state.Timeouts

	// We determine the technical app name as this is needed by the API
	technicalAppName, _, err := rs.determineAppNames(ctx, state.SubaccountId.ValueString(), state.AppName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error determining subscription app name (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	cliRes, rawRes, err := rs.cli.Accounts.Subscription.Get(ctx, state.SubaccountId.ValueString(), technicalAppName, state.PlanName.ValueString())
	if err != nil || cliRes.State == saas_manager_service.StateNotSubscribed {
		handleReadErrors(ctx, rawRes, cliRes, resp, err, "Resource Subscription (Subaccount)")
		return
	}

	newState, diags := subaccountSubscriptionValueFrom(ctx, cliRes)
	newState.AppName = state.AppName
	newState.Timeouts = timeoutsLocal

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

	var identity SubaccountSubscriptionResourceIdentityModel

	diags = req.Identity.Get(ctx, &identity)
	if diags.HasError() {
		identity = SubaccountSubscriptionResourceIdentityModel{
			SubaccountId: types.StringValue(state.SubaccountId.ValueString()),
			AppName:      types.StringValue(cliRes.AppName),
			PlanName:     types.StringValue(cliRes.PlanName),
		}

		diags = resp.Identity.Set(ctx, identity)
		resp.Diagnostics.Append(diags...)
	}
}

func (rs *subaccountSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountSubscriptionType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// We determine the technical app name as this is needed by the API
	technicalAppName, _, err := rs.determineAppNames(ctx, plan.SubaccountId.ValueString(), plan.AppName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error determining subscription app name (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	_, _, err = rs.cli.Accounts.Subaccount.Subscribe(ctx, plan.SubaccountId.ValueString(), technicalAppName, plan.PlanName.ValueString(), plan.Parameters.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	createStateConf, diags := rs.CreateStateChange(ctx, plan, "create", technicalAppName)
	resp.Diagnostics.Append(diags...)

	updatedRes, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
	}

	if updatedRes == nil {
		return
	}

	updatedPlan, diags := subaccountSubscriptionValueFrom(ctx, updatedRes.(saas_manager_service.EntitledApplicationsResponseObject))
	// We must override the API values with the plan values as we might have had to change the app name
	// due to a mismatch of the technical and commercial app name
	updatedPlan.AppName = plan.AppName
	updatedPlan.Parameters = plan.Parameters
	updatedPlan.Timeouts = plan.Timeouts
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedPlan)
	resp.Diagnostics.Append(diags...)

	identity := SubaccountSubscriptionResourceIdentityModel{
		SubaccountId: types.StringValue(plan.SubaccountId.ValueString()),
		AppName:      types.StringValue(updatedPlan.AppName.ValueString()),
		PlanName:     types.StringValue(updatedPlan.PlanName.ValueString()),
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountSubscriptionType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state subaccountSubscriptionType
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.PlanName.ValueString() != state.PlanName.ValueString() && !state.SupportsPlanUpdates.ValueBool() {
		resp.Diagnostics.AddError("API Error Updating Resource Subscription (Subaccount)", "Plan name is not supposed to be updated for this resource")
		return
	}
	if plan.Parameters.ValueString() != state.Parameters.ValueString() && !state.SupportsParametersUpdates.ValueBool() {
		resp.Diagnostics.AddError("API Error Updating Resource Subscription (Subaccount)", "Parameters are not supposed to be updated for this resource")
		return
	}

	updateType := checkForChanges(plan, state)

	switch updateType {
	case invalidUpdateRequest:
		// The update tries to access fields, which are not supposed to be updated
		resp.Diagnostics.AddError("API Error Updating Subscription (Subaccount)", "This provided parameters are not supposed to be updated")
	case updateSubscriptionResource:
		_, _, err := rs.cli.Accounts.Subscription.Update(ctx, plan.SubaccountId.ValueString(), plan.AppName.ValueString(), plan.PlanName.ValueString(), plan.Parameters.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Updating Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
			return
		}

		updateStateConf, diags := rs.UpdateStateChange(ctx, plan, "update")
		resp.Diagnostics.Append(diags...)

		updatedRes, err := updateStateConf.WaitForStateContext(ctx)
		if err != nil {
			resp.Diagnostics.AddError("API Error Updating Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
		}

		if updatedRes == nil {
			return
		}

		updatedPlan, diags := subaccountSubscriptionValueFrom(ctx, updatedRes.(saas_manager_service.EntitledApplicationsResponseObject))
		updatedPlan.Parameters = plan.Parameters
		updatedPlan.Timeouts = plan.Timeouts
		resp.Diagnostics.Append(diags...)

		diags = resp.State.Set(ctx, &updatedPlan)
		resp.Diagnostics.Append(diags...)

	case updateTimeoutOnly:
		// The update only tries to access the timeouts, so we do not need to call the API, we just set the state values back into the state and update the timeouts only
		// Reason: The API implementations are not all idempotent which could lead to errors if we call the API for an UPDATe even if no fields were changed
		updatedPlan := state
		updatedPlan.Timeouts = plan.Timeouts
		diags = resp.State.Set(ctx, &updatedPlan)
		resp.Diagnostics.Append(diags...)
	}

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

	// We determine the technical app name as this is needed by the API
	technicalAppName, _, err := rs.determineAppNames(ctx, state.SubaccountId.ValueString(), state.AppName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error determining subscription app name (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	_, _, err = rs.cli.Accounts.Subaccount.Unsubscribe(ctx, state.SubaccountId.ValueString(), technicalAppName)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	deleteStateConf, diags := rs.DeleteStateChange(ctx, state, "delete", technicalAppName)
	resp.Diagnostics.Append(diags...)

	_, err = deleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Subscription (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" {
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
		return
	}

	var identityData SubaccountSubscriptionResourceIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), identityData.SubaccountId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_name"), identityData.AppName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("plan_name"), identityData.PlanName)...)
}

func (rs *subaccountSubscriptionResource) CreateStateChange(ctx context.Context, plan subaccountSubscriptionType, operation string, technicalAppName string) (tfutils.StateChangeConf, diag.Diagnostics) {
	var summary diag.Diagnostics

	timeoutsLocal := plan.Timeouts

	createTimeout, diags := timeoutsLocal.Create(ctx, tfutils.DefaultTimeout)
	summary.Append(diags...)
	delay, minTimeout := tfutils.CalculateDelayAndMinTimeOut(createTimeout)

	return tfutils.StateChangeConf{
			Pending: []string{saas_manager_service.StateInProcess},
			Target:  []string{saas_manager_service.StateSubscribed},
			Refresh: func() (interface{}, string, error) {
				subRes, cmdRes, err := rs.cli.Accounts.Subscription.Get(ctx, plan.SubaccountId.ValueString(), technicalAppName, plan.PlanName.ValueString())

				if cmdRes.StatusCode == http.StatusTooManyRequests {
					// Retry in case of rate limiting
					return subRes, saas_manager_service.StateInProcess, nil
				}

				if err != nil {
					return subRes, "", err
				}

				// No error returned even is subscription failed
				if subRes.State == saas_manager_service.StateSubscribeFailed {
					errorMessage := getErrorFromResponse(subRes)
					return subRes, subRes.State, errors.New(errorMessage)
				}

				return subRes, subRes.State, nil
			},
			Timeout:    createTimeout,
			Delay:      delay,
			MinTimeout: minTimeout,
		},
		summary
}

func (rs *subaccountSubscriptionResource) DeleteStateChange(ctx context.Context, state subaccountSubscriptionType, operation string, technicalAppName string) (tfutils.StateChangeConf, diag.Diagnostics) {

	var summary diag.Diagnostics

	timeoutsLocal := state.Timeouts

	deleteTimeout, diags := timeoutsLocal.Delete(ctx, tfutils.DefaultTimeout)
	summary.Append(diags...)
	delay, minTimeout := tfutils.CalculateDelayAndMinTimeOut(deleteTimeout)

	return tfutils.StateChangeConf{
			Pending: []string{saas_manager_service.StateInProcess},
			Target:  []string{saas_manager_service.StateNotSubscribed},
			Refresh: func() (interface{}, string, error) {
				subRes, cmdRes, err := rs.cli.Accounts.Subscription.Get(ctx, state.SubaccountId.ValueString(), technicalAppName, state.PlanName.ValueString())

				if cmdRes.StatusCode == http.StatusTooManyRequests {
					// Retry in case of rate limiting
					return subRes, saas_manager_service.StateInProcess, nil
				}

				if err != nil {
					return subRes, subRes.State, err
				}

				// No error returned even is unsubscribe failed
				if subRes.State == saas_manager_service.StateUnsubscribeFailed {
					return subRes, subRes.State, errors.New("undefined API error during unsubscription")
				}

				return subRes, subRes.State, nil
			},
			Timeout:    deleteTimeout,
			Delay:      delay,
			MinTimeout: minTimeout,
		},
		summary
}

func (rs *subaccountSubscriptionResource) UpdateStateChange(ctx context.Context, plan subaccountSubscriptionType, operation string) (tfutils.StateChangeConf, diag.Diagnostics) {

	var summary diag.Diagnostics

	timeoutsLocal := plan.Timeouts

	updateTimeout, diags := timeoutsLocal.Update(ctx, tfutils.DefaultTimeout)
	summary.Append(diags...)
	delay, minTimeout := tfutils.CalculateDelayAndMinTimeOut(updateTimeout)

	return tfutils.StateChangeConf{
			Pending: []string{saas_manager_service.StateInProcess},
			Target:  []string{saas_manager_service.StateSubscribed},
			Refresh: func() (interface{}, string, error) {
				subRes, cmdRes, err := rs.cli.Accounts.Subscription.Get(ctx, plan.SubaccountId.ValueString(), plan.AppName.ValueString(), plan.PlanName.ValueString())

				if cmdRes.StatusCode == http.StatusTooManyRequests {
					// Retry in case of rate limiting
					return subRes, saas_manager_service.StateInProcess, nil
				}

				if err != nil {
					return subRes, subRes.State, err
				}

				// No error returned even is subscription failed
				if subRes.State == saas_manager_service.StateSubscribeFailed {
					return subRes, subRes.State, errors.New("undefined API error during updating subscription")
				}

				return subRes, subRes.State, nil
			},
			Timeout:    updateTimeout,
			Delay:      delay,
			MinTimeout: minTimeout,
		},
		summary
}

func getErrorFromResponse(subRes saas_manager_service.EntitledApplicationsResponseObject) (errorMessage string) {
	errorMessage = "undefined API error during subscription"

	if subRes.SubscriptionError == nil {
		return errorMessage
	}

	if subRes.SubscriptionError.AppError != "" {
		return subRes.SubscriptionError.AppError
	} else {
		return errorMessage
	}
}

func (rs *subaccountSubscriptionResource) determineAppNames(ctx context.Context, subaccountId string, planAppName string) (technicalAppName string, commercialAppName string, err error) {
	// The caller might hand over the technical or commercial app name
	// We ensure that the right name is used for the subscription namely in the API calls
	// This only works in consumer subaccounts as the List command will filter the apps to make then unique
	appList, _, err := rs.cli.Accounts.Subscription.List(ctx, subaccountId)

	if err != nil {
		return "", "", err
	}

	for _, app := range appList {
		if app.AppName == planAppName || app.CommercialAppName == planAppName {
			technicalAppName = app.AppName
			commercialAppName = app.CommercialAppName
			return technicalAppName, commercialAppName, nil
		}
	}

	// We did not find the app name in the list of applications, default to the one handed over by the user for the technical app name
	// The API will return an error if the app name is not valid
	return planAppName, "", nil
}

func checkForChanges(plan subaccountSubscriptionType, state subaccountSubscriptionType) string {
	// We check what kind of update is requested to distinguish if an API call is necessary or not
	if plan.PlanName.ValueString() != state.PlanName.ValueString() || plan.Parameters.ValueString() != state.Parameters.ValueString() {
		return updateSubscriptionResource
	}

	if !plan.Timeouts.Equal(state.Timeouts) {
		// An update of the timeouts can especially happen during import of the resource
		return updateTimeoutOnly
	}

	return invalidUpdateRequest
}
