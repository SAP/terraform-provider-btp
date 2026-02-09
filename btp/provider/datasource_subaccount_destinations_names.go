package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountDestinationsNamesDataSource() datasource.DataSource {
	return &subaccountDestinationsNamesDataSource{}
}

type subaccountDestinationsNamesDataSourceConfig struct {
	/* INPUT */
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
	/* OUTPUT */
	DestinationNames []subaccountDestinationName `tfsdk:"destination_names"`
}

type subaccountDestinationsNamesDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountDestinationsNamesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destinations_names", req.ProviderTypeName)
}

func (ds *subaccountDestinationsNamesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountDestinationsNamesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets a list of all subaccount destination names.
__Tip:__
You must have the appropriate connectivity and destination permissions, such as:

Subaccount Administrator
Destination Administrator
Destination Viewer
Connectivity and Destination Administrator`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The service instance that becomes part of the path used to access the destination of the subaccount.",
				Optional:            true,
			},
			"destination_names": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The descriptive name of the destination for subaccount",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountDestinationsNamesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountDestinationsNamesDataSourceConfig
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Connectivity.Destination.ListNamesBySubaccount(ctx, data.SubaccountID.ValueString(), data.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading destinations", fmt.Sprintf("%s", err))
		return
	}

	data.DestinationNames = []subaccountDestinationName{}
	for _, destination := range cliRes {
		result := subaccountDestinationName{
			Name: types.StringValue(destination["Name"]),
		}
		data.DestinationNames = append(data.DestinationNames, result)
	}
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

}
