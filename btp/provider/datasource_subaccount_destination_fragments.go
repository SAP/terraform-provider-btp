package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/connectivity"
)

func newSubaccountDestinationFragmentsDataSource() datasource.DataSource {
	return &subaccountDestinationFragmentsDataSource{}
}

type subaccountDestinationFragmentsType struct {
	/* INPUT */
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
	/* OUTPUT */
	ID     types.String                               `tfsdk:"id"`
	Values []subaccountDestinationFragmentContentType `tfsdk:"values"`
}

type subaccountDestinationFragmentContentType struct {
	FragmentContent map[string]string `tfsdk:"fragment_content"`
}

type subaccountDestinationFragmentsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountDestinationFragmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_destination_fragments", req.ProviderTypeName)
}

func (ds *subaccountDestinationFragmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountDestinationFragmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a list of subaccount destination fragments.

__Tip:__
You must have the appropriate connectivity and destination permissions, such as:
- Subaccount Administrator  
- Destination Administrator  
- Destination Viewer  
- Connectivity and Destination Administrator

__Scope:__
- **Subaccount-level fragments**: Specify only the 'subaccount_id' and 'name' attribute.
- **Service instance-level fragments**: Specify the 'subaccount_id', 'service_instance_id' and 'name' attributes.

__Notes:__
- 'service_instance_id' is optional. When omitted, the fragments are searched at the subaccount level.`,
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
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service instance associated with the destination fragment.",
				Optional:            true,
				Computed:            true,
			},
			"values": schema.ListNestedAttribute{
				MarkdownDescription: "The list of destination fragments.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fragment_content": schema.MapAttribute{
							MarkdownDescription: "The content of the destination fragment.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (ds *subaccountDestinationFragmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountDestinationFragmentsType
	var destinationFragmentsDetails []connectivity.DestinationFragment
	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasServiceInstance := !data.ServiceInstanceID.IsNull() && !data.ServiceInstanceID.IsUnknown() && data.ServiceInstanceID.ValueString() != ""
	var err error

	if hasServiceInstance {
		destinationFragmentsDetails, _, err = ds.cli.Connectivity.DestinationFragment.ListByServiceInstance(ctx, data.SubaccountID.ValueString(), data.ServiceInstanceID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading Destination Fragments at Service Instance Level", fmt.Sprintf("%s", err))
			return
		}
	} else {
		destinationFragmentsDetails, _, err = ds.cli.Connectivity.DestinationFragment.ListBySubaccount(ctx, data.SubaccountID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Reading Destination Fragments at Subaccount Level", fmt.Sprintf("%s", err))
			return
		}
	}

	data.ID = types.StringValue(data.SubaccountID.ValueString())
	data.Values = []subaccountDestinationFragmentContentType{}

	for _, fragment := range destinationFragmentsDetails {
		content := subaccountDestinationFragmentContentType{
			FragmentContent: fragment.Content,
		}
		data.Values = append(data.Values, content)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
