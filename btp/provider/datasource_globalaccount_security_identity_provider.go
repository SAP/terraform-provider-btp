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

type globalaccountIdentityProviderDataSourceModel struct {
	TenantId     types.String `tfsdk:"tenant_id"`
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
	Region       types.String `tfsdk:"region"`
	Status       types.String `tfsdk:"status"`
}

func newGlobalaccountSecurityIdentityProviderDataSource() datasource.DataSource {
	return &globalaccountSecurityIdentityProviderDataSource{}
}

type globalaccountSecurityIdentityProviderDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountSecurityIdentityProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_identity_provider", req.ProviderTypeName)
}

func (ds *globalaccountSecurityIdentityProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountSecurityIdentityProviderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific available identity provider for a globalaccount.`,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The host of the identity provider.",
				Required:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the tenant.",
				Computed:            true,
			},
			"tenant_type": schema.StringAttribute{
				MarkdownDescription: "The type of the tenant.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the identity provider.",
				Computed:            true,
			},
			"common_host": schema.StringAttribute{
				MarkdownDescription: "The common host of the identity provider.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the identity provider.",
				Computed:            true,
			},
			"custom_host": schema.StringAttribute{
				MarkdownDescription: "The custom host of the identity provider.",
				Computed:            true,
			},
			"customer_name": schema.StringAttribute{
				MarkdownDescription: "The name of the customer.",
				Computed:            true,
			},
			"cost_center_id": schema.Int64Attribute{
				MarkdownDescription: "The cost center ID associated with the entity.",
				Computed:            true,
			},
			"data_center_id": schema.StringAttribute{
				MarkdownDescription: "The data center ID.",
				Computed:            true,
			},
			"customer_id": schema.StringAttribute{
				MarkdownDescription: "The customer ID.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region where the identity provider is located.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The current status of the identity provider.",
				Computed:            true,
			},
		},
	}
}

func (ds *globalaccountSecurityIdentityProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountIdentityProviderDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Idp.GetByGlobalAccount(
		ctx,
		data.Host.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("API Error Reading IdP Details", err.Error())
		return
	}

	state, diags := globalaccountIdentityProviderDataSourceValueFrom(cliRes)
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func globalaccountIdentityProviderDataSourceValueFrom(val xsuaa_authz.Idp) (data globalaccountIdentityProviderDataSourceModel, diags diag.Diagnostics) {
	data.TenantId = types.StringValue(val.TenantId)
	data.TenantType = types.StringValue(val.TenantType)
	data.DisplayName = types.StringPointerValue(val.DisplayName)
	data.CommonHost = types.StringValue(val.CommonHost)
	data.Description = types.StringValue(val.Description)
	data.CustomHost = types.StringPointerValue(val.CustomHost)
	data.CustomerName = types.StringPointerValue(val.CustomerName)
	data.CostCenterId = types.Int64Value(int64(val.CostCenterId))
	data.DataCenterId = types.StringValue(val.DataCenterId)
	data.Host = types.StringValue(val.Host)
	data.CustomerId = types.StringPointerValue(val.CustomerId)
	data.Region = types.StringValue(val.Region)
	data.Status = types.StringValue(val.Status)

	return
}
