package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ list.ListResourceWithConfigure = &subaccountDestinationGenericListResource{}

type subaccountDestinationGenericListResource struct {
	client *btpcli.ClientFacade
}

type subaccountDestinationGenericListResourceFilter struct {
	SubaccountID      types.String `tfsdk:"subaccount_id"`
	ServiceInstanceID types.String `tfsdk:"service_instance_id"`
}

func NewSubaccountDestinationGenericListResource() list.ListResource {
	return &subaccountDestinationGenericListResource{}
}

func (r *subaccountDestinationGenericListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_destination_generic" // must match managed resource
}

func (r *subaccountDestinationGenericListResource) Configure(_ context.Context,
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

func (r *subaccountDestinationGenericListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all destinations available for given subaccount.",
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

// List streams all destinations available for given subaccount from the API
func (r *subaccountDestinationGenericListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {
	var (
		filter subaccountDestinationGenericListResourceFilter
	)

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cliRes, _, err := r.client.Connectivity.Destination.ListBySubaccount(ctx, filter.SubaccountID.ValueString(), filter.ServiceInstanceID.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Destination Generic (SubAccount)",
			fmt.Sprintf("Failed to list destinations generic: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, destination := range cliRes {

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountID)
			result.Identity.SetAttribute(ctx, path.Root("name"), types.StringValue(destination.DestinationConfiguration["Name"]))
			result.Identity.SetAttribute(ctx, path.Root("service_instance_id"), filter.ServiceInstanceID)

			if req.IncludeResource {
				resDestination, diags := destinationGenericResourceValueFrom(destination, filter.SubaccountID, filter.ServiceInstanceID, destination.DestinationConfiguration["Name"])

				result.Diagnostics.Append(diags...)

				// Set the resource information on the result
				if !result.Diagnostics.HasError() {
					result.Diagnostics.Append(result.Resource.Set(ctx, resDestination)...)
				}
			}

			if !push(result) {
				return
			}
		}
	}
}
