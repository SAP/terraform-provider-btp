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

type subaccountRoleCollectionDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	/* OUTPUT */
	Name        types.String                       `tfsdk:"name"`
	IsReadOnly  types.Bool                         `tfsdk:"read_only"`
	Description types.String                       `tfsdk:"description"`
	Roles       []subaccountRoleCollectionRoleType `tfsdk:"roles"`
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
		MarkdownDescription: `Gets details about a specific subaccount role collection.`,
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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

func (ds *subaccountRoleCollectionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountRoleCollectionDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rolecollection, _, err := ds.cli.Security.RoleCollection.GetBySubaccount(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.SubaccountId
	data.Name = types.StringValue(rolecollection.Name)
	data.Description = types.StringValue(rolecollection.Description)
	data.IsReadOnly = types.BoolValue(rolecollection.IsReadOnly)

	data.Roles = []subaccountRoleCollectionRoleType{}
	for _, ref := range rolecollection.RoleReferences {
		data.Roles = append(data.Roles, subaccountRoleCollectionRoleType{
			RoleTemplateName:  types.StringValue(ref.RoleTemplateName),
			RoleTemplateAppId: types.StringValue(ref.RoleTemplateAppId),
			Description:       types.StringValue(ref.Description),
			Name:              types.StringValue(ref.Name),
		})
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
