package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountSecurityIdentityProvidersDataSource() datasource.DataSource {
	return &globalaccountSecurityIdentityProvidersDataSource{}
}

type globalaccountSecurityIdentityProvidersDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountSecurityIdentityProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_identity_providers", req.ProviderTypeName)
}

func (ds *globalaccountSecurityIdentityProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountSecurityIdentityProvidersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists the available identity providers for a global account.`,
		Attributes: map[string]schema.Attribute{
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

func (ds *globalaccountSecurityIdentityProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountIdentityProvidersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Idp.ListByGlobalAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Available IdPs", err.Error())
		return
	}

	state, diags := globalaccountIdentityProvidersDataSourceValueFrom(cliRes)
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
