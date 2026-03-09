package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ list.ListResourceWithConfigure = &subaccountDestinationFragmentListResource{}

type subaccountDestinationFragmentListResource struct {
	client *btpcli.ClientFacade
}

type subaccountDestinationFragmentListResourceFilter struct {
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
}

func NewSubaccountDestinationFragmentListResource() list.ListResource {
	return &subaccountDestinationFragmentListResource{}
}

func (r *subaccountDestinationFragmentListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_destination_fragment" // must match managed resource
}

func (r *subaccountDestinationFragmentListResource) Configure(_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*btpcli.ClientFacade)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *btpcli.ClientFacade, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *subaccountDestinationFragmentListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all destination fragments available for given subaccount.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service instance.",
				Optional:            true,
			},
		},
	}
}

func (r *subaccountDestinationFragmentListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {
	var (
		filter subaccountDestinationFragmentListResourceFilter
		err    error
		cliRes []connectivity.DestinationFragment
	)

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	hasServiceInstance := !filter.ServiceInstanceID.IsNull() && !filter.ServiceInstanceID.IsUnknown() && filter.ServiceInstanceID.ValueString() != ""

	if hasServiceInstance {
		cliRes, _, err = r.client.Connectivity.DestinationFragment.ListByServiceInstance(ctx, filter.SubaccountID.ValueString(), filter.ServiceInstanceID.ValueString())
	} else {
		cliRes, _, err = r.client.Connectivity.DestinationFragment.ListBySubaccount(ctx, filter.SubaccountID.ValueString())
	}

	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Destination Fragment (Subaccount)",
			fmt.Sprintf("Failed to list destination fragments: %s", err),
		)
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, fragment := range cliRes {

			name, ok := fragment.Content["FragmentName"]
			if !ok {
				continue
			}

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountID)
			result.Identity.SetAttribute(ctx, path.Root("name"), types.StringValue(name))
			result.Identity.SetAttribute(ctx, path.Root("service_instance_id"), filter.ServiceInstanceID)

			if req.IncludeResource {
				contentValue, diags := types.MapValueFrom(ctx, types.StringType, fragment.Content)
				result.Diagnostics.Append(diags...)

				resDestination := subaccountDestinationFragmentResourceConfig{
					SubaccountID:        filter.SubaccountID,
					Name:                types.StringValue(name),
					ServiceInstanceID:   filter.ServiceInstanceID,
					ID:                  types.StringValue(name),
					DestinationFragment: contentValue,
				}

				// 5. Set the resource information on the result
				if !result.Diagnostics.HasError() {
					result.Diagnostics.Append(result.Resource.Set(ctx, resDestination)...)
				}
			}

			// Push to the stream
			if !push(result) {
				return
			}
		}
	}
}
