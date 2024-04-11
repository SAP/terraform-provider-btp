package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/SAP/terraform-provider-btp/internal/validation/jsonvalidator"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountServiceInstanceResource() resource.Resource {
	return &subaccountServiceInstanceResource{}
}

// We need to manually model the Service Manager Error structure
// Reason: the Swagger document is not consistent with the actual API response
type svcMgrError struct {
	Error        string
	Description  string
	Broker_Error brokerError
}

type brokerError struct {
	StatusCode    int32
	ErrorMessage  string
	Description   string
	ResponseError string
}

type subaccountServiceInstanceResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountServiceInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_instance", req.ProviderTypeName)
}

func (rs *subaccountServiceInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountServiceInstanceResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a service instance in a subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service instance.",
				Required:            true,
			},
			"serviceplan_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service plan.",
				Required:            true,
			},
			"labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				MarkdownDescription: "The set of words or phrases assigned to the service instance.",
				Computed:            true,
				Optional:            true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "The configuration parameters for the service instance.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					jsonvalidator.ValidJSON(),
				},
			},
			"shared": schema.BoolAttribute{
				MarkdownDescription: "The configuration parameter for service instance sharing. Ensure that the instance is created with a plan that supports instance sharing.",
				Optional:            true,
				Computed:            true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create:            true,
				CreateDescription: "Timeout for creating the service instance.",
				Update:            true,
				UpdateDescription: "Timeout for updating the service instance.",
				Delete:            true,
				DeleteDescription: "Timeout for deleting the service instance.",
			}),
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service instance.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ready": schema.BoolAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"platform_id": schema.StringAttribute{
				MarkdownDescription: "The platform ID.",
				Computed:            true,
			},
			"referenced_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the instance to which the service instance refers.",
				Computed:            true,
			},
			"context": schema.StringAttribute{
				MarkdownDescription: "Contextual data for the resource.",
				Computed:            true,
			},
			"usable": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the resource can be used.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the service instance.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
		},
	}
}

func (rs *subaccountServiceInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountServiceInstanceType
	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	timeoutsLocal := state.Timeouts

	cliRes, rawRes, err := rs.cli.Services.Instance.GetById(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		handleReadErrors(ctx, rawRes, resp, err, "Resource Service Instance (Subaccount)")
		return
	}

	newState, diags := subaccountServiceInstanceValueFrom(ctx, cliRes)
	newState.Timeouts = timeoutsLocal

	// Handle resource import
	if cliRes.Parameters != "" && state.Parameters.ValueString() == "" {
		newState.Parameters = types.StringValue(cliRes.Parameters)
	} else {
		newState.Parameters = state.Parameters
	}

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountServiceInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountServiceInstanceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliReq := btpcli.ServiceInstanceCreateInput{
		Subaccount:    plan.SubaccountId.ValueString(),
		Name:          plan.Name.ValueString(),
		ServicePlanId: plan.ServicePlanId.ValueString(),
	}

	if !plan.Parameters.IsNull() {
		params := plan.Parameters.ValueString()
		cliReq.Parameters = &params
	}

	if !plan.Labels.IsNull() {
		var labels map[string][]string
		plan.Labels.ElementsAs(ctx, &labels, false)

		cliReq.Labels = labels
	}

	cliRes, _, err := rs.cli.Services.Instance.Create(ctx, &cliReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Service Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	createStateConf, diags := rs.CreateStateChange(ctx, cliRes, plan, "creation")
	resp.Diagnostics.Append(diags...)
	updatedRes, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Service Instance (Subaccount)", fmt.Sprintf("%s", err))
	}

	state, diags := subaccountServiceInstanceValueFrom(ctx, updatedRes.(servicemanager.ServiceInstanceResponseObject))
	state.Parameters = plan.Parameters
	state.Timeouts = plan.Timeouts
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if plan.Shared.ValueBool() {
		cliReq := btpcli.ServiceInstanceShareInput{
			Id:         updatedRes.(servicemanager.ServiceInstanceResponseObject).Id,
			Subaccount: updatedRes.(servicemanager.ServiceInstanceResponseObject).SubaccountId,
			Name:       updatedRes.(servicemanager.ServiceInstanceResponseObject).Name,
		}

		cliRes, _, err = rs.cli.Services.Instance.Share(ctx, &cliReq)
		if err != nil {
			resp.Diagnostics.AddError("API Error Sharing Resource Service Instance (Subaccount) while Creating", fmt.Sprintf("%s", err))
			return
		}

		createStateConf, diags := rs.CreateStateChange(ctx, cliRes, plan, "sharing")
		resp.Diagnostics.Append(diags...)
		updatedRes, err := createStateConf.WaitForStateContext(ctx)
		if err != nil {
			resp.Diagnostics.AddError("API Error Sharing Resource Service Instance (Subaccount) while Creating", fmt.Sprintf("%s", err))
		}

		state, diags = subaccountServiceInstanceValueFrom(ctx, updatedRes.(servicemanager.ServiceInstanceResponseObject))
		state.Parameters = plan.Parameters
		state.Timeouts = plan.Timeouts
		resp.Diagnostics.Append(diags...)

		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
	}
}

