package provider

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/SAP/terraform-provider-btp/internal/validation/labelvalidator"
)

func newSubaccountResource() resource.Resource {
	return &subaccountResource{}
}

type subaccountResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount", req.ProviderTypeName)
}

func (rs *subaccountResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a subaccount in a global account or directory.

__Tip:__
You must be assigned admin role of the global account or directory.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/account-model>`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "A descriptive name of the subaccount for customer-facing UIs.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^\/]{1,255}$`), "must not contain '/', not be empty and not exceed 255 characters"),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the subaccount for customer-facing UIs.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(300),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region in which the subaccount was created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "The subdomain that becomes part of the path used to access the authorization tenant of the subaccount. Must be unique within the defined region and cannot be changed after the subaccount has been created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount’s parent entity. If the subaccount is located directly in the global account (not in a directory), then this is the ID of the global account.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				MarkdownDescription: "The set of words or phrases assigned to the subaccount.",
				Optional:            true,
				Validators: []validator.Map{
					labelvalidator.ValidLabels(),
				},
			},
			"beta_enabled": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the subaccount can use beta services and applications.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_auto_entitlement": schema.BoolAttribute{
				MarkdownDescription: "Specifies if the subaccount creation excludes the auto-assignment of base entitlements, allowing quicker setup with potentially reduced resource consumption. When not set or set to 'false' the standard auto-assigned plans are included.",
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The details of the user that created the subaccount.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"parent_features": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The features of parent entity of the subaccount.",
				Computed:            true,
				// No plan modifier as field can change due to changes on parent level
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the subaccount. Possible values are: \n" +
					getFormattedValueAsTableRow("state", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`OK`", "The CRUD operation or series of operations completed successfully.") +
					getFormattedValueAsTableRow("`STARTED`", "CRUD operation on the subaccount has started.") +
					getFormattedValueAsTableRow("`CANCELED`", "The operation or processing was canceled by the operator.") +
					getFormattedValueAsTableRow("`PROCESSING`", "A series of operations related to the subaccount are in progress.") +
					getFormattedValueAsTableRow("`PROCESSING_FAILED`", "The processing operations failed.") +
					getFormattedValueAsTableRow("`CREATING`", "Creating the subaccount is in progress.") +
					getFormattedValueAsTableRow("`CREATION_FAILED`", "The creation operation failed, and the subaccount was not created or was created but cannot be used.") +
					getFormattedValueAsTableRow("`UPDATING`", "Updating the subaccount is in progress.") +
					getFormattedValueAsTableRow("`UPDATE_FAILED`", "The update operation failed, and the subaccount was not updated.") +
					getFormattedValueAsTableRow("`UPDATE_DIRECTORY_TYPE_FAILED`", "The update of the directory type failed.") +
					getFormattedValueAsTableRow("`UPDATE_ACCOUNT_TYPE_FAILED`", "The update of the account type failed.") +
					getFormattedValueAsTableRow("`DELETING`", "Deleting the subaccount is in progress.") +
					getFormattedValueAsTableRow("`DELETION_FAILED`", "The deletion of the subaccount failed, and the subaccount was not deleted.") +
					getFormattedValueAsTableRow("`MOVING`", "Moving the subaccount is in progress.") +
					getFormattedValueAsTableRow("`MOVE_FAILED`", "The moving of the subaccount failed.") +
					getFormattedValueAsTableRow("`MOVING_TO_OTHER_GA`", "Moving the subaccount to another global account is in progress.") +
					getFormattedValueAsTableRow("`MOVE_TO_OTHER_GA_FAILED`", "Moving the subaccount to another global account failed.") +
					getFormattedValueAsTableRow("`PENDING_REVIEW`", "The processing operation has been stopped for reviewing and can be restarted by the operator.") +
					getFormattedValueAsTableRow("`MIGRATING`", "Migrating the subaccount from Neo to Cloud Foundry.") +
					getFormattedValueAsTableRow("`MIGRATED`", "The migration of the subaccount completed.") +
					getFormattedValueAsTableRow("`MIGRATION_FAILED`", "The migration of the subaccount failed and the subaccount was not migrated.") +
					getFormattedValueAsTableRow("`ROLLBACK_MIGRATION_PROCESSING`", "The migration of the subaccount was rolled back and the subaccount is not migrated.") +
					getFormattedValueAsTableRow("`SUSPENSION_FAILED`", "The suspension operations failed."),
				Computed: true,
			},
			"usage": schema.StringAttribute{
				MarkdownDescription: "Shows whether the subaccount is used for production purposes. This flag can help your cloud operator to take appropriate action when handling incidents that are related to mission-critical accounts in production systems. Do not apply for subaccounts that are used for nonproduction purposes, such as development, testing, and demos. Applying this setting this does not modify the subaccount. Possible values are: \n" +
					getFormattedValueAsTableRow("value", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`UNSET`", "Global account or subaccount admin has not set the production-relevancy flag (default value).") +
					getFormattedValueAsTableRow("`NOT_USED_FOR_PRODUCTION`", "The subaccount is not used for production purposes.") +
					getFormattedValueAsTableRow("`USED_FOR_PRODUCTION`", "The subaccount is used for production purposes."),
				Computed: true,
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						"USED_FOR_PRODUCTION",
						"NOT_USED_FOR_PRODUCTION",
						"UNSET",
					}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type subaccountResourceIdentityModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
}

