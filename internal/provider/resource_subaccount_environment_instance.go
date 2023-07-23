package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/provisioning"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountEnvironmentInstanceResource() resource.Resource {
	return &subaccountEnvironmentInstanceResource{}
}

type subaccountEnvironmentInstanceResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountEnvironmentInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_environment_instance", req.ProviderTypeName)
}

func (rs *subaccountEnvironmentInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountEnvironmentInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates an environment instance, such as a Cloud Foundry org, in a subaccount.

__Tips:__
* You must be assigned to the subaccount admin role.
* Quota-based environments, such as Kyma, must first be assigned as entitlements to the subaccount.

__Further documentation:__
* Cloud Foundry: <https://help.sap.com/docs/btp/sap-business-technology-platform/org-management-using-sap-btp-command-line-interface-btp-cli>
* Kyma: <https://help.sap.com/docs/btp/sap-business-technology-platform/available-plans-in-kyma-environment>
* Concept: <https://help.sap.com/docs/btp/sap-business-technology-platform/environments>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the environment instance.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"environment_type": schema.StringAttribute{
				MarkdownDescription: "The type of the environment instance that is used.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plan_name": schema.StringAttribute{
				MarkdownDescription: "The name of the service plan for the environment instance in the corresponding service broker's catalog.",
				Required:            true,
			},
			"service_name": schema.StringAttribute{
				MarkdownDescription: "The name of the service for the environment instance in the corresponding service broker's catalog.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"landscape_label": schema.StringAttribute{
				MarkdownDescription: "The name of the landscape within the logged in region on which the environment instance is created.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "The configuration parameters for the environment instance.",
				Optional:            true,
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the environment instance.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"broker_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the associated environment broker.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"custom_labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				MarkdownDescription: "The set of words or phrases assigned to the environment instance.",
				Computed:            true,
			},
			"dashboard_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the service dashboard, which is a web-based management user interface for the service instances.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the environment instance.",
				Computed:            true,
			},
			"labels": schema.StringAttribute{
				MarkdownDescription: "The Broker-specified key-value pairs that specify attributes of an environment instance.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"operation": schema.StringAttribute{
				MarkdownDescription: "An identifier that represents the last operation. This ID is returned by the environment brokers.",
				Computed:            true,
			},
			"plan_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service plan for the environment instance in the corresponding service broker's catalog.",
				Computed:            true,
			},
			"platform_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the platform for the environment instance in the corresponding service broker's catalog.",
				Computed:            true,
			},
			"service_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service for the environment instance in the corresponding service broker's catalog.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the environment instance. Possible values are: \n" +
					getFormattedValueAsTableRow("state", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`OK`", "The CRUD operation or series of operations completed successfully.") +
					getFormattedValueAsTableRow("`CREATING`", "Creating of the environment instance is in progress.") +
					getFormattedValueAsTableRow("`CREATION_FAILED`", "The creation of the environment instance failed, and the environment instance was not created or was created but cannot be used.") +
					getFormattedValueAsTableRow("`UPDATING`", "Updating of the environment instance is in progress.") +
					getFormattedValueAsTableRow("`UPDATE_FAILED`", "The update of the environment instance failed, and  the environment instance was not updated.") +
					getFormattedValueAsTableRow("`DELETING`", "Deleting of the environment instance is in progress.") +
					getFormattedValueAsTableRow("`DELETION_FAILED`", "The deletion of the environment instance failed, and the environment instance was not deleted."), Computed: true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant that owns the environment instance.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The last provisioning operation on the environment instance. Possible values are: \n" +
					getFormattedValueAsTableRow("type", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`Provision`", "The environment instance is created.") +
					getFormattedValueAsTableRow("`Update`", "The environment instance is changed.") +
					getFormattedValueAsTableRow("`Deprovision`", "The environment instance is deleted."),
				Computed: true,
			},
		},
	}
}

