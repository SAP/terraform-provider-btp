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
)

func newSubaccountRoleCollectionDataSource() datasource.DataSource {
	return &subaccountRoleCollectionDataSource{}
}

type subaccountRoleCollectionRoleType struct {
	/* OUTPUT */
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	RoleTemplateAppId types.String `tfsdk:"role_template_app_id"`
	Description       types.String `tfsdk:"description"`
	Name              types.String `tfsdk:"name"`
}

type subaccountRoleCollectionAttributeMappingsType struct {
	/* OUTPUT */
	IdentityProvider types.String `tfsdk:"identity_provider"`
	Attribute        types.String `tfsdk:"attribute"`
	Operator         types.String `tfsdk:"operator"`
	Value            types.String `tfsdk:"value"`
}

type subaccountRoleCollectionDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	/* OUTPUT */
	Name                  types.String                                    `tfsdk:"name"`
	IsReadOnly            types.Bool                                      `tfsdk:"read_only"`
	Description           types.String                                    `tfsdk:"description"`
	Roles                 []subaccountRoleCollectionRoleType              `tfsdk:"roles"`
	ShowAttributeMappings types.Bool                                      `tfsdk:"show_attribute_mappings"`
	AttributeMappings     []subaccountRoleCollectionAttributeMappingsType `tfsdk:"attribute_mappings"`
}

type subaccountRoleCollectionDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountRoleCollectionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role_collection", req.ProviderTypeName)
}

func (ds *subaccountRoleCollectionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountRoleCollectionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific subaccount role collection.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` attribute instead",
				MarkdownDescription: "The ID of the subaccount.",
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

func (ds *subaccountRoleCollectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountRoleCollectionDataSourceConfig
	var roleCollectionDetails, roleCollectionAttributeMappings xsuaa_authz.RoleCollection
	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleCollectionDetails, _, err := ds.cli.Security.RoleCollection.GetBySubaccount(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	if data.ShowAttributeMappings.ValueBool() {
		roleCollectionAttributeMappings, _, err = ds.cli.Security.RoleCollection.GetBySubaccountWithAttributeMappings(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Subaccount) Attribute Mappings", fmt.Sprintf("%s", err))
			return
		}
	}

	data.Id = data.SubaccountId
	data.Name = types.StringValue(roleCollectionDetails.Name)
	data.Description = types.StringValue(roleCollectionDetails.Description)
	data.IsReadOnly = types.BoolValue(roleCollectionDetails.IsReadOnly)

	data.Roles = []subaccountRoleCollectionRoleType{}
	for _, ref := range roleCollectionDetails.RoleReferences {
		data.Roles = append(data.Roles, subaccountRoleCollectionRoleType{
			RoleTemplateName:  types.StringValue(ref.RoleTemplateName),
			RoleTemplateAppId: types.StringValue(ref.RoleTemplateAppId),
			Description:       types.StringValue(ref.Description),
			Name:              types.StringValue(ref.Name),
		})
	}

	if data.ShowAttributeMappings.ValueBool() {
		// Attribute mappings
		for _, am := range roleCollectionAttributeMappings.SamlAttributeAssignment {
			data.AttributeMappings = append(data.AttributeMappings, subaccountRoleCollectionAttributeMappingsType{
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
