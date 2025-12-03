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

func newSubaccountDestinationDataSource() datasource.DataSource {
	return &subaccountDestinationDataSource{}
}

type subaccountDestinationDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountDestinationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destination", req.ProviderTypeName)
}

func (ds *subaccountDestinationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountDestinationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific subaccount destination.
__Tip:__
You must be assigned to the admin or viewer role of the subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A name of the destination for subaccount",
				Required:            true,
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The service instance that becomes part of the path used to access the destination of the subaccount.",
				Optional:            true,
			},
			"creation_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"etag": schema.StringAttribute{
				MarkdownDescription: "etag for the destination resource",
				Computed:            true,
			},
			"modification_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
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
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the destination.",
				Computed:            true,
			},
			"additional_configuration": schema.StringAttribute{
				MarkdownDescription: "The additional configuration parameters for the destination.",
				Computed:            true,
				CustomType:          jsontypes.NormalizedType{},
				Validators: []validator.String{
					jsonvalidator.ValidJSON(),
				},
			},
		},
	}
}

func (ds *subaccountDestinationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountDestinationType
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Connectivity.Destination.Get(ctx, data.SubaccountID.ValueString(), data.Name.ValueString(), data.ServiceInstanceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading destination", fmt.Sprintf("%s", err))
		return
	}

	data, diags = destinationValueFrom(cliRes, data.SubaccountID, data.ServiceInstanceID)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
