package provider

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis_entitlements"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newDirectoryEntitlementResource() resource.Resource {
	return &directoryEntitlementResource{}
}

type directoryEntitlementResource struct {
	cli *btpcli.ClientFacade
}

func (rs *directoryEntitlementResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_entitlement", req.ProviderTypeName)
}

func (rs *directoryEntitlementResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *directoryEntitlementResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assigns the entitlement plan of a service, multitenant application, or environment, to a directory. Note that some environments, such as Cloud Foundry, are available by default to all global accounts and their directorys, and therefore are not made available as entitlements.

__Tip:__
You must be assigned to the global account admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/entitlements-and-quotas>`,
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the entitled service plan.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service_name": schema.StringAttribute{
				MarkdownDescription: "The name of the entitled service.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plan_name": schema.StringAttribute{
				MarkdownDescription: "The name of the entitled service plan.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"amount": schema.Int64Attribute{
				MarkdownDescription: "The quota assigned to the directory.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 2000000000),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"auto_assign": schema.BoolAttribute{
				MarkdownDescription: "Determines whether the plans of entitlements that have a numeric quota with the amount specified in `auto_distribute_amount` are automatically allocated to any new subaccount that is added to the directory in the future. For entitlements without a numeric quota, it shows if the plan are assigned to any new subaccount that is added to the directory in the future (`auto_distribute_amount` is not needed). If the `distribute` parameter is set, the same assignment is also made to all subaccounts currently in the directory. Entitlements are subject to available quota in the directory.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_distribute_amount": schema.Int64Attribute{
				MarkdownDescription: "The quota of the specified plan automatically allocated to any new subaccount that is created in the future in the directory. When applying this option, `auto_assign` and/or `distribute` must also be set. Applies only to entitlements that have a numeric quota.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"distribute": schema.BoolAttribute{
				MarkdownDescription: "Defines the assignment of the plan with the quota specified in `auto_distribute_amount` to subaccounts currently located in the specified directory. For entitlements without a numeric quota, the plan is assigned to the subaccounts currently located in the directory (`auto_distribute_amount` is not needed). When applying this option, `auto_assign` must also be set.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The current state of the entitlement. Possible values are: \n " +
					getFormattedValueAsTableRow("value", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`PLATFORM`", " A service required for using a specific platform; for example, Application Runtime is required for the Cloud Foundry platform.") +
					getFormattedValueAsTableRow("`SERVICE`", "A commercial or technical service. that has a numeric quota (amount) when entitled or assigned to a resource. When assigning entitlements of this type, use the 'amount' option.") +
					getFormattedValueAsTableRow("`ELASTIC_SERVICE`", "A commercial or technical service that has no numeric quota (amount) when entitled or assigned to a resource. Generally this type of service can be as many times as needed when enabled, but may in some cases be restricted by the service owner.") +
					getFormattedValueAsTableRow("`ELASTIC_LIMITED`", "An elastic service that can be enabled for only one subaccount per global account.") +
					getFormattedValueAsTableRow("`APPLICATION`", "A multitenant application to which consumers can subscribe. As opposed to applications defined as a 'QUOTA_BASED_APPLICATION', these applications do not have a numeric quota and are simply enabled or disabled as entitlements per subaccount.") +
					getFormattedValueAsTableRow("`QUOTA_BASED_APPLICATION`", "A multitenant application to which consumers can subscribe. As opposed to applications defined as 'APPLICATION', these applications have an numeric quota that limits consumer usage of the subscribed application per subaccount.") +
					getFormattedValueAsTableRow("`ENVIRONMENT`", " An environment service; for example, Cloud Foundry."),
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"plan_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the entitled service plan.",
				Computed:            true,
			},
		},
	}
}

func (rs *directoryEntitlementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state directoryEntitlementType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entitlement, _, err := rs.cli.Accounts.Entitlement.GetEntitledByDirectory(ctx, state.DirectoryId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Entitlement (Directory)", fmt.Sprintf("%s", err))
		return
	}

	if entitlement == nil {
		resp.Diagnostics.AddError("API Error Reading Resource Entitlement (Directory)", "Resource not found")
		return
	}

	updatedState, diags := directoryEntitlementValueFrom(ctx, *entitlement, state.DirectoryId.ValueString(), state.Distribute.ValueBool())

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (rs *directoryEntitlementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	rs.createOrUpdate(ctx, req.Plan, &resp.Diagnostics, &resp.State, "Creating")
}

func (rs *directoryEntitlementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	rs.createOrUpdate(ctx, req.Plan, &resp.Diagnostics, &resp.State, "Updating")
}

func (rs *directoryEntitlementResource) createOrUpdate(ctx context.Context, requestPlan tfsdk.Plan, responseDiagnostics *diag.Diagnostics, responseState *tfsdk.State, action string) {
	var plan directoryEntitlementType
	diags := requestPlan.Get(ctx, &plan)
	responseDiagnostics.Append(diags...)
	if responseDiagnostics.HasError() {
		return
	}

	var err error
	if !hasPlanQuotaDir(plan) {
		_, err = rs.cli.Accounts.Entitlement.EnableInDirectory(ctx, plan.DirectoryId.ValueString(), plan.ServiceName.ValueString(), plan.PlanName.ValueString(), plan.Distribute.ValueBool(), plan.AutoAssign.ValueBool())
	} else {

		dirAssignmentInput := btpcli.DirectoryAssignmentInput{
			DirectoryId:          plan.DirectoryId.ValueString(),
			ServiceName:          plan.ServiceName.ValueString(),
			ServicePlanName:      plan.PlanName.ValueString(),
			Amount:               int(plan.Amount.ValueInt64()),
			Distribute:           plan.Distribute.ValueBool(),
			AutoAssign:           plan.AutoAssign.ValueBool(),
			AutoDistributeAmount: int(plan.AutoDistributeAmount.ValueInt64()),
		}
		_, err = rs.cli.Accounts.Entitlement.AssignToDirectory(ctx, dirAssignmentInput)
	}

	if err != nil {
		responseDiagnostics.AddError(fmt.Sprintf("API Error %s Resource Entitlement (Directory)", action), fmt.Sprintf("%s", err))
		return
	}

	// wait for the entitlement to become effective
	createStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis_entitlements.StateStarted, cis_entitlements.StateProcessing},
		Target:  []string{cis_entitlements.StateOK},
		Refresh: func() (interface{}, string, error) {
			entitlement, _, err := rs.cli.Accounts.Entitlement.GetEntitledByDirectory(ctx, plan.DirectoryId.ValueString(), plan.ServiceName.ValueString(), plan.PlanName.ValueString())

			if err != nil {
				return nil, "", err
			}

			if entitlement == nil {
				return nil, cis_entitlements.StateProcessing, nil
			}

			if !reflect.ValueOf(entitlement).IsNil() {
				if checkForTargetStateReached(ctx, *entitlement, plan) {
					return *entitlement, cis_entitlements.StateOK, nil
				}
			}

			return *entitlement, cis_entitlements.StateProcessing, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	entitlement, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		responseDiagnostics.AddError(fmt.Sprintf("API Error %s Resource Entitlement (Directory)", action), fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := directoryEntitlementValueFrom(ctx, entitlement.(btpcli.UnfoldedEntitlement), plan.DirectoryId.ValueString(), plan.Distribute.ValueBool())
	responseDiagnostics.Append(diags...)

	diags = responseState.Set(ctx, &updatedState)
	responseDiagnostics.Append(diags...)
}

func (rs *directoryEntitlementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state directoryEntitlementType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !hasPlanQuotaDir(state) {
		_, err = rs.cli.Accounts.Entitlement.DisableInDirectory(ctx, state.DirectoryId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString(), state.Distribute.ValueBool(), state.AutoAssign.ValueBool())
	} else {

		dirAssignmentInput := btpcli.DirectoryAssignmentInput{
			DirectoryId:          state.DirectoryId.ValueString(),
			ServiceName:          state.ServiceName.ValueString(),
			ServicePlanName:      state.PlanName.ValueString(),
			Amount:               0,
			Distribute:           state.Distribute.ValueBool(),
			AutoAssign:           state.AutoAssign.ValueBool(),
			AutoDistributeAmount: 0,
		}
		_, err = rs.cli.Accounts.Entitlement.AssignToDirectory(ctx, dirAssignmentInput)
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Entitlement (Directory)", fmt.Sprintf("%s", err))
		return
	}

	deleteStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis_entitlements.StateStarted, cis_entitlements.StateProcessing},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {

			entitlement, _, err := rs.cli.Accounts.Entitlement.GetEntitledByDirectory(ctx, state.DirectoryId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString())

			if reflect.ValueOf(entitlement).IsNil() {
				return entitlement, "DELETED", nil
			}

			if err != nil {
				return entitlement, cis_entitlements.StateProcessingFailed, err
			}

			return entitlement, cis_entitlements.StateProcessing, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = deleteStateConf.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Entitlement (Directory)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *directoryEntitlementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	//Known Gap: The DISTRIBUTE flag cannot be fetched via the platform APIs. Hence, we cannot import the value, it will always be FALSE
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: directory,service_name,plan_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("directory_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("plan_name"), idParts[2])...)
}

func hasPlanQuotaDir(state directoryEntitlementType) bool {

	// Case 1: CREATE with a explicitly non-specified amount by caller
	if state.Amount.ValueInt64() == 0 {
		return false
	}

	// Case 2: Categories that allow enabling/disabling only
	planCategory := state.Category.ValueString()
	if planCategory == "ELASTIC_SERVICE" || planCategory == "ELASTIC_LIMITED" || planCategory == "APPLICATION" {
		return false
	}

	return true
}

// As we do not get valid state information from the CLI we need to check for changes manually esp. for updates
func checkForTargetStateReached(ctx context.Context, entitlement btpcli.UnfoldedEntitlement, plan directoryEntitlementType) bool {

	updatedState, _ := directoryEntitlementValueFrom(ctx, entitlement, plan.DirectoryId.String(), plan.Distribute.ValueBool())
	// Execute the check on all fields that are potential input
	if updatedState.ServiceName.ValueString() == plan.ServiceName.ValueString() &&
		updatedState.PlanName.ValueString() == plan.PlanName.ValueString() &&
		updatedState.AutoAssign.ValueBool() == plan.AutoAssign.ValueBool() &&
		updatedState.AutoDistributeAmount.ValueInt64() == plan.AutoDistributeAmount.ValueInt64() &&
		// Special case for Amount as this is changed to 1 if plan is 0
		(updatedState.Amount.ValueInt64() == plan.Amount.ValueInt64() || plan.Amount.ValueInt64() == 0) {
		return true
	} else {
		return false
	}

}
