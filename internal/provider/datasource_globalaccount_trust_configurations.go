package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountTrustConfigurationsDataSource() datasource.DataSource {
	return &globalaccountTrustConfigurationsDataSource{}
}

type globalaccountTrustConfigurationsDataSourceConfig struct {
	Values []globalaccountTrustConfigurationType `tfsdk:"values"`
}

type globalaccountTrustConfigurationsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountTrustConfigurationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_trust_configurations", req.ProviderTypeName)
}

func (ds *globalaccountTrustConfigurationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountTrustConfigurationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `List all trust configurations that are configured for your global account.

__Tip__
You must be viewer or administrator of the global account.

__Further documentation__
https://help.sap.com/docs/BTP/65de2977205c403bbc107264b8eccf4b/cb1bc8f1bd5c482e891063960d7acd78.html`,
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"origin": schema.StringAttribute{
							MarkdownDescription: "The origin of the identity provider.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the trust configuration.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the trust configuration.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the trust configuration.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The trust type.",
							Computed:            true,
						},
						"identity_provider": schema.StringAttribute{
							MarkdownDescription: "The name of the identity provider.",
							Computed:            true,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "The protocol used to establish trust with the identity provider.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Whether the identity provider is currently active or not.",
							Computed:            true,
						},
						"read_only": schema.BoolAttribute{
							MarkdownDescription: "Whether the trust configuration can be modified.",
							Computed:            true,
						},
					},
				},
				MarkdownDescription: "Trust configurations associated with the global account.",
				Computed:            true,
			},
		},
	}
}

func (ds *globalaccountTrustConfigurationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountTrustConfigurationsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Trust.ListByGlobalAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Trust Configurations (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	data.Values = []globalaccountTrustConfigurationType{}

	for _, trustConfig := range cliRes {
		trustConfigValue, diags := globalaccountTrustConfigurationFromValue(ctx, trustConfig)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, trustConfigValue)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
