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

func newSubaccountDestinationsDataSource() datasource.DataSource {
	return &subaccountDestinationsDataSource{}
}

type subaccountDestinationsDataSourceConfig struct {
	/* INPUT */
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
	NamesOnly         types.Bool   `tfsdk:"names_only"`
	/* OUTPUT */
	Values           []subaccountDestinationType `tfsdk:"values"`
	DestinationNames []subaccountDestinationName `tfsdk:"destination_names"`
}

type subaccountDestinationsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountDestinationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destinations", req.ProviderTypeName)
}

func (ds *subaccountDestinationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountDestinationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"names_only": schema.BoolAttribute{
				MarkdownDescription: "The Bool value for getting names only. Default value is false.",
				Optional:            true,
			},
			"destination_names": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The descriptive name of the destination for subaccount",
							Required:            true,
						},
					},
				},
				Computed: true,
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
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of request from destination.",
							Computed:            true,
						},
						"proxy_type": schema.StringAttribute{
							MarkdownDescription: "The proxytype of the destination.",
							Computed:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "The url of the destination.",
							Computed:            true,
						},
						"authentication": schema.StringAttribute{
							MarkdownDescription: "The authentication of the destination.",
							Computed:            true,
						},
						"service_instance_id": schema.StringAttribute{
							MarkdownDescription: "The service instance that becomes part of the path used to access the destination of the subaccount.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the destination.",
							Computed:            true,
						},
						"additional_configuration": schema.StringAttribute{
							MarkdownDescription: "The additional configuration parameters for the destination.",
							Optional:            true,
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

func (ds *subaccountDestinationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountDestinationsDataSourceConfig
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.NamesOnly.ValueBool() {
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
	} else {
		cliRes, _, err := ds.cli.Connectivity.Destination.ListBySubaccount(ctx, data.SubaccountID.ValueString(), data.ServiceInstanceID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading destinations", fmt.Sprintf("%s", err))
			return
		}
		for _, destination := range cliRes {
			result, diags := destinationDatasourceValueFrom(destination, data.SubaccountID, data.ServiceInstanceID)
			resp.Diagnostics.Append(diags...)

			data.Values = append(data.Values, result)
		}
		diags = resp.State.Set(ctx, &data)
		resp.Diagnostics.Append(diags...)
	}
}
