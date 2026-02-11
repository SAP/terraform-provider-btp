package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountSecurityIdentityProviderDataSource() datasource.DataSource {
	return &subaccountSecurityIdentityProviderDataSource{}
}

type subaccountSecurityIdentityProviderDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountSecurityIdentityProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_identity_provider", req.ProviderTypeName)
}

func (ds *subaccountSecurityIdentityProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountSecurityIdentityProviderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific available identity provider for a subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
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

func (ds *subaccountSecurityIdentityProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountIdentityProviderDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Idp.GetBySubaccount(
		ctx,
		data.SubaccountId.ValueString(),
		data.Host.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("API Error Reading IdP Details", err.Error())
		return
	}

	state, diags := subaccountIdentityProviderDataSourceValueFrom(cliRes)
	state.SubaccountId = data.SubaccountId
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
