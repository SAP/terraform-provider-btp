package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newDirectoryRoleCollectionResource() resource.Resource {
	return &directoryRoleCollectionType{}
}

type directoryRoleCollectionRoleRefType struct {
	Name              types.String `tfsdk:"name"`
	RoleTemplateAppId types.String `tfsdk:"role_template_app_id"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
}

type directoryRoleCollectionTypeConfig struct {
	DirectoryId types.String                         `tfsdk:"directory_id"`
	Name        types.String                         `tfsdk:"name"`
	Description types.String                         `tfsdk:"description"`
	Roles       []directoryRoleCollectionRoleRefType `tfsdk:"roles"`
}

type directoryRoleCollectionType struct {
	cli *btpcli.ClientFacade
}

func (rs *directoryRoleCollectionType) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_role_collection", req.ProviderTypeName)
}

func (rs *directoryRoleCollectionType) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *directoryRoleCollectionType) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Create a role collection in a directory.

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
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Whether the role collection is readonly.",
				Optional:            true,
				Computed:            true,
			},
			"roles": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced role.",
							Required:            true,
						},
						"role_template_app_id": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced template app id",
							Required:            true,
						},
						"role_template_name": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced role template.",
							Required:            true,
						},
					},
				},
				Required: true,
			},
		},
	}
}

func (rs *directoryRoleCollectionType) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state directoryRoleCollectionTypeConfig

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.GetByDirectory(ctx, state.DirectoryId.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Directory)", fmt.Sprintf("%s", err))
		return
	}

	state.Name = types.StringValue(cliRes.Name)
	state.Description = types.StringValue(cliRes.Description)

	state.Roles = []directoryRoleCollectionRoleRefType{}
	for _, role := range cliRes.RoleReferences {
		state.Roles = append(state.Roles, directoryRoleCollectionRoleRefType{
			RoleTemplateName:  types.StringValue(role.RoleTemplateName),
			RoleTemplateAppId: types.StringValue(role.RoleTemplateAppId),
			Name:              types.StringValue(role.Name),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *directoryRoleCollectionType) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan directoryRoleCollectionTypeConfig
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.CreateByDirectory(ctx, plan.DirectoryId.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role Collection (Directory)", fmt.Sprintf("%s", err))
		return
	}

	for _, role := range plan.Roles {
		_, err := rs.cli.Security.Role.AddByDirectory(ctx, plan.DirectoryId.ValueString(), plan.Name.ValueString(), role.Name.ValueString(), role.RoleTemplateAppId.ValueString(), role.RoleTemplateName.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("API Error Assigning Role To Role Collection (Directory)", fmt.Sprintf("%s", err))
		}
	}

	plan.Name = types.StringValue(cliRes.Name)
	plan.Description = types.StringValue(cliRes.Description)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (rs *directoryRoleCollectionType) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan directoryRoleCollectionTypeConfig
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Error Updating Resource Role Collection (Directory)", "Update is not yet implemented.")

	/*TODO cliRes, err := rs.cli.Execute(ctx, btpcli.Update, rs.command, plan)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Role Collection (Directory)", fmt.Sprintf("%s", err))
		return
	}*/

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *directoryRoleCollectionType) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state directoryRoleCollectionTypeConfig
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.RoleCollection.DeleteByDirectory(ctx, state.DirectoryId.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Role Collection (Directory)", fmt.Sprintf("%s", err))
		return
	}
}