func (rs *subaccountEnvironmentInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountEnvironmentInstanceType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Accounts.EnvironmentInstance.Get(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Environment Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := subaccountEnvironmentInstanceValueFrom(ctx, cliRes)

	if !state.Parameters.IsNull() {
		updatedState.Parameters = state.Parameters
	}

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountEnvironmentInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountEnvironmentInstanceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	parameters := plan.Parameters.ValueString()

	cliRes, _, err := rs.cli.Accounts.EnvironmentInstance.Create(ctx, &btpcli.SubaccountEnvironmentInstanceCreateInput{
		SubaccountID:    plan.SubaccountId.ValueString(),
		DisplayName:     plan.Name.ValueString(),
		Service:         plan.ServiceName.ValueString(),
		Plan:            plan.PlanName.ValueString(),
		EnvironmentType: plan.EnvironmentType.ValueString(),
		Landscape:       plan.LandscapeLabel.ValueString(),
		Parameters:      parameters,
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Environment Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	plan, diags = subaccountEnvironmentInstanceValueFrom(ctx, cliRes)
	plan.Parameters = types.StringValue(parameters)
	resp.Diagnostics.Append(diags...)

	createStateConf := &tfutils.StateChangeConf{
		Pending: []string{provisioning.StateCreating},
		Target:  []string{provisioning.StateOK, provisioning.StateCreationFailed},
		Refresh: func() (interface{}, string, error) {
			subRes, _, err := rs.cli.Accounts.EnvironmentInstance.Get(ctx, plan.SubaccountId.ValueString(), cliRes.Id)

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
		resp.Diagnostics.AddError("API Error Creating Resource Environment Instance (Subaccount)", fmt.Sprintf("%s", err))
	}

	plan, diags = subaccountEnvironmentInstanceValueFrom(ctx, updatedRes.(provisioning.EnvironmentInstanceResponseObject))
	plan.Parameters = types.StringValue(parameters)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountEnvironmentInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountEnvironmentInstanceType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Accounts.EnvironmentInstance.Update(ctx, &btpcli.SubaccountEnvironmentInstanceUpdateInput{
		EnvironmentID: plan.Id.ValueString(),
		Parameters:    plan.Parameters.ValueString(),
		Plan:          plan.PlanName.ValueString(),
		SubaccountID:  plan.SubaccountId.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Environment Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	updateStateConf := &tfutils.StateChangeConf{
		Pending: []string{provisioning.StateUpdating},
		Target:  []string{provisioning.StateOK, provisioning.StateUpdateFailed},
		Refresh: func() (interface{}, string, error) {
			subRes, _, err := rs.cli.Accounts.EnvironmentInstance.Get(ctx, plan.SubaccountId.ValueString(), plan.Id.ValueString())

			if err != nil {
				return subRes, "", err
			}

			return subRes, subRes.State, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	updatedRes, err := updateStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Environment Instance (Subaccount)", fmt.Sprintf("%s", err))
	}

	state, diags := subaccountEnvironmentInstanceValueFrom(ctx, updatedRes.(provisioning.EnvironmentInstanceResponseObject))
	// TODO: this temporary workaround ignores the actual "parameters" value which is diverging from the planned state by an additional "status" attribute
	state.Parameters = plan.Parameters
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountEnvironmentInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountEnvironmentInstanceType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Accounts.EnvironmentInstance.Delete(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Environment Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	deleteStateConf := &tfutils.StateChangeConf{
		Pending: []string{provisioning.StateDeleting},
		Target:  []string{"DELETED", provisioning.StateDeletionFailed},
		Refresh: func() (interface{}, string, error) {
			subRes, comRes, err := rs.cli.Accounts.EnvironmentInstance.Get(ctx, state.SubaccountId.ValueString(), cliRes.Id)

			if comRes.StatusCode == http.StatusNotFound {
				return subRes, "DELETED", nil
			}

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
		resp.Diagnostics.AddError("API Error Deleting Resource Environment Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountEnvironmentInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: subaccount_id,environment_instance_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}
