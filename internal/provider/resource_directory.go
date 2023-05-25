package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newDirectoryResource() resource.Resource {
	return &directoryResource{}
}

type directoryResource struct {
	cli *btpcli.ClientFacade
}

func (rs *directoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory", req.ProviderTypeName)
}

func (rs *directoryResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *directoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Directories allow you to organize and manage your subaccounts according to your technical and business needs. The use of directories is optional.

You can create up to 5 levels of directories in your account hierarchy. If you have directories, you can still create subaccounts directly under your global account.

Directory features: Set the '--features' parameter to specify which features to enable for the directory. Use either the feature name or its character.
* DEFAULT (D): (Required) All directories provide the following basic features: (i) Group and filter subaccounts for reports and filters, (ii) monitor usage and costs on a directory level (costs only available for contracts that use the consumption-based commercial model), and (iii) assign labels to the directory for identification and reporting purposes.
* ENTITLEMENTS (E): (Optional) Enables the assignment of a quota for services and applications to the directory from the global account quota for distribution to the directory's subaccounts and subdirectories.
* AUTHORIZATIONS (A): (Optional) Allows you to assign users as administrators or viewers of this directory. For example, allow certain users to manage the directory's entitlements. You can enable this feature only in combination with the ENTITLEMENTS (E) feature.

__Tips__
* You must be assigned to the global account admin role, or the directory admin if the directory is configured to manage its authorizations.
* A directory path in the account hierarchy can have only one directory that is enabled with the ENTITLEMENTS (E) or AUTHORIZATIONS (A) features. If such a directory exists, other directories in that path can only be enabled with the DEFAULT (D) features.

__Further documentation__
https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/8ed4a705efa0431b910056c0acdbf377.html`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The display name of the directory.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[^\/]{1,255}$`), "must not contain '/', not be empty and not exceed 255 characters"),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the directory.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(300),
				},
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "The GUID of the directory's parent entity. Typically this is the global account.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "Applies only to directories that have the user authorization management feature enabled. The subdomain becomes part of the path used to access the authorization tenant of the directory. Unique within the defined region.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("^[a-z0-9](?:[a-z0-9|-]{0,61}[a-z0-9])?$"), "must only contain letters (a-z), digits (0-9), and hyphens (not at the start or end)"),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "Details of the user that created the directory.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},

			"features": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The features that are enabled for the directory. Valid values: - DEFAULT: (Mandatory) All directories have the following basic feature enabled. (1) Group and filter subaccounts for reports and filters, (2) monitor usage and costs on a directory level (costs only available for contracts that use the consumption-based commercial model), and (3) set custom properties and tags to the directory for identification and reporting purposes. - ENTITLEMENTS: (Optional) Allows the assignment of a quota for services and applications to the directory from the global account quota for distribution to the subaccounts under this directory.  - AUTHORIZATIONS: (Optional) Allows the assignment of users as administrators or viewers of this directory. You must apply this feature in combination with the ENTITLEMENTS feature. <br/><b>Valid values:</b>  [DEFAULT] [DEFAULT,ENTITLEMENTS] [DEFAULT,ENTITLEMENTS,AUTHORIZATIONS]<br/>",
				Computed:            true,
			},
			"labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				MarkdownDescription: "Contains information about the labels assigned to a specified global account. Labels are represented in a JSON array of key-value pairs; each key has up to 10 corresponding values. This field replaces the deprecated \"customProperties\" field, which supports only single values per key.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},

			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the directory. * <b>STARTED:</b> CRUD operation on an entity has started. * <b>CREATING:</b> Creating entity operation is in progress. * <b>UPDATING:</b> Updating entity operation is in progress. * <b>MOVING:</b> Moving entity operation is in progress. * <b>PROCESSING:</b> A series of operations related to the entity is in progress. * <b>DELETING:</b> Deleting entity operation is in progress. * <b>OK:</b> The CRUD operation or series of operations completed successfully. * <b>PENDING REVIEW:</b> The processing operation has been stopped for reviewing and can be restarted by the operator. * <b>CANCELLED:</b> The operation or processing was canceled by the operator. * <b>CREATION_FAILED:</b> The creation operation failed, and the entity was not created or was created but cannot be used. * <b>UPDATE_FAILED:</b> The update operation failed, and the entity was not updated. * <b>PROCESSING_FAILED:</b> The processing operations failed. * <b>DELETION_FAILED:</b> The delete operation failed, and the entity was not deleted. * <b>MOVE_FAILED:</b> Entity could not be moved to a different location. * <b>MIGRATING:</b> Migrating entity from NEO to CF.",
				Computed:            true,
			},
		},
	}
}

func (rs *directoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state directoryType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Accounts.Directory.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Directory", fmt.Sprintf("%s", err))
		return
	}

	state, diags = directoryValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *directoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan directoryType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	args := btpcli.DirectoryCreateInput{
		DisplayName: plan.Name.ValueString(),
	}

	if !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		args.Description = &description
	}

	if !plan.ParentID.IsUnknown() {
		parentID := plan.ParentID.ValueString()
		args.ParentID = &parentID
	}

	if !plan.Subdomain.IsUnknown() {
		subdomain := plan.Subdomain.ValueString()
		args.Subdomain = &subdomain
	}

	cliRes, _, err := rs.cli.Accounts.Directory.Create(ctx, &args)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Directory", fmt.Sprintf("%s", err))
		return
	}

	plan, diags = directoryValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	createStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis.StateCreating, cis.StateStarted},
		Target:  []string{cis.StateOK, cis.StateCreationFailed, cis.StateCanceled},
		Refresh: func() (interface{}, string, error) {
			subRes, _, err := rs.cli.Accounts.Directory.Get(ctx, cliRes.Guid)

			if err != nil {
				return subRes, "", err
			}

			return subRes, subRes.EntityState, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	updatedRes, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Directory", fmt.Sprintf("%s", err))
	}

	plan, diags = directoryValueFrom(ctx, updatedRes.(cis.DirectoryResponseObject))
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (rs *directoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan directoryType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Resource Directory", "Update is not yet implemented.")

	/*TODO: cliRes, err := rs.cli.Execute(ctx, btpcli.Update, rs.command, plan)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Directory", fmt.Sprintf("%s", err))
		return
	}*/

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *directoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state directoryType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Accounts.Directory.Delete(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Directory", fmt.Sprintf("%s", err))
		return
	}

	deleteStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis.StateDeleting, cis.StateStarted},
		Target:  []string{cis.StateOK, cis.StateDeletionFailed, cis.StateCanceled, "DELETED"},
		Refresh: func() (interface{}, string, error) {
			subRes, comRes, err := rs.cli.Accounts.Directory.Get(ctx, cliRes.Guid)

			if err != nil {
				return subRes, subRes.EntityState, err
			}

			if comRes.StatusCode == http.StatusNotFound || comRes.StatusCode == http.StatusForbidden {
				return subRes, "DELETED", nil
			}

			return subRes, subRes.EntityState, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = deleteStateConf.WaitForStateContext(ctx)

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Directory", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *directoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
