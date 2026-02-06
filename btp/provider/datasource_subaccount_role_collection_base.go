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

func newSubaccountRoleCollectionBaseDataSource() datasource.DataSource {
	return &subaccountRoleCollectionBaseDataSource{}
}

type subaccountRoleCollectionBaseDataSourceConfig struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	IsReadOnly   types.Bool   `tfsdk:"read_only"`
}

type subaccountRoleCollectionBaseDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountRoleCollectionBaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role_collection_base", req.ProviderTypeName)
}

func (ds *subaccountRoleCollectionBaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountRoleCollectionBaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets metadata for a specific subaccount role collection.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the role collection.",
				Computed:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Whether the role collection is read-only.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountRoleCollectionBaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountRoleCollectionBaseDataSourceConfig
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleCollection, _, err := ds.cli.Security.RoleCollection.GetBySubaccount(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Role Collection Base", fmt.Sprintf("%s", err))
		return
	}

	data.Name = types.StringValue(roleCollection.Name)
	data.Description = types.StringValue(roleCollection.Description)
	data.IsReadOnly = types.BoolValue(roleCollection.IsReadOnly)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
