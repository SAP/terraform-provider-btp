package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/SAP/terraform-provider-btp/internal/validation/jsonvalidator"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountServiceBindingResource() resource.Resource {
	return &subaccountServiceBindingResource{}
}

type subaccountServiceBindingResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountServiceBindingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_binding", req.ProviderTypeName)
}

func (rs *subaccountServiceBindingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountServiceBindingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a service binding, i.e. generates access details to consume a service.

__Tip:__
You must be assigned to the admin or the service administrator role of the subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service instance associated with the binding.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service binding.",
				Required:            true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "The parameters of the service binding as a valid JSON object.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(`{}`),
				CustomType:          jsontypes.NormalizedType{},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					jsonvalidator.ValidJSON(),
				},
			},
			"labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				MarkdownDescription: "The set of words or phrases assigned to the service binding.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service binding.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ready": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service binding is ready.",
				Computed:            true,
			},
			"context": schema.StringAttribute{
				MarkdownDescription: "The contextual data for the resource.",
				Computed:            true,
			},
			"bind_resource": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Contains the resources associated with the binding.",
				Computed:            true,
			},
			"credentials": schema.StringAttribute{
				MarkdownDescription: "The credentials to access the binding.",
				Computed:            true,
				Sensitive:           true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the service binding. Possible values are: \n" +
					getFormattedValueAsTableRow("state", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`in progress`", "The operation or processing is in progress") +
					getFormattedValueAsTableRow("`failed`", "The operation or processing failed") +
					getFormattedValueAsTableRow("`succeeded`", "The operation or processing succeeded"),
				Computed: true,
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

func (rs *subaccountServiceBindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountServiceBindingType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Services.Binding.GetById(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		handleReadErrors(ctx, rawRes, cliRes, resp, err, "Resource Service Binding (Subaccount)")
		return
	}

	updatedState, diags := subaccountServiceBindingValueFrom(ctx, cliRes)

	if updatedState.Parameters.IsNull() && !state.Parameters.IsNull() {
		// The parameters are not returned by the API so we transfer the existing state to the read result if not existing
		updatedState.Parameters = state.Parameters
	} else if updatedState.Parameters.IsNull() && state.Parameters.IsNull() {
		// During the import of the resource both values might be empty, so we need to apply the default value form the schema if not existing
		updatedState.Parameters = jsontypes.NewNormalizedValue("{}")
	}

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountServiceBindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountServiceBindingType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliReq := btpcli.SubaccountServiceBindingCreateInput{
		Subaccount:        plan.SubaccountId.ValueString(),
		ServiceInstanceId: plan.ServiceInstanceId.ValueString(),
		Name:              plan.Name.ValueString(),
		Parameters:        plan.Parameters.ValueString(),
	}

	if !plan.Labels.IsNull() {
		var labels map[string][]string
		plan.Labels.ElementsAs(ctx, &labels, false)

		cliReq.Labels = labels
	}

	cliRes, _, err := rs.cli.Services.Binding.Create(ctx, cliReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Service Binding (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	updatedPlan, diags := subaccountServiceBindingValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	createStateConf := &tfutils.StateChangeConf{
		Pending: []string{servicemanager.StateInProgress},
		Target:  []string{servicemanager.StateSucceeded},
		Refresh: func() (interface{}, string, error) {
			subRes, _, err := rs.cli.Services.Binding.GetById(ctx, plan.SubaccountId.ValueString(), cliRes.Id)

			if err != nil {
				return subRes, "", err
			}

			// No error returned even if operation failed
			if subRes.LastOperation.State == servicemanager.StateFailed {
				return subRes, subRes.LastOperation.State, errors.New("undefined API error during service binding creation")
			}

			return subRes, subRes.LastOperation.State, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	updatedRes, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Service Binding (Subaccount)", fmt.Sprintf("%s", err))
	}

	updatedPlan, diags = subaccountServiceBindingValueFrom(ctx, updatedRes.(servicemanager.ServiceBindingResponseObject))
	updatedPlan.Parameters = plan.Parameters
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedPlan)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountServiceBindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountServiceBindingType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Resource Service Binding (Subaccount)", "This resource is not supposed to be updated")
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *subaccountServiceBindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountServiceBindingType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Services.Binding.Delete(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Service Binding (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	deleteStateConf := &tfutils.StateChangeConf{
		Pending: []string{servicemanager.StateInProgress},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			subRes, comRes, err := rs.cli.Services.Binding.GetById(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())

			if comRes.StatusCode == http.StatusNotFound {
				return subRes, "DELETED", nil
			}

			if err != nil {
				return subRes, servicemanager.StateFailed, err
			}

			// No error returned even if operation failed
			if subRes.LastOperation.State == servicemanager.StateFailed {
				return subRes, subRes.LastOperation.State, errors.New("undefined API error during service binding deletion")
			}

			return subRes, subRes.LastOperation.State, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = deleteStateConf.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Service Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountServiceBindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: subaccount,service_binding_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
