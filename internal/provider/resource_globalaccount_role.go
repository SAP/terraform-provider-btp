package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountRoleResource() resource.Resource {
	return &globalaccountRoleResource{}
}

type globalaccountRoleResource struct {
	cli *btpcli.ClientFacade
}

func (rs *globalaccountRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_role", req.ProviderTypeName)
}

func (rs *globalaccountRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *globalaccountRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Create a role in a global account.

__Further documentation__
https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/0039cf082d3d43eba9200fe15647922a.html`,
		Attributes: map[string]schema.Attribute{
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
				MarkdownDescription: "Whether the role can be modified or not.",
				Computed:            true,
			},
		},
	}
}

func (rs *globalaccountRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state globalaccountRoleType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.Role.GetByGlobalAccount(ctx,
		state.Name.ValueString(),
		state.RoleTemplateAppId.ValueString(),
		state.RoleTemplateName.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := globalaccountRoleFromValue(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan globalaccountRoleType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.Role.CreateByGlobalAccount(ctx, &btpcli.GlobalAccountRoleCreateInput{
		RoleName:         plan.Name.ValueString(),
		AppId:            plan.RoleTemplateAppId.ValueString(),
		RoleTemplateName: plan.RoleTemplateName.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	updatedPlan, diags := globalaccountRoleFromValue(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedPlan)
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan globalaccountRoleType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Resource Role (Global Account)", "this resource is not supposed to be updated")
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *globalaccountRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state globalaccountRoleType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.Role.DeleteByGlobalAccount(ctx, state.Name.ValueString(), state.RoleTemplateAppId.ValueString(), state.RoleTemplateName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Role (Global Account)", fmt.Sprintf("%s", err))
		return
	}
}
