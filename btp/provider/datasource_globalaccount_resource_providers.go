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
	return &globalaccountResourceProvidersDataSource{}
}

type globalaccountResourceProvidersValue struct {
	Provider      types.String `tfsdk:"provider_type"`
	TechnicalName types.String `tfsdk:"technical_name"`
	DisplayName   types.String `tfsdk:"display_name"`
	Description   types.String `tfsdk:"description"`
}

type globalaccountResourceProvidersDataSourceConfig struct {
	/* INPUT */
	/* OUTPUT */
	Id     types.String                          `tfsdk:"id"`
	Values []globalaccountResourceProvidersValue `tfsdk:"values"`
}

type globalaccountResourceProvidersDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountResourceProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_resource_providers", req.ProviderTypeName)
}

func (ds *globalaccountResourceProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountResourceProvidersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists all the resource provider instances in a global account.

__Tip:__
You must be assigned to the global account admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/managing-resource-providers>`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `btp_globalaccount` datasource instead",
				MarkdownDescription: "The ID of the global account",
				Computed:            true,
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"provider_type": schema.StringAttribute{
							MarkdownDescription: "The cloud vendor from which to consume services through your subscribed account. Possible values are: \n" +
								getFormattedValueAsTableRow("value", "description") +
								getFormattedValueAsTableRow("---", "---") +
								getFormattedValueAsTableRow("`AWS`", "Amazon Web Services") +
								getFormattedValueAsTableRow("`AZURE`", "Microsoft Azure"),
							Computed: true,
						},
						"technical_name": schema.StringAttribute{
							MarkdownDescription: "The unique technical name of the resource provider.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "The descriptive name of the resource provider.",
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

func (ds *globalaccountResourceProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountResourceProvidersDataSourceConfig

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

	data.Id = types.StringValue(ds.cli.GetGlobalAccountSubdomain())
	data.Values = []globalaccountResourceProvidersValue{}

	for _, provider := range cliRes {
		resourceProvider := globalaccountResourceProvidersValue{
			Provider:      types.StringValue(provider.ResourceProvider),
			TechnicalName: types.StringValue(provider.TechnicalName),
			DisplayName:   types.StringValue(provider.DisplayName),
			Description:   types.StringValue(provider.Description),
		}

		data.Values = append(data.Values, resourceProvider)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
