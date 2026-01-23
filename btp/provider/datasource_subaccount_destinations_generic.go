package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/jsonvalidator"
)

func newSubaccountDestinationsGenericDataSource() datasource.DataSource {
	return &subaccountDestinationsGenericDataSource{}
}

type subaccountDestinationsGenericDataSourceConfig struct {
	/* INPUT */
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
	/* OUTPUT */
	Values []subaccountDestinationGenericType `tfsdk:"values"`
}

type subaccountDestinationsGenericDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountDestinationsGenericDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destinations_generic", req.ProviderTypeName)
}

func (ds *subaccountDestinationsGenericDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountDestinationsGenericDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets list of all subaccount destinations details/names.
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
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"subaccount_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the subaccount.",
							Computed:            true,
						},
						"creation_time": schema.StringAttribute{
							MarkdownDescription: "The date and time when the resource was created.",
							Computed:            true,
						},
						"etag": schema.StringAttribute{
							MarkdownDescription: "The etag for the destination resource",
							Computed:            true,
						},
						"modification_time": schema.StringAttribute{
							MarkdownDescription: "The date and time when the resource was modified.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The descriptive name of the destination for subaccount",
							Computed:            true,
						},
						"service_instance_id": schema.StringAttribute{
							MarkdownDescription: "The service instance that becomes part of the path used to access the destination of the subaccount.",
							Computed:            true,
						},
						"destination_configuration": schema.StringAttribute{
							MarkdownDescription: "The configuration parameters for the destination.",
							Computed:            true,
							CustomType:          jsontypes.NormalizedType{},
							Validators: []validator.String{
								jsonvalidator.ValidJSON(),
							},
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountDestinationsGenericDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountDestinationsGenericDataSourceConfig
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Connectivity.Destination.ListBySubaccount(ctx, data.SubaccountID.ValueString(), data.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading destinations", fmt.Sprintf("%s", err))
		return
	}

	for _, destination := range cliRes {
		result, diags := destinationGenericDatasourceValueFrom(destination, data.SubaccountID, data.ServiceInstanceID)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, result)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
