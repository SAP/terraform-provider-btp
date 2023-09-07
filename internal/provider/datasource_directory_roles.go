package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newDirectoryRolesDataSource() datasource.DataSource {
	return &directoryRolesDataSource{}
}

type directoryRolesValue struct {
	Name              types.String `tfsdk:"name"`
	RoleTemplateAppId types.String `tfsdk:"app_id"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	/* OUTPUT */
	Description types.String         `tfsdk:"description"`
	IsReadOnly  types.Bool           `tfsdk:"read_only"`
	AppName     types.String         `tfsdk:"app_name"`
	Scopes      []directoryRoleScope `tfsdk:"scopes"`
}

type directoryRolesDataSourceConfig struct {
	/* INPUT */
	DirectoryId types.String `tfsdk:"directory_id"`
	Id          types.String `tfsdk:"id"`
	/* OUTPUT */
	Values []directoryRolesValue `tfsdk:"values"`
}

type directoryRolesDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *directoryRolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_roles", req.ProviderTypeName)
}

func (ds *directoryRolesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *directoryRolesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets all roles.

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
				DeprecationMessage:  "Use the `directory_id` attribute instead",
				MarkdownDescription: "The ID of the directory.",
				Computed:            true,
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the role.",
							Computed:            true,
						},
						"app_id": schema.StringAttribute{
							MarkdownDescription: "The id of the application that provides the role template and the role.",
							Computed:            true,
						},
						"role_template_name": schema.StringAttribute{
							MarkdownDescription: "The name of the role template.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the role.",
							Computed:            true,
						},
						"read_only": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the role can be modified or not.",
							Computed:            true,
						},
						"app_name": schema.StringAttribute{
							MarkdownDescription: "The name of the application that provides the role template and the role.",
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
				},
				Computed: true,
			},
		},
	}
}

func (ds *directoryRolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryRolesDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Role.ListByDirectory(ctx, data.DirectoryId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Roles (Directory)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.DirectoryId
	data.Values = []directoryRolesValue{}

	for _, role := range cliRes {
		roleVal := directoryRolesValue{
			Name:              types.StringValue(role.Name),
			RoleTemplateAppId: types.StringValue(role.RoleTemplateAppId),
			RoleTemplateName:  types.StringValue(role.RoleTemplateName),
			Description:       types.StringValue(role.Description),
			IsReadOnly:        types.BoolValue(role.IsReadOnly),
			AppName:           types.StringValue(role.AppName),
			Scopes:            []directoryRoleScope{},
		}

		for _, scope := range role.Scopes {
			scopeVal := directoryRoleScope{
				Name:        types.StringValue(scope.Name),
				Description: types.StringValue(scope.Description),
			}

			scopeVal.CustomGrantAsAuthorityToApps, diags = types.SetValueFrom(ctx, types.StringType, scope.CustomGrantAsAuthorityToApps)
			resp.Diagnostics.Append(diags...)

			scopeVal.CustomGrantedApps, diags = types.SetValueFrom(ctx, types.StringType, scope.CustomGrantedApps)
			resp.Diagnostics.Append(diags...)

			scopeVal.GrantAsAuthorityToApps, diags = types.SetValueFrom(ctx, types.StringType, scope.GrantAsAuthorityToApps)
			resp.Diagnostics.Append(diags...)

			scopeVal.GrantedApps, diags = types.SetValueFrom(ctx, types.StringType, scope.GrantedApps)
			resp.Diagnostics.Append(diags...)

			roleVal.Scopes = append(roleVal.Scopes, scopeVal)
		}

		data.Values = append(data.Values, roleVal)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
