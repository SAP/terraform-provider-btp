package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newDirectoryRoleDataSource() datasource.DataSource {
	return &directoryRoleDataSource{}
}

type directoryRoleDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *directoryRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_role", req.ProviderTypeName)
}

func (ds *directoryRoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *directoryRoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific directory role.
		
__Tip:__
You must be assigned to the admin or viewer role of the global account, directory.`,
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `directory_id` attribute instead",
				MarkdownDescription: "The ID of the directory.",
				Computed:            true,
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
				MarkdownDescription: "The description of the role.",
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
				MarkdownDescription: "The scopes available with this role.",
				Computed:            true,
			},
		},
	}
}

func (ds *directoryRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryRoleType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Role.GetByDirectory(ctx, data.DirectoryId.ValueString(), data.Name.ValueString(), data.RoleTemplateAppId.ValueString(), data.RoleTemplateName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role (Directory)", fmt.Sprintf("%s", err))
		return
	}

	state, diags := directoryRoleFromValue(ctx, cliRes)
	state.DirectoryId = data.DirectoryId
	state.Id = data.DirectoryId

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
