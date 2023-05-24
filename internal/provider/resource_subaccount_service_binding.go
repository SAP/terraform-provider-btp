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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	resourcehelper "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/tfdefaults"
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
		MarkdownDescription: `Create a service binding between a service instance and an application.`,
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
				PlanModifiers: []planmodifier.String{
					tfdefaults.DefaultStringValue(`{}`),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					jsonvalidator.ValidJSON(),
				},
			},
			"labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				MarkdownDescription: "Set of words or phrases assigned to service binding.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service binding.",
				Computed:            true,
			},
			"ready": schema.BoolAttribute{
				MarkdownDescription: "Whether the service binding is ready.",
				Computed:            true,
			},
			"context": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Contextual data for the resource.",
				Computed:            true,
			},
			"bind_resource": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Contains the resources associated with the binding.",
				Computed:            true,
			},
			"credentials": schema.StringAttribute{
				MarkdownDescription: "Credentials to access the binding.",
				Computed:            true,
				Sensitive:           true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "Current state of the service binding.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
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

	cliRes, _, err := rs.cli.Services.Binding.GetById(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Binding (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := subaccountServiceBindingValueFrom(ctx, cliRes)
	updatedState.Parameters = state.Parameters
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

	cliRes, _, err := rs.cli.Services.Binding.Create(ctx, cliReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Service Binding (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	updatedPlan, diags := subaccountServiceBindingValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	createStateConf := &resourcehelper.StateChangeConf{
		Pending: []string{servicemanager.StateInProgress},
		Target:  []string{servicemanager.StateSucceeded, servicemanager.StateFailed},
		Refresh: func() (interface{}, string, error) {
			subRes, _, err := rs.cli.Services.Binding.GetById(ctx, plan.SubaccountId.ValueString(), cliRes.Id)

			if err != nil {
				return subRes, "", err
			}

			return subRes, subRes.LastOperation.State, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	updatedRes, err := createStateConf.WaitForStateContext(ctx)

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

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