func (rs *subaccountServiceInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var stateCurrent, plan subaccountServiceInstanceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diagsState := req.State.Get(ctx, &stateCurrent)
	resp.Diagnostics.Append(diagsState...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliReq := btpcli.ServiceInstanceUpdateInput{
		Subaccount:    plan.SubaccountId.ValueString(),
		Id:            plan.Id.ValueString(),
		NewName:       plan.Name.ValueString(),
		ServicePlanId: plan.ServicePlanId.ValueString(),
	}

	if !plan.Parameters.IsNull() {
		params := plan.Parameters.ValueString()
		cliReq.Parameters = &params
	}

	// Labels of plan and state need to be transferred as a delta must be computed for the update operation
	if !plan.Labels.IsNull() {
		var labelsFromPlan map[string][]string
		plan.Labels.ElementsAs(ctx, &labelsFromPlan, false)

		cliReq.LabelsPlan = labelsFromPlan
	}

	if !stateCurrent.Labels.IsNull() {
		var labelsFromState map[string][]string
		stateCurrent.Labels.ElementsAs(ctx, &labelsFromState, false)

		cliReq.LabelsState = labelsFromState
	}

	cliRes, _, err := rs.cli.Services.Instance.Update(ctx, &cliReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Service Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	updateStateConf, diags := rs.UpdateStateChange(ctx, cliRes, plan, "update")
	resp.Diagnostics.Append(diags...)
	updatedRes, err := updateStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Service Instance (Subaccount)", fmt.Sprintf("%s", err))
	}

	state, diags := subaccountServiceInstanceValueFrom(ctx, updatedRes.(servicemanager.ServiceInstanceResponseObject))
	state.Parameters = plan.Parameters
	state.Timeouts = plan.Timeouts
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

	if !plan.Shared.IsUnknown() && plan.Shared != stateCurrent.Shared {

		cliReq := btpcli.ServiceInstanceShareInput{
			Id:         cliRes.Id,
			Subaccount: cliRes.SubaccountId,
			Name:       cliRes.Name,
		}

		if plan.Shared.ValueBool() && !stateCurrent.Shared.ValueBool() {
			cliRes, _, err = rs.cli.Services.Instance.Share(ctx, &cliReq)
			if err != nil {
				resp.Diagnostics.AddError("API Error Sharing Resource Service Instance (Subaccount) while Updating", fmt.Sprintf("%s", err))
				return
			}
		}

		if !plan.Shared.ValueBool() && stateCurrent.Shared.ValueBool() {
			cliRes, _, err = rs.cli.Services.Instance.Unshare(ctx, &cliReq)
			if err != nil {
				resp.Diagnostics.AddError("API Error Unsharing Resource Service Instance (Subaccount) while Updating", fmt.Sprintf("%s", err))
				return
			}
		}

		updateStateConf, diags := rs.UpdateStateChange(ctx, cliRes, plan, "sharing/unsharing")
		resp.Diagnostics.Append(diags...)
		updatedRes, err := updateStateConf.WaitForStateContext(ctx)
		if err != nil {
			resp.Diagnostics.AddError("API Error Sharing Resource Service Instance (Subaccount) while Updating", fmt.Sprintf("%s", err))
		}

		state, diags := subaccountServiceInstanceValueFrom(ctx, updatedRes.(servicemanager.ServiceInstanceResponseObject))
		state.Parameters = plan.Parameters
		state.Timeouts = plan.Timeouts
		resp.Diagnostics.Append(diags...)

		diags = resp.State.Set(ctx, state)
		resp.Diagnostics.Append(diags...)
	}
}

func (rs *subaccountServiceInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountServiceInstanceType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := rs.cli.Services.Instance.Delete(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Service Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, tfutils.DefaultTimeout)
	resp.Diagnostics.Append(diags...)
	delay, minTimeout := tfutils.CalculateDelayAndMinTimeOut(deleteTimeout)

	deleteStateConf := &tfutils.StateChangeConf{
		Pending: []string{servicemanager.StateInProgress},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			subRes, comRes, err := rs.cli.Services.Instance.GetById(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())

			if comRes.StatusCode == http.StatusNotFound {
				return subRes, "DELETED", nil
			}

			if err != nil {
				return subRes, subRes.LastOperation.State, err
			}

			// No error returned even if operation failed
			if subRes.LastOperation.State == servicemanager.StateFailed {
				opsError := extractDetailedError(subRes.LastOperation.Errors, "deletion")
				return subRes, subRes.LastOperation.State, opsError
			}

			return subRes, subRes.LastOperation.State, nil
		},
		Timeout:    deleteTimeout,
		Delay:      delay,
		MinTimeout: minTimeout,
	}

	_, err = deleteStateConf.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Service Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountServiceInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: subaccount_id,id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (rs *subaccountServiceInstanceResource) CreateStateChange(ctx context.Context, cliRes servicemanager.ServiceInstanceResponseObject, plan subaccountServiceInstanceType, operation string) (tfutils.StateChangeConf, diag.Diagnostics) {

	var summary diag.Diagnostics

	state, diags := subaccountServiceInstanceValueFrom(ctx, cliRes)
	state.Parameters = plan.Parameters
	summary.Append(diags...)

	timeoutsLocal := plan.Timeouts
	createTimeout, diags := timeoutsLocal.Create(ctx, tfutils.DefaultTimeout)
	summary.Append(diags...)
	delay, minTimeout := tfutils.CalculateDelayAndMinTimeOut(createTimeout)

	return tfutils.StateChangeConf{
			Pending: []string{servicemanager.StateInProgress},
			Target:  []string{servicemanager.StateSucceeded},
			Refresh: func() (interface{}, string, error) {
				subRes, _, err := rs.cli.Services.Instance.GetById(ctx, state.SubaccountId.ValueString(), cliRes.Id)

				if err != nil {
					return subRes, "", err
				}

				// No error returned even if operation failed
				if subRes.LastOperation.State == servicemanager.StateFailed {
					opsError := extractDetailedError(subRes.LastOperation.Errors, operation)
					return subRes, subRes.LastOperation.State, opsError
				}

				return subRes, subRes.LastOperation.State, nil
			},
			Timeout:    createTimeout,
			Delay:      delay,
			MinTimeout: minTimeout,
		},
		summary
}

func (rs *subaccountServiceInstanceResource) UpdateStateChange(ctx context.Context, cliRes servicemanager.ServiceInstanceResponseObject, plan subaccountServiceInstanceType, operation string) (tfutils.StateChangeConf, diag.Diagnostics) {

	var summary diag.Diagnostics

	timeoutsLocal := plan.Timeouts
	state, diags := subaccountServiceInstanceValueFrom(ctx, cliRes)
	state.Parameters = plan.Parameters
	summary.Append(diags...)

	updateTimeout, diags := timeoutsLocal.Update(ctx, tfutils.DefaultTimeout)
	summary.Append(diags...)
	delay, minTimeout := tfutils.CalculateDelayAndMinTimeOut(updateTimeout)

	return tfutils.StateChangeConf{
		Pending: []string{servicemanager.StateInProgress},
		Target:  []string{servicemanager.StateSucceeded},
		Refresh: func() (interface{}, string, error) {
			subRes, _, err := rs.cli.Services.Instance.GetById(ctx, state.SubaccountId.ValueString(), cliRes.Id)

			if err != nil {
				return subRes, "", err
			}

			// No error returned even if operation failed
			if subRes.LastOperation.State == servicemanager.StateFailed {
				opsError := extractDetailedError(subRes.LastOperation.Errors, operation)
				return subRes, subRes.LastOperation.State, opsError
			}

			return subRes, subRes.LastOperation.State, nil
		},
		Timeout:    updateTimeout,
		Delay:      delay,
		MinTimeout: minTimeout,
	}, summary
}

func extractDetailedError(errorAsRawJson json.RawMessage, operation string) error {
	// Manual extraction of the error message details (human readable error mesage) from the API response
	var modelErr svcMgrError
	if err := json.Unmarshal(errorAsRawJson, &modelErr); err == nil {
		return errors.New("API error during service instance " + operation + " - " + modelErr.Broker_Error.Description)
	}

	return errors.New("undefined API error during service instance " + operation)
}
