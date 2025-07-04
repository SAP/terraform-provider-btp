package provider

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

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

const entitlementCallRetrySucceeded = "retryCallSucceeded"
const entitlementCallRetryPending = "retryCallPending"
const entitlementCallRetryFailed = "retryCallFailed"

func newSubaccountEntitlementResource() resource.Resource {
	return &subaccountEntitlementResource{}
}

type subaccountEntitlementResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountEntitlementResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_entitlement", req.ProviderTypeName)
}

func (rs *subaccountEntitlementResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountEntitlementResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assigns the entitlement plan of a service, multitenant application, or environment, to a subaccount. Note that some environments, such as Cloud Foundry, are available by default to all global accounts and their subaccounts, and therefore are not made available as entitlements.

__Tip:__
You must be assigned to the admin role of the global account.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/entitlements-and-quotas>`,
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
			"plan_unique_identifier": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the service plan.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"amount": schema.Int64Attribute{
				MarkdownDescription: "The quota assigned to the subaccount.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 2000000000),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the entitlement. Possible values are: \n " +
					getFormattedValueAsTableRow("state", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`OK`", "The CRUD operation or series of operations completed successfully.") +
					getFormattedValueAsTableRow("`STARTED`", "The processing operation started") +
					getFormattedValueAsTableRow("`PROCESSING`", "The processing operation is in progress") +
					getFormattedValueAsTableRow("`PROCESSING_FAILED`", "The processing operation failed"),
				Computed: true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
		},
	}
}

type SubaccountEntitlementResourceIdentityModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	ServiceName  types.String `tfsdk:"service_name"`
	PlanName     types.String `tfsdk:"plan_name"`
}

func (rs *subaccountEntitlementResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"subaccount_id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"service_name": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"plan_name": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (rs *subaccountEntitlementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountEntitlementType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the parent of the subaccount
	// In case of a directory with feature "ENTITLEMENTS" enabled we must hand over the ID in the GetAssignedBySubaccount call
	subaccountData, _, _ := rs.cli.Accounts.Subaccount.Get(ctx, state.SubaccountId.ValueString())
	parentId, isParentGlobalAccount, err := determineParentIdForEntitlement(rs.cli, ctx, subaccountData.ParentGUID)
	if err != nil {
		resp.Diagnostics.AddError("API Error determining parent features for entitlements", fmt.Sprintf("%s", err))
		return
	}

	readStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis_entitlements.StateStarted, cis_entitlements.StateProcessing},
		Target:  []string{cis_entitlements.StateOK, cis_entitlements.StateProcessingFailed},
		Refresh: func() (interface{}, string, error) {

			entitlement, _, err := rs.cli.Accounts.Entitlement.GetAssignedBySubaccount(ctx, state.SubaccountId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString(), isParentGlobalAccount, parentId)

			if tfutils.IsRetriableErrorForEntitlement(err) {
				return nil, cis_entitlements.StateProcessing, nil
			}

			if err != nil {
				return nil, "", err
			}

			if entitlement == nil {
				return nil, cis_entitlements.StateProcessing, nil
			}

			return *entitlement, entitlement.Assignment.EntityState, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	entitlement, err := readStateConf.WaitForStateContext(ctx)

	if err != nil {
		if notFoundErr(err) {
			// Treat "Not Found" as a signal to recreate resource see https://developer.hashicorp.com/terraform/plugin/framework/resources/read#recommendations
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(("API Error Resource Entitlement (Subaccount) Read"), fmt.Sprintf("%s", err))
		return
	}

	if entitlement == nil {
		// Treat "Not Found" as a signal to recreate resource see https://developer.hashicorp.com/terraform/plugin/framework/resources/read#recommendations
		resp.State.RemoveResource(ctx)
		return
	}

	updatedState, diags := subaccountEntitlementValueFrom(ctx, entitlement.(btpcli.UnfoldedAssignment))
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)

	var identity SubaccountEntitlementResourceIdentityModel

	diags = req.Identity.Get(ctx, &identity)
	if diags.HasError() {
		identity = SubaccountEntitlementResourceIdentityModel{
			SubaccountId: types.StringValue(state.SubaccountId.ValueString()),
			ServiceName:  types.StringValue(updatedState.ServiceName.ValueString()),
			PlanName:     types.StringValue(updatedState.PlanName.ValueString()),
		}

		diags = resp.Identity.Set(ctx, identity)
		resp.Diagnostics.Append(diags...)
	}
}

func (rs *subaccountEntitlementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	rs.createOrUpdate(ctx, req.Plan, &resp.Diagnostics, &resp.State, &resp.Identity, "Creating")
}

func (rs *subaccountEntitlementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	rs.createOrUpdate(ctx, req.Plan, &resp.Diagnostics, &resp.State, &resp.Identity, "Updating")
}

func (rs *subaccountEntitlementResource) createOrUpdate(ctx context.Context, requestPlan tfsdk.Plan, responseDiagnostics *diag.Diagnostics, responseState *tfsdk.State, responseIdentity **tfsdk.ResourceIdentity, action string) {
	var plan subaccountEntitlementType
	diags := requestPlan.Get(ctx, &plan)
	responseDiagnostics.Append(diags...)
	if responseDiagnostics.HasError() {
		return
	}

	// Determine the parent of the subaccount
	subaccountData, _, _ := rs.cli.Accounts.Subaccount.Get(ctx, plan.SubaccountId.ValueString())

	//Determine if the parent is a directory and if it has authorization enabled
	parentId, isParentGlobalAccount, err := determineParentIdForAuthorization(rs.cli, ctx, subaccountData.ParentGUID)
	if err != nil {
		responseDiagnostics.AddError("API Error determining parent features for authorization", fmt.Sprintf("%s", err))
		return
	}

	var directoryId string
	if !isParentGlobalAccount {
		directoryId = parentId
	}

	// We call the API in a retry mode as the API may return a locking error
	RetryApiCallConf := &tfutils.StateChangeConf{
		Pending: []string{entitlementCallRetryPending},
		Target:  []string{entitlementCallRetryFailed, entitlementCallRetrySucceeded},
		Refresh: func() (interface{}, string, error) {

			var callResult btpcli.CommandResponse
			var err error

			if !hasPlanQuota(plan) {
				callResult, err = rs.cli.Accounts.Entitlement.EnableInSubaccount(ctx, directoryId, plan.SubaccountId.ValueString(), plan.ServiceName.ValueString(), plan.PlanName.ValueString(), plan.PlanUniqueIdentifier.ValueString())
			} else {
				callResult, err = rs.cli.Accounts.Entitlement.AssignToSubaccount(ctx, directoryId, plan.SubaccountId.ValueString(), plan.ServiceName.ValueString(), plan.PlanName.ValueString(), plan.PlanUniqueIdentifier.ValueString(), int(plan.Amount.ValueInt64()))
			}

			if err == nil {
				return callResult, entitlementCallRetrySucceeded, nil
			}

			if tfutils.IsRetriableErrorForEntitlement(err) {
				return callResult, entitlementCallRetryPending, nil
			}

			if err != nil {
				return callResult, entitlementCallRetryFailed, err
			}

			return callResult, entitlementCallRetryPending, err
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = RetryApiCallConf.WaitForStateContext(ctx)

	if err != nil {
		responseDiagnostics.AddError(fmt.Sprintf("API Error %s Resource Entitlement (Subaccount)", action), fmt.Sprintf("%s", err))
		return
	}

	// In case of a directory with feature "ENTITLEMENTS" enabled we must hand over the ID in the GetAssignedBySubaccount call
	parentId, isParentGlobalAccount, err = determineParentIdForEntitlement(rs.cli, ctx, subaccountData.ParentGUID)
	if err != nil {
		responseDiagnostics.AddError("API Error determining parent features for entitlements", fmt.Sprintf("%s", err))
		return
	}

	// wait for the entitlement to become effective
	createStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis_entitlements.StateStarted, cis_entitlements.StateProcessing},
		Target:  []string{cis_entitlements.StateOK},
		Refresh: func() (interface{}, string, error) {

			entitlement, _, err := rs.cli.Accounts.Entitlement.GetAssignedBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.ServiceName.ValueString(), plan.PlanName.ValueString(), isParentGlobalAccount, parentId)

			if tfutils.IsRetriableErrorForEntitlement(err) {
				return nil, cis_entitlements.StateProcessing, nil
			}

			if err != nil {
				return nil, "", err
			}

			if entitlement == nil {
				return nil, cis_entitlements.StateProcessing, nil
			}
			// No error returned even if operation failed
			if entitlement.Assignment.EntityState == cis_entitlements.StateProcessingFailed {
				return *entitlement, entitlement.Assignment.EntityState, errors.New("undefined API error during entitlement processing")
			}

			return *entitlement, entitlement.Assignment.EntityState, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	entitlement, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		responseDiagnostics.AddError(fmt.Sprintf("API Error %s Resource Entitlement (Subaccount)", action), fmt.Sprintf("%s", err))
		return
	}

	// The amount field is always set, even if not specified. Distinguish between operations via category
	updatedState, diags := subaccountEntitlementValueFrom(ctx, entitlement.(btpcli.UnfoldedAssignment))
	responseDiagnostics.Append(diags...)

	diags = responseState.Set(ctx, &updatedState)
	responseDiagnostics.Append(diags...)

	identity := SubaccountEntitlementResourceIdentityModel{
		SubaccountId: types.StringValue(plan.SubaccountId.ValueString()),
		ServiceName:  types.StringValue(updatedState.ServiceName.ValueString()),
		PlanName:     types.StringValue(updatedState.PlanName.ValueString()),
	}

	diags = (*responseIdentity).Set(ctx, identity)
	responseDiagnostics.Append(diags...)
}

func (rs *subaccountEntitlementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountEntitlementType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine the parent of the subaccount
	subaccountData, _, _ := rs.cli.Accounts.Subaccount.Get(ctx, state.SubaccountId.ValueString())
	//Determine if the parent is a directory and if it has authoization enabled
	parentId, isParentGlobalAccount, err := determineParentIdForAuthorization(rs.cli, ctx, subaccountData.ParentGUID)
	if err != nil {
		resp.Diagnostics.AddError("API Error determining parent features for authorization", fmt.Sprintf("%s", err))
		return
	}

	var directoryId string

	if !isParentGlobalAccount {
		directoryId = parentId
	}

	// We call the API in a retry mode as the API may return a locking error
	RetryApiCallConf := &tfutils.StateChangeConf{
		Pending: []string{entitlementCallRetryPending},
		Target:  []string{entitlementCallRetryFailed, entitlementCallRetrySucceeded},
		Refresh: func() (interface{}, string, error) {
			var callResult btpcli.CommandResponse
			var err error

			if !hasPlanQuota(state) {
				callResult, err = rs.cli.Accounts.Entitlement.DisableInSubaccount(ctx, directoryId, state.SubaccountId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString())
			} else {
				callResult, err = rs.cli.Accounts.Entitlement.AssignToSubaccount(ctx, directoryId, state.SubaccountId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString(), state.PlanUniqueIdentifier.ValueString(), 0)
			}

			if err == nil {
				return callResult, entitlementCallRetrySucceeded, nil
			}

			if tfutils.IsRetriableErrorForEntitlement(err) {
				return callResult, entitlementCallRetryPending, nil
			}

			if err != nil {
				return callResult, entitlementCallRetryFailed, err
			}

			return callResult, entitlementCallRetryPending, err
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = RetryApiCallConf.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Entitlement (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	// In case of a directory with feature "ENTITLEMENTS" enabled we must hand over the ID in the GetAssignedBySubaccount call
	parentId, isParentGlobalAccount, err = determineParentIdForEntitlement(rs.cli, ctx, subaccountData.ParentGUID)
	if err != nil {
		resp.Diagnostics.AddError("API Error determining parent features for entitlements", fmt.Sprintf("%s", err))
		return
	}

	deleteStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis_entitlements.StateStarted, cis_entitlements.StateProcessing},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {

			entitlement, _, err := rs.cli.Accounts.Entitlement.GetAssignedBySubaccount(ctx, state.SubaccountId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString(), isParentGlobalAccount, parentId)

			if reflect.ValueOf(entitlement).IsNil() {
				return entitlement, "DELETED", nil
			}

			if tfutils.IsRetriableErrorForEntitlement(err) {
				return nil, cis_entitlements.StateProcessing, nil
			}

			if err != nil {
				return entitlement, cis_entitlements.StateProcessingFailed, err
			}

			// No error returned even if operation failed
			if entitlement.Assignment.EntityState == cis_entitlements.StateProcessingFailed {
				return *entitlement, entitlement.Assignment.EntityState, errors.New(entitlement.Assignment.StateMessage)
			}

			return entitlement, cis_entitlements.StateProcessing, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = deleteStateConf.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Entitlement (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountEntitlementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" {
		idParts := strings.Split(req.ID, ",")

		if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier",
				fmt.Sprintf("Expected import identifier with format: subaccount,service_name,plan_name. Got: %q", req.ID),
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), idParts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("plan_name"), idParts[2])...)
		return
	}
	var identityData SubaccountEntitlementResourceIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), identityData.SubaccountId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), identityData.ServiceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("plan_name"), identityData.PlanName)...)
}

func hasPlanQuota(state subaccountEntitlementType) bool {

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

func notFoundErr(err error) bool {
	if err.Error() != "" && strings.Contains(err.Error(), "couldn't find resource") {
		return true
	}
	return false
}
