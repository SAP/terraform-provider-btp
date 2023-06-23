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
)

func newGlobalaccountRoleCollectionDataSource() datasource.DataSource {
	return &globalaccountRoleCollectionDataSource{}
}

type globalaccountRoleCollectionRoleType struct {
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	RoleTemplateAppId types.String `tfsdk:"role_template_app_id"`
	Description       types.String `tfsdk:"description"`
	Name              types.String `tfsdk:"name"`
}

type globalaccountRoleCollectionDataSourceConfig struct {
	Id types.String `tfsdk:"id"`

	/* OUTPUT */
	Name        types.String                          `tfsdk:"name"`
	IsReadOnly  types.Bool                            `tfsdk:"read_only"`
	Description types.String                          `tfsdk:"description"`
	Roles       []globalaccountRoleCollectionRoleType `tfsdk:"roles"`
}

type globalaccountRoleCollectionDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountRoleCollectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_role_collection", req.ProviderTypeName)
}

func (ds *globalaccountRoleCollectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountRoleCollectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific global account role collection.`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				DeprecationMessage:  "Use the `btp_globalaccount` datasource instead",
				MarkdownDescription: "The ID of the global account.",
				Computed:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Whether the role collection is read-only.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the role collection.",
				Computed:            true,
			},
			"roles": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role_template_name": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced role template.",
							Computed:            true,
						},
						"role_template_app_id": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced template app id",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the referenced role",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced role.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *globalaccountRoleCollectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountRoleCollectionDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.RoleCollection.GetByGlobalAccount(ctx, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = types.StringValue(ds.cli.GetGlobalAccountSubdomain())
	data.Name = types.StringValue(cliRes.Name)
	data.Description = types.StringValue(cliRes.Description)
	data.IsReadOnly = types.BoolValue(cliRes.IsReadOnly)

	data.Roles = []globalaccountRoleCollectionRoleType{}
	for _, ref := range cliRes.RoleReferences {
		data.Roles = append(data.Roles, globalaccountRoleCollectionRoleType{
			RoleTemplateName:  types.StringValue(ref.RoleTemplateName),
			RoleTemplateAppId: types.StringValue(ref.RoleTemplateAppId),
			Description:       types.StringValue(ref.Description),
			Name:              types.StringValue(ref.Name),
		})
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
