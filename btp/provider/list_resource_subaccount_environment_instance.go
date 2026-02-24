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

var _ list.ListResourceWithConfigure = &subaccountEnvironmentInstanceListResource{}

type subaccountEnvironmentInstanceListResource struct {
	client *btpcli.ClientFacade
}

type subaccountEnvironmentInstanceListResourceFilter struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
}

func NewSubaccountEnvironmentInstanceListResource() list.ListResource {
	return &subaccountEnvironmentInstanceListResource{}
}

func (r *subaccountEnvironmentInstanceListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_environment_instance" // must match managed resource
}

func (r *subaccountEnvironmentInstanceListResource) Configure(_ context.Context,
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

func (r *subaccountEnvironmentInstanceListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all environment instances for a subaccount.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
		},
	}
}

// List streams all environment instances for a subaccount from the API
func (r *subaccountEnvironmentInstanceListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	var (
		filter subaccountEnvironmentInstanceListResourceFilter
	)

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cliRes, _, err := r.client.Accounts.EnvironmentInstance.List(ctx, filter.SubaccountId.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Environment Instance (Subaccount)",
			fmt.Sprintf("Failed to list environment instances: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, environmentInstance := range cliRes.EnvironmentInstances {

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountId)
			result.Identity.SetAttribute(ctx, path.Root("id"), types.StringValue(environmentInstance.Id))

			if req.IncludeResource {
				resEnvironmentInstance, diags := subaccountEnvironmentInstanceListValueFrom(ctx, environmentInstance)

				result.Diagnostics.Append(diags...)

				// Set the resource information on the result
				if !result.Diagnostics.HasError() {
					result.Diagnostics.Append(result.Resource.Set(ctx, resEnvironmentInstance)...)
				}
			}

			if !push(result) {
				return
			}
		}
	}
}
