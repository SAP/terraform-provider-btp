package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountRolesDataSource() datasource.DataSource {
	return &globalaccountRolesDataSource{}
}

type globalaccountRolesValue struct {
	Name              types.String `tfsdk:"name"`
	RoleTemplateAppId types.String `tfsdk:"app_id"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	/* OUTPUT */
	Description types.String             `tfsdk:"description"`
	IsReadOnly  types.Bool               `tfsdk:"read_only"`
	Scopes      []globalaccountRoleScope `tfsdk:"scopes"`
}

type globalaccountRolesDataSourceConfig struct {
	/* OUTPUT */
	Id     types.String              `tfsdk:"id"`
	Values []globalaccountRolesValue `tfsdk:"values"`
}

type globalaccountRolesDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountRolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_roles", req.ProviderTypeName)
}

func (ds *globalaccountRolesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountRolesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists all roles.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/role-collections-and-roles-in-global-accounts-directories-and-subaccounts>`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{ // required hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `btp_globalaccount` datasource instead",
				MarkdownDescription: "The ID of the global account.",
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
							MarkdownDescription: "The ID of the xsuaa application.",
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

func (ds *globalaccountRolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountRolesDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Role.ListByGlobalAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = types.StringValue(ds.cli.GetGlobalAccountSubdomain())
	data.Values = []globalaccountRolesValue{}

	for _, role := range cliRes {
		roleVal := globalaccountRolesValue{
			Name:              types.StringValue(role.RoleTemplateName),
			RoleTemplateAppId: types.StringValue(role.RoleTemplateAppId),
			RoleTemplateName:  types.StringValue(role.RoleTemplateName),
			Description:       types.StringValue(role.Description),
			IsReadOnly:        types.BoolValue(role.IsReadOnly),
			Scopes:            []globalaccountRoleScope{},
		}

		for _, scope := range role.Scopes {
			scopeVal := globalaccountRoleScope{
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
