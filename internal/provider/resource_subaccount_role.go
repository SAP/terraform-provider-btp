package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountRoleResource() resource.Resource {
	return &subaccountRoleResource{}
}

type subaccountRoleResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role", req.ProviderTypeName)
}

func (rs *subaccountRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a role in a subaccount.

__Tip:__
You must be assigned to the admin role of the subaccount.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/role-collections-and-roles-in-global-accounts-directories-and-subaccounts>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id`, `name`, `role_template_name` and `app_id` attributes instead",
				MarkdownDescription: "The combined unique ID of the role.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the xsuaa application.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"role_template_name": schema.StringAttribute{
				MarkdownDescription: "The name of the role template.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The role description.",
				Optional:            true,
				Computed:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the role can be modified or not.",
				Computed:            true,
			},
		},
	}
}

func (rs *subaccountRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountRoleType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Security.Role.GetBySubaccount(ctx,
		state.SubaccountId.ValueString(),
		state.Name.ValueString(),
		state.RoleTemplateAppId.ValueString(),
		state.RoleTemplateName.ValueString(),
	)
	if err != nil {
		handleReadErrors(ctx, rawRes, resp, err, "Resource Role (Subaccount)")
		return
	}

	updatedState, diags := subaccountRoleFromValue(ctx, cliRes)
	updatedState.SubaccountId = state.SubaccountId

	if updatedState.Id.IsNull() || updatedState.Id.IsUnknown() {
		// Setting ID of state - required by hashicorps terraform plugin testing framework for Import. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
		updatedState.Id = types.StringValue(fmt.Sprintf("%s,%s,%s,%s", state.SubaccountId.ValueString(), state.Name.ValueString(), state.RoleTemplateName.ValueString(), state.RoleTemplateAppId.ValueString()))
	}

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountRoleType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.Role.CreateBySubaccount(ctx, &btpcli.SubaccountRoleCreateInput{
		RoleName:         plan.Name.ValueString(),
		AppId:            plan.RoleTemplateAppId.ValueString(),
		RoleTemplateName: plan.RoleTemplateName.ValueString(),
		SubaccountId:     plan.SubaccountId.ValueString(),
		Description:      plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	updatedPlan, diags := subaccountRoleFromValue(ctx, cliRes)
	updatedPlan.SubaccountId = plan.SubaccountId

	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	updatedPlan.Id = types.StringValue(fmt.Sprintf("%s,%s,%s,%s", plan.SubaccountId.ValueString(), plan.Name.ValueString(), plan.RoleTemplateName.ValueString(), plan.RoleTemplateAppId.ValueString()))

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedPlan)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountRoleType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Resource Role (Subaccount)", "This resource is not supposed to be updated")
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *subaccountRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountRoleType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.Role.DeleteBySubaccount(ctx, state.SubaccountId.ValueString(), state.Name.ValueString(), state.RoleTemplateAppId.ValueString(), state.RoleTemplateName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Role (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

// Create the function for the state import
func (rs *subaccountRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 4 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: subaccount_id, name, role_template_name, app_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_template_name"), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), idParts[3])...)
}
