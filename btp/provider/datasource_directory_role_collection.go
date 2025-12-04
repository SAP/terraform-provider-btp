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
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newDirectoryRoleCollectionDataSource() datasource.DataSource {
	return &directoryRoleCollectionDataSource{}
}

type directoryRoleCollectionRoleType struct {
	/* OUTPUT */
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	RoleTemplateAppId types.String `tfsdk:"role_template_app_id"`
	Description       types.String `tfsdk:"description"`
	Name              types.String `tfsdk:"name"`
}

type directoryRoleCollectionAttributeMappingsType struct {
	/* OUTPUT */
	IdentityProvider types.String `tfsdk:"identity_provider"`
	Attribute        types.String `tfsdk:"attribute"`
	Operator         types.String `tfsdk:"operator"`
	Value            types.String `tfsdk:"value"`
}

type directoryRoleCollectionDataSourceConfig struct {
	/* INPUT */
	DirectoryId types.String `tfsdk:"directory_id"`
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	/* OUTPUT */
	IsReadOnly            types.Bool                                     `tfsdk:"read_only"`
	Description           types.String                                   `tfsdk:"description"`
	Roles                 []directoryRoleCollectionRoleType              `tfsdk:"roles"`
	ShowAttributeMappings types.Bool                                     `tfsdk:"show_attribute_mappings"`
	AttributeMappings     []directoryRoleCollectionAttributeMappingsType `tfsdk:"attribute_mappings"`
}

type directoryRoleCollectionDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *directoryRoleCollectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_role_collection", req.ProviderTypeName)
}

func (ds *directoryRoleCollectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *directoryRoleCollectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific directory role collection.

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
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
			"show_attribute_mappings": schema.BoolAttribute{
				MarkdownDescription: "If set to true, the data source will also return which user attributes and user groups provided by an identity provider effectively grant this role collection.",
				Optional:            true,
			},
			"attribute_mappings": schema.SetNestedAttribute{
				MarkdownDescription: "List of user attributes and user groups from identity providers that effectively grant this role collection.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identity_provider": schema.StringAttribute{
							MarkdownDescription: "The display name of the identity provider from which the attribute or group mapping originates.",
							Computed:            true,
						},
						"attribute": schema.StringAttribute{
							MarkdownDescription: "The user attribute or group name used in the mapping.",
							Computed:            true,
						},
						"operator": schema.StringAttribute{
							MarkdownDescription: "The operator applied in the attribute mapping. Only `equals` is currently supported.",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("EQUALS, equals"),
							},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "The value of the user attribute or group that grants the role collection.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (ds *directoryRoleCollectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryRoleCollectionDataSourceConfig
	var roleCollectionDetails, roleCollectionAttributeMappings xsuaa_authz.RoleCollection
	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleCollectionDetails, _, err := ds.cli.Security.RoleCollection.GetByDirectory(ctx, data.DirectoryId.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Directory)", fmt.Sprintf("%s", err))
		return
	}

	if data.ShowAttributeMappings.ValueBool() {
		roleCollectionAttributeMappings, _, err = ds.cli.Security.RoleCollection.GetByDirectoryWithAttributeMappings(ctx, data.DirectoryId.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Directory) Attribute Mappings", fmt.Sprintf("%s", err))
			return
		}
	}

	data.Id = data.DirectoryId
	data.Name = types.StringValue(roleCollectionDetails.Name)
	data.Description = types.StringValue(roleCollectionDetails.Description)
	data.IsReadOnly = types.BoolValue(roleCollectionDetails.IsReadOnly)

	data.Roles = []directoryRoleCollectionRoleType{}
	for _, ref := range roleCollectionDetails.RoleReferences {
		data.Roles = append(data.Roles, directoryRoleCollectionRoleType{
			RoleTemplateName:  types.StringValue(ref.RoleTemplateName),
			RoleTemplateAppId: types.StringValue(ref.RoleTemplateAppId),
			Description:       types.StringValue(ref.Description),
			Name:              types.StringValue(ref.Name),
		})
	}

	if data.ShowAttributeMappings.ValueBool() {
		// Attribute mappings
		for _, am := range roleCollectionAttributeMappings.SamlAttributeAssignment {
			data.AttributeMappings = append(data.AttributeMappings, directoryRoleCollectionAttributeMappingsType{
				IdentityProvider: types.StringValue(am.IdentityProvider),
				Attribute:        types.StringValue(am.AttributeName),
				Operator:         types.StringValue(am.ComparisonOperator),
				Value:            types.StringValue(am.SamlAttributeValue),
			})
		}
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