func (rs *subaccountResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"subaccount_id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (rs *subaccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data subaccountType

	diags := req.State.Get(ctx, &data)

	// As we manipulate the state afterwards we make a copy for later use
	originalState := data

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Accounts.Subaccount.Get(ctx, data.ID.ValueString())
	if err != nil {
		handleReadErrors(ctx, rawRes, cliRes, resp, err, "Resource Subaccount")
		return
	}

	data, diags = subaccountValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	if !originalState.SkipAutoEntitlement.IsUnknown() {
		// The value was explicitly set in the state and is not returned via the READ API call, so we keep it
		data.SkipAutoEntitlement = originalState.SkipAutoEntitlement
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	var identity subaccountResourceIdentityModel

	diags = req.Identity.Get(ctx, &identity)
	if diags.HasError() {
		identity = subaccountResourceIdentityModel{
			SubaccountId: data.ID,
		}

		diags = resp.Identity.Set(ctx, identity)
		resp.Diagnostics.Append(diags...)
	}
}

func (rs *subaccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// As we manipulate the plan afterwards we make a copy for later use
	originalPlan := plan

	// We check the subdomain value against the recommended format of the API and the BTP Administrator's guide
	// However we cannot enforce it via the schema as the cockpit allows also deviating formats

	if !subdomainRegex.MatchString(plan.Subdomain.ValueString()) {
		resp.Diagnostics.AddAttributeWarning(path.Root("subdomain"), "subdomain is in a non-recommended format", "subdomain should only contain letters (a-z), digits (0-9), and hyphens (not at the start or end)")
	}

	args := btpcli.SubaccountCreateInput{
		DisplayName: plan.Name.ValueString(),
		Subdomain:   plan.Subdomain.ValueString(),
		Region:      plan.Region.ValueString(),
	}

	if !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		args.Description = description
	}

	if !plan.ParentID.IsUnknown() {
		parentID := plan.ParentID.ValueString()
		args.Directory = parentID

		//Check which parent ID needs to be used for authorization
		parentId, isParentGlobalAccount, err := determineParentIdForAuthorization(rs.cli, ctx, parentID)
		if err != nil {
			resp.Diagnostics.AddError("API Error determining parent features for authorization", fmt.Sprintf("%s", err))
			return
		}

		if !isParentGlobalAccount && parentId != "" {
			//if the parent is a managed directory, the directoryId must be set to make sure the right authorizations are validated
			args.AdminDirectoryId = parentId
		}
	}

	if !plan.BetaEnabled.IsUnknown() {
		betaEnabled := plan.BetaEnabled.ValueBool()
		args.BetaEnabled = betaEnabled
	}

	if !plan.SkipAutoEntitlement.IsUnknown() {
		skipAutoEntitlement := plan.SkipAutoEntitlement.ValueBool()
		args.SkipAutoEntitlement = skipAutoEntitlement
	}

	var labels map[string][]string
	plan.Labels.ElementsAs(ctx, &labels, false)
	args.Labels = map[string][]string{}
	maps.Copy(args.Labels, labels)

	args.UsedForProduction = mapUsageToUsedForProduction(plan.Usage.ValueString())

	cliRes, _, err := rs.cli.Accounts.Subaccount.Create(ctx, &args)

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Subaccount", fmt.Sprintf("%s", err))
		return
	}

	plan, diags = subaccountValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	// The retryable HTTP client already handles transient network and HTTP errors.
	// However, the BTP API may still respond with "not ready" or "processing" errors after a successful request.
	// Keeping this check ensures Terraform continues polling until the resource reaches a stable state.
	createStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis.StateCreating, cis.StateStarted},
		Target:  []string{cis.StateOK, cis.StateCreationFailed, cis.StateCanceled},
		Refresh: func() (any, string, error) {
			subRes, _, err := rs.cli.Accounts.Subaccount.Get(ctx, cliRes.Guid)

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
		updatedRes = cliRes
		resp.Diagnostics.AddError("API Error Creating Resource Subaccount", fmt.Sprintf("%s", err))
	}

	plan, diags = subaccountValueFrom(ctx, updatedRes.(cis.SubaccountResponseObject))

	if !originalPlan.SkipAutoEntitlement.IsUnknown() {
		// The value was explicitly set in the plan, so we keep it
		plan.SkipAutoEntitlement = originalPlan.SkipAutoEntitlement
	}

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	identity := subaccountResourceIdentityModel{
		SubaccountId: plan.ID,
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountType
	var state subaccountType

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	args := btpcli.SubaccountUpdateInput{
		BetaEnabled:  plan.BetaEnabled.ValueBool(),
		Description:  plan.Description.ValueString(),
		Directory:    plan.ParentID.ValueString(),
		DisplayName:  plan.Name.ValueString(),
		SubaccountId: plan.ID.ValueString(),
	}

	var labels map[string][]string
	plan.Labels.ElementsAs(ctx, &labels, false)
	args.Labels = map[string][]string{}
	maps.Copy(args.Labels, labels)

	args.UsedForProduction = mapUsageToUsedForProduction(plan.Usage.ValueString())

	cliRes, _, err := rs.cli.Accounts.Subaccount.Update(ctx, &args)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Subaccount", fmt.Sprintf("%s", err))
		return
	}

	plan, diags = subaccountValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	// The retryable HTTP client already handles transient network and HTTP errors.
	// However, the BTP API may still respond with "not ready" or "processing" errors after a successful request.
	// Keeping this check ensures Terraform continues polling until the resource reaches a stable state.
	updateStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis.StateUpdating, cis.StateStarted},
		Target:  []string{cis.StateOK, cis.StateUpdateFailed, cis.StateCanceled},
		Refresh: func() (any, string, error) {
			subRes, _, err := rs.cli.Accounts.Subaccount.Get(ctx, cliRes.Guid)

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
		resp.Diagnostics.AddError("API Error Updating Resource Subaccount", fmt.Sprintf("%s", err))
	}

	plan, diags = subaccountValueFrom(ctx, updatedRes.(cis.SubaccountResponseObject))

	if !state.SkipAutoEntitlement.IsUnknown() {
		// The value was explicitly set in the state, so we keep it
		plan.SkipAutoEntitlement = state.SkipAutoEntitlement
	}

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentId, isParentGlobalAccount, err := determineParentIdForAuthorization(rs.cli, ctx, state.ParentID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error determining parent features for authorization", fmt.Sprintf("%s", err))
		return
	}

	var directoryId string

	if !isParentGlobalAccount {
		//if the parent is a managed directory, the directoryId must be set to make sure the right authorizations are validated
		directoryId = parentId
	}

	cliRes, _, err := rs.cli.Accounts.Subaccount.Delete(ctx, state.ID.ValueString(), directoryId)
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Subaccount", fmt.Sprintf("%s", err))
		return
	}

	// The retryable HTTP client already handles transient network and HTTP errors.
	// However, the BTP API may still respond with "not ready" or "processing" errors after a successful request.
	// Keeping this check ensures Terraform continues polling until the resource reaches a stable state.
	deleteStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis.StateDeleting, cis.StateStarted},
		Target:  []string{cis.StateOK, cis.StateDeletionFailed, cis.StateCanceled, "DELETED"},
		Refresh: func() (any, string, error) {
			subRes, comRes, err := rs.cli.Accounts.Subaccount.Get(ctx, cliRes.Guid)

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
		resp.Diagnostics.AddError("API Error Deleting Resource Subaccount", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" {
		resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
		return
	}
	resource.ImportStatePassthroughWithIdentity(ctx, path.Root("id"), path.Root("subaccount_id"), req, resp)
}

func mapUsageToUsedForProduction(subaccountUsage string) string {
	// The BTP CLI and CIS use different parameters for the subaccount usage
	// To trigger the right usage creation in CREATE and avoid unwanted state changes in UPDATE  we must distinguish if and how to set the value
	// Options: "" == ignored in request, "true" == boolean true in request, "false" == boolean false in request
	switch subaccountUsage {
	case "UNSET":
		return ""
	case "NOT_USED_FOR_PRODUCTION":
		return "false"
	case "USED_FOR_PRODUCTION":
		return "true"
	}
	return ""
}
