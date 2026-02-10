package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
)

type subaccountIdentityProvidersDataSourceModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Values       []idpModel   `tfsdk:"values"`
}

type idpModel struct {
	TenantType   types.String `tfsdk:"tenant_type"`
	DisplayName  types.String `tfsdk:"display_name"`
	CommonHost   types.String `tfsdk:"common_host"`
	Description  types.String `tfsdk:"description"`
	CustomHost   types.String `tfsdk:"custom_host"`
	CustomerName types.String `tfsdk:"customer_name"`
	CostCenterId types.Int64  `tfsdk:"cost_center_id"`
	DataCenterId types.String `tfsdk:"data_center_id"`
	Host         types.String `tfsdk:"host"`
	CustomerId   types.String `tfsdk:"customer_id"`
	TenantId     types.String `tfsdk:"tenant_id"`
	Region       types.String `tfsdk:"region"`
	Status       types.String `tfsdk:"status"`
}

func newSubaccountSecurityIdentityProvidersDataSource() datasource.DataSource {
	return &subaccountSecurityIdentityProvidersDataSource{}
}

type subaccountSecurityIdentityProvidersDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountSecurityIdentityProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_identity_providers", req.ProviderTypeName)
}

func (ds *subaccountSecurityIdentityProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountSecurityIdentityProvidersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists the available identity providers for a subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"values": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tenant_type":    schema.StringAttribute{Computed: true},
						"display_name":   schema.StringAttribute{Computed: true},
						"common_host":    schema.StringAttribute{Computed: true},
						"description":    schema.StringAttribute{Computed: true},
						"custom_host":    schema.StringAttribute{Computed: true},
						"customer_name":  schema.StringAttribute{Computed: true},
						"cost_center_id": schema.Int64Attribute{Computed: true},
						"data_center_id": schema.StringAttribute{Computed: true},
						"host":           schema.StringAttribute{Computed: true},
						"customer_id":    schema.StringAttribute{Computed: true},
						"tenant_id":      schema.StringAttribute{Computed: true},
						"region":         schema.StringAttribute{Computed: true},
						"status":         schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (ds *subaccountSecurityIdentityProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountIdentityProvidersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Idp.ListBySubaccount(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Available IdPs", err.Error())
		return
	}

	state, diags := subaccountIdentityProvidersDataSourceValueFrom(cliRes)
	state.SubaccountId = data.SubaccountId
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func subaccountIdentityProvidersDataSourceValueFrom(value []xsuaa_authz.Idp) (data subaccountIdentityProvidersDataSourceModel, diags diag.Diagnostics) {
	data.Values = []idpModel{}

	for _, val := range value {
		idp := idpModel{
			TenantType:   types.StringValue(val.TenantType),
			DisplayName:  types.StringPointerValue(val.DisplayName),
			CommonHost:   types.StringValue(val.CommonHost),
			Description:  types.StringValue(val.Description),
			CustomHost:   types.StringPointerValue(val.CustomHost),
			CustomerName: types.StringPointerValue(val.CustomerName),
			CostCenterId: types.Int64Value(int64(val.CostCenterId)),
			DataCenterId: types.StringValue(val.DataCenterId),
			Host:         types.StringValue(val.Host),
			CustomerId:   types.StringPointerValue(val.CustomerId),
			TenantId:     types.StringValue(val.TenantId),
			Region:       types.StringValue(val.Region),
			Status:       types.StringValue(val.Status),
		}
		data.Values = append(data.Values, idp)
	}

	return
}
