package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountResourceProvidersDataSource() datasource.DataSource {
	return &globalaccountGlobalaccountResourceProvidersDataSource{}
}

type globalaccountGlobalaccountResourceProvidersValue struct {
	ResourceProvider types.String `tfsdk:"resource_provider"`
	Id               types.String `tfsdk:"id"`
	DisplayName      types.String `tfsdk:"display_name"`
	Description      types.String `tfsdk:"description"`
}

type globalaccountGlobalaccountResourceProvidersDataSourceConfig struct {
	/* INPUT */
	/* OUTPUT */
	Values []globalaccountGlobalaccountResourceProvidersValue `tfsdk:"values"`
}

type globalaccountGlobalaccountResourceProvidersDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountGlobalaccountResourceProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_resource_providers", req.ProviderTypeName)
}

func (ds *globalaccountGlobalaccountResourceProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountGlobalaccountResourceProvidersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `List all the resource provider instances in a global account.

__Tips__
You must be assigned to the global account admin or viewer role.

__Further documentation__
https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/e2c250dc5abd468a81f4f619206157a2.html`,
		Attributes: map[string]schema.Attribute{
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_provider": schema.StringAttribute{
							MarkdownDescription: "Provider of the requested resource. For example: AWS, AZURE.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique technical name of the resource provider.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Descriptive name of the resource provider.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the resource provider.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *globalaccountGlobalaccountResourceProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountGlobalaccountResourceProvidersDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.ResourceProvider.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Resource Providers (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	data.Values = []globalaccountGlobalaccountResourceProvidersValue{}

	for _, provider := range cliRes {
		data.Values = append(data.Values, globalaccountGlobalaccountResourceProvidersValue{
			ResourceProvider: types.StringValue(provider.ResourceProvider),
			Id:               types.StringValue(provider.ResourceTechnicalName),
			DisplayName:      types.StringValue(provider.DisplayName),
			Description:      types.StringValue(provider.Description),
		})
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
