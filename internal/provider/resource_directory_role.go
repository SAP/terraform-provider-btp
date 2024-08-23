package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

var directoryScopeObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"custom_grant_as_authority_to_apps": types.SetType{
			ElemType: types.StringType,
		},
		"custom_granted_apps": types.SetType{
			ElemType: types.StringType,
		},
		"grant_as_authority_to_apps": types.SetType{
			ElemType: types.StringType,
		},
		"granted_apps": types.SetType{
			ElemType: types.StringType,
		},
	},
}

func newDirectoryRoleResource() resource.Resource {
	return &directoryRoleResource{}
}

type directoryRoleResource struct {
	cli *btpcli.ClientFacade
}

func (rs *directoryRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_role", req.ProviderTypeName)
}

func (rs *directoryRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *directoryRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a role in a directory.

__Tip:__
You must be assigned to the admin role of the global account or the directory.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/role-collections-and-roles-in-global-accounts-directories-and-subaccounts>`,
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `directory_id`, `name`, `role_template_name` and `app_id` attributes instead",
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
			"scopes": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the scope.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the scope.",
							Computed:            true,
						},
						"custom_grant_as_authority_to_apps": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"custom_granted_apps": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"grant_as_authority_to_apps": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"granted_apps": schema.SetAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
					},
				},
				MarkdownDescription: "Scopes available with this role.",
				Computed:            true,
			},
		},
	}
}

func (rs *directoryRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state directoryRoleType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.Role.GetByDirectory(ctx,
		state.DirectoryId.ValueString(),
		state.Name.ValueString(),
		state.RoleTemplateAppId.ValueString(),
		state.RoleTemplateName.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role (Directory)", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := directoryRoleFromValue(ctx, cliRes)
	updatedState.DirectoryId = state.DirectoryId

	if updatedState.Id.IsNull() || updatedState.Id.IsUnknown() {
		// Setting ID of state - required by hashicorps terraform plugin testing framework for Import. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
		updatedState.Id = types.StringValue(fmt.Sprintf("%s,%s,%s,%s", state.DirectoryId.ValueString(), state.Name.ValueString(), state.RoleTemplateName.ValueString(), state.RoleTemplateAppId.ValueString()))
	}

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (rs *directoryRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan directoryRoleType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.Role.CreateByDirectory(ctx, &btpcli.DirectoryRoleCreateInput{
		RoleName:         plan.Name.ValueString(),
		AppId:            plan.RoleTemplateAppId.ValueString(),
		RoleTemplateName: plan.RoleTemplateName.ValueString(),
		DirectoryId:      plan.DirectoryId.ValueString(),
		Description:      plan.Description.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role (Directory)", fmt.Sprintf("%s", err))
		return
	}

	updatedPlan, diags := directoryRoleFromValue(ctx, cliRes)
	updatedPlan.DirectoryId = plan.DirectoryId

	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	updatedPlan.Id = types.StringValue(fmt.Sprintf("%s,%s,%s,%s", plan.DirectoryId.ValueString(), plan.Name.ValueString(), plan.RoleTemplateName.ValueString(), plan.RoleTemplateAppId.ValueString()))

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedPlan)
	resp.Diagnostics.Append(diags...)
}

func (rs *directoryRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan directoryRoleType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Resource Role (Directory)", "This resource is not supposed to be updated")
	if resp.Diagnostics.HasError() {
		return
	}

}

func (rs *directoryRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state directoryRoleType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.Role.DeleteByDirectory(ctx, state.DirectoryId.ValueString(), state.Name.ValueString(), state.RoleTemplateAppId.ValueString(), state.RoleTemplateName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Role (Directory)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *directoryRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 4 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: directory_id, name, role_template_name, app_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("directory_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_template_name"), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("app_id"), idParts[3])...)
}
