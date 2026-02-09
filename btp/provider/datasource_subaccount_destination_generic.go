package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/jsonvalidator"
)

func newSubaccountDestinationGenericDataSource() datasource.DataSource {
	return &subaccountDestinationGenericDataSource{}
}

type subaccountDestinationGenericDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountDestinationGenericDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destination_generic", req.ProviderTypeName)
}

func (ds *subaccountDestinationGenericDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountDestinationGenericDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific subaccount destination.
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the destination for subaccount",
				Required:            true,
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The service instance that becomes part of the path used to access the destination of the subaccount.",
				Optional:            true,
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
			"destination_configuration": schema.StringAttribute{
				MarkdownDescription: "The configuration parameters for the destination.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
				Validators: []validator.String{
					jsonvalidator.ValidJSON(),
				},
			},
		},
	}
}

func (ds *subaccountDestinationGenericDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountDestinationGenericType
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Connectivity.Destination.GetBySubaccount(ctx, data.SubaccountID.ValueString(), data.Name.ValueString(), data.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading destination", fmt.Sprintf("%s", err))
		return
	}

	data, diags = destinationGenericDatasourceValueFrom(cliRes, data.SubaccountID, data.ServiceInstanceID)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
