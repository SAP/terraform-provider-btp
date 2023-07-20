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

func newDirectoryRoleCollectionsDataSource() datasource.DataSource {
	return &directoryRoleCollectionsDataSource{}
}

type directoryRoleCollectionsRoleType struct {
	/* OUTPUT */
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	RoleTemplateAppId types.String `tfsdk:"role_template_app_id"`
	Description       types.String `tfsdk:"description"`
	Name              types.String `tfsdk:"name"`
}

type directoryRoleCollectionsValueConfig struct {
	/* OUTPUT */
	Name        types.String                       `tfsdk:"name"`
	IsReadOnly  types.Bool                         `tfsdk:"read_only"`
	Description types.String                       `tfsdk:"description"`
	Roles       []directoryRoleCollectionsRoleType `tfsdk:"roles"`
}

type directoryRoleCollectionsDataSourceConfig struct {
	DirectoryId types.String `tfsdk:"directory_id"`
	Id          types.String `tfsdk:"id"`
	/* OUTPUT */
	Values []directoryRoleCollectionsValueConfig `tfsdk:"values"`
}

type directoryRoleCollectionsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *directoryRoleCollectionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_role_collections", req.ProviderTypeName)
}

func (ds *directoryRoleCollectionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *directoryRoleCollectionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets all role collections.`,
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
							MarkdownDescription: "The name of the role collection.",
							Computed:            true,
						},
						"read_only": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the role collection is read-only.",
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
				},
				Computed: true,
			},
		},
	}
}

func (ds *directoryRoleCollectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryRoleCollectionsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.RoleCollection.ListByDirectory(ctx, data.DirectoryId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collections (Directory)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.DirectoryId
	data.Values = []directoryRoleCollectionsValueConfig{}

	for _, rolecollection := range cliRes {
		val := directoryRoleCollectionsValueConfig{
			Name:        types.StringValue(rolecollection.Name),
			Description: types.StringValue(rolecollection.Description),
			IsReadOnly:  types.BoolValue(rolecollection.IsReadOnly),
			Roles:       []directoryRoleCollectionsRoleType{},
		}

		for _, ref := range rolecollection.RoleReferences {
			val.Roles = append(val.Roles, directoryRoleCollectionsRoleType{
				RoleTemplateName:  types.StringValue(ref.RoleTemplateName),
				RoleTemplateAppId: types.StringValue(ref.RoleTemplateAppId),
				Description:       types.StringValue(ref.Description),
				Name:              types.StringValue(ref.Name),
			})
		}
		// TODO map additional attributes
		data.Values = append(data.Values, val)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
