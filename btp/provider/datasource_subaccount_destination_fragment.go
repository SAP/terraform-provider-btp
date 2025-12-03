package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/connectivity"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountDestinationFragmentDataSource() datasource.DataSource {
	return &subaccountDestinationFragmentDataSource{}
}

type subaccountDestinationFragmentType struct {
	/* INPUT */
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	Name              types.String `tfsdk:"name"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
	/* OUTPUT */
	ID                  types.String `tfsdk:"id"`
	DestinationFragment types.Map    `tfsdk:"fragment_content"`
}

type subaccountDestinationFragmentDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountDestinationFragmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destination_fragment", req.ProviderTypeName)
}

func (ds *subaccountDestinationFragmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountDestinationFragmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific subaccount destination fragment.

__Tip:__
You must be assigned admin role of the subaccount and destination service.

__Scope:__
- **Subaccount-level fragment**: Specify only the 'subaccount_id' and 'name' attribute.
- **Service instance-level fragment**: Specify the 'subaccount_id', 'service_instance_id' and 'name' attributes.

__Notes:__
- 'service_instance_id' is optional. When omitted, the fragment is searched at the subaccount level.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` attribute instead",
				MarkdownDescription: "The ID of the subaccount.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the destination fragment.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service instance associated with the destination fragment.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(uuidvalidator.UuidRegexp, "value must be a valid UUID"),
				},
			},
			"fragment_content": schema.MapAttribute{
				MarkdownDescription: "The content of the destination fragment.",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountDestinationFragmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountDestinationFragmentType
	var destinationFragmentDetails connectivity.DestinationFragment
	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasServiceInstance := !data.ServiceInstanceID.IsNull() && !data.ServiceInstanceID.IsUnknown() && data.ServiceInstanceID.ValueString() != ""
	var err error

	if hasServiceInstance {
		destinationFragmentDetails, _, err = ds.cli.Connectivity.DestinationFragment.GetByServiceInstance(ctx, data.SubaccountID.ValueString(), data.Name.ValueString(), data.ServiceInstanceID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading Destination Fragment at Service Instance Level", fmt.Sprintf("%s", err))
			return
		}
	} else {
		destinationFragmentDetails, _, err = ds.cli.Connectivity.DestinationFragment.GetBySubaccount(ctx, data.SubaccountID.ValueString(), data.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading Destination Fragment at Subaccount Level", fmt.Sprintf("%s", err))
			return
		}
	}

	data.ID = types.StringValue(data.SubaccountID.ValueString())

	data.DestinationFragment, diags = types.MapValueFrom(ctx, types.StringType, destinationFragmentDetails.Content)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
