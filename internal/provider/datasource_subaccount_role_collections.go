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

func newSubaccountRoleCollectionsDataSource() datasource.DataSource {
	return &subaccountRoleCollectionsDataSource{}
}

type subaccountRoleCollectionsRoleType struct {
	/* OUTPUT */
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	RoleTemplateAppId types.String `tfsdk:"role_template_app_id"`
	Description       types.String `tfsdk:"description"`
	Name              types.String `tfsdk:"name"`
}

type subaccountRoleCollectionsValueConfig struct {
	/* OUTPUT */
	Name        types.String                        `tfsdk:"name"`
	IsReadOnly  types.Bool                          `tfsdk:"read_only"`
	Description types.String                        `tfsdk:"description"`
	Roles       []subaccountRoleCollectionsRoleType `tfsdk:"roles"`
}

type subaccountRoleCollectionsDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	/* OUTPUT */
	Values []subaccountRoleCollectionsValueConfig `tfsdk:"values"`
}

type subaccountRoleCollectionsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountRoleCollectionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role_collections", req.ProviderTypeName)
}

func (ds *subaccountRoleCollectionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountRoleCollectionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `List all role collections.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` attribute instead",
				MarkdownDescription: "The ID of the subaccount.",
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
							MarkdownDescription: "Whether the role collection is readonly.",
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

func (ds *subaccountRoleCollectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountRoleCollectionsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.RoleCollection.ListBySubaccount(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collections (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.SubaccountId
	data.Values = []subaccountRoleCollectionsValueConfig{}

	for _, rolecollection := range cliRes {
		val := subaccountRoleCollectionsValueConfig{
			Name:        types.StringValue(rolecollection.Name),
			Description: types.StringValue(rolecollection.Description),
			IsReadOnly:  types.BoolValue(rolecollection.IsReadOnly),
			Roles:       []subaccountRoleCollectionsRoleType{},
		}

		for _, ref := range rolecollection.RoleReferences {
			val.Roles = append(val.Roles, subaccountRoleCollectionsRoleType{
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
