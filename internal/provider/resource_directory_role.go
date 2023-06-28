package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

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
			"id": schema.StringAttribute{ // required hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `directory_id` attribute instead",
				MarkdownDescription: "The ID of the directory.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role.",
				Required:            true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the xsuaa application.",
				Required:            true,
			},
			"role_template_name": schema.StringAttribute{
				MarkdownDescription: "The name of the role template.",
				Required:            true,
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

			// TODO: this was commented out for now as user has no direct influence
			//"scopes": schema.SetNestedAttribute{
			//	NestedObject: schema.NestedAttributeObject{
			//		Attributes: map[string]schema.Attribute{
			//			"name": schema.StringAttribute{
			//				MarkdownDescription: "The name of the scope.",
			//				Computed:    true,
			//			},
			//			"description": schema.StringAttribute{
			//				MarkdownDescription: "The description of the scope.",
			//				Computed:    true,
			//			},
			//			"custom_grant_as_authority_to_apps": schema.SetAttribute{
			//				ElementType: types.StringType,
			//				Computed:    true,
			//			},
			//			"custom_granted_apps": schema.SetAttribute{
			//				ElementType: types.StringType,
			//				Computed:    true,
			//			},
			//			"grant_as_authority_to_apps": schema.SetAttribute{
			//				ElementType: types.StringType,
			//				Computed:    true,
			//			},
			//			"granted_apps": schema.SetAttribute{
			//				ElementType: types.StringType,
			//				Computed:    true,
			//			},
			//		},
			//	},
			//	MarkdownDescription: "Scopes available with this role.",
			//	Computed:    true,
			//},
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
	updatedState.Id = state.DirectoryId
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
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role (Directory)", fmt.Sprintf("%s", err))
		return
	}

	updatedPlan, diags := directoryRoleFromValue(ctx, cliRes)
	updatedPlan.DirectoryId = plan.DirectoryId
	updatedPlan.Id = plan.DirectoryId
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
