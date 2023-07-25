package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountRoleCollectionAssignmentResource() resource.Resource {
	return &subaccountRoleCollectionAssignmentResource{}
}

type subaccountRoleCollectionAssignmentType struct {
	SubaccountId       types.String `tfsdk:"subaccount_id"`
	Id                 types.String `tfsdk:"id"`
	RoleCollectionName types.String `tfsdk:"role_collection_name"`
	Username           types.String `tfsdk:"user_name"`
	Groupname          types.String `tfsdk:"group_name"`
	Origin             types.String `tfsdk:"origin"`
}

type subaccountRoleCollectionAssignmentResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountRoleCollectionAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role_collection_assignment", req.ProviderTypeName)
}

func (rs *subaccountRoleCollectionAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountRoleCollectionAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assigns a user to a role collection on a subaccount level.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role_collection_name": schema.StringAttribute{
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` and `role_collection_name` attributes instead",
				MarkdownDescription: "The combined unique ID of the role collection.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_name": schema.StringAttribute{
				MarkdownDescription: "The username of the user to assign.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("user_name"), path.MatchRoot("group_name")),
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"group_name": schema.StringAttribute{
				MarkdownDescription: "The name of the group to assign.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The identity provider that hosts the user or a group. The default value is `ldap`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("ldap"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (rs *subaccountRoleCollectionAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountRoleCollectionAssignmentType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// This resource is not supposed to be read by definition. However nothing the user can do about that, hence no error message is raised via resp.Diagnostics.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountRoleCollectionAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountRoleCollectionAssignmentType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !plan.Username.IsNull() {
		// assign user
		_, _, err = rs.cli.Security.RoleCollection.AssignUserBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.RoleCollectionName.ValueString(), plan.Username.ValueString(), plan.Origin.ValueString())
	} else {
		// assign group
		_, _, err = rs.cli.Security.RoleCollection.AssignGroupBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.RoleCollectionName.ValueString(), plan.Groupname.ValueString(), plan.Origin.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role Collection Assignment (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	plan.Id = types.StringValue(fmt.Sprintf("%s,%s,%s", plan.SubaccountId.ValueString(), plan.RoleCollectionName.ValueString(), plan.Username.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountRoleCollectionAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountRoleCollectionAssignmentType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// since all the attributes are marked to be replaced in case of update, this should never be reached.
	resp.Diagnostics.AddError("API Error Updating Role Collection Assignment (Subaccount)", "This resource is not supposed to be updated")

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *subaccountRoleCollectionAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountRoleCollectionAssignmentType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if !state.Username.IsNull() {
		// unassign user
		_, _, err = rs.cli.Security.RoleCollection.UnassignUserBySubaccount(ctx, state.SubaccountId.ValueString(), state.RoleCollectionName.ValueString(), state.Username.ValueString(), state.Origin.ValueString())
	} else {
		// unassign group
		_, _, err = rs.cli.Security.RoleCollection.UnassignGroupBySubaccount(ctx, state.SubaccountId.ValueString(), state.RoleCollectionName.ValueString(), state.Groupname.ValueString(), state.Origin.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Role Collection Assignment (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountRoleCollectionAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError(
		"Import Not Supported",
		"Import is not supported for this resource. Use the resource subaccount_role_collection instead.",
	)
}
