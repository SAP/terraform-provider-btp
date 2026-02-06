package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountRoleCollectionBasesDataSource() datasource.DataSource {
	return &subaccountRoleCollectionBasesDataSource{}
}

type subaccountRoleCollectionBasesValueModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IsReadOnly  types.Bool   `tfsdk:"read_only"`
}

var roleBasesAttrTypes = map[string]attr.Type{
	"name":        types.StringType,
	"description": types.StringType,
	"read_only":   types.BoolType,
}

type subaccountRoleCollectionBasesDataSourceModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Values       types.Set    `tfsdk:"values"`
}

type subaccountRoleCollectionBasesDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountRoleCollectionBasesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role_collection_bases", req.ProviderTypeName)
}

func (ds *subaccountRoleCollectionBasesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountRoleCollectionBasesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists metadata of all role collections in a subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				Required: true,
			},
			"values": schema.SetNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":        schema.StringAttribute{Computed: true},
						"description": schema.StringAttribute{Computed: true},
						"read_only":   schema.BoolAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (ds *subaccountRoleCollectionBasesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountRoleCollectionBasesDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	roleCollections, _, err := ds.cli.Security.RoleCollection.ListBySubaccount(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Role Collection Bases", fmt.Sprintf("%s", err))
		return
	}

	var bases []subaccountRoleCollectionBasesValueModel
	for _, rc := range roleCollections {
		bases = append(bases, subaccountRoleCollectionBasesValueModel{
			Name:        types.StringValue(rc.Name),
			Description: types.StringValue(rc.Description),
			IsReadOnly:  types.BoolValue(rc.IsReadOnly),
		})
	}

	baseSet, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: roleBasesAttrTypes}, bases)
	resp.Diagnostics.Append(diags...)

	data.Values = baseSet
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
