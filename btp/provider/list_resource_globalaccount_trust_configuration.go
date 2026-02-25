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

var _ list.ListResourceWithConfigure = &globalaccountTrustConfigurationListResource{}

type globalaccountTrustConfigurationListResource struct {
	client *btpcli.ClientFacade
}

func NewGlobalaccountTrustConfigurationListResource() list.ListResource {
	return &globalaccountTrustConfigurationListResource{}
}

func (r *globalaccountTrustConfigurationListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_globalaccount_trust_configuration" // must match managed resource
}

func (r *globalaccountTrustConfigurationListResource) Configure(_ context.Context,
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

func (r *globalaccountTrustConfigurationListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all trust configurations available for given global account.",
	}
}

// List streams all trust configurations for global account from the API
func (r *globalaccountTrustConfigurationListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {
	cliRes, _, err := r.client.Security.Trust.ListByGlobalAccount(ctx)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Trust Configuration (Global Account)",
			fmt.Sprintf("Failed to list trust configurations: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, trustConfig := range cliRes {

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("origin"), types.StringValue(trustConfig.OriginKey))

			if req.IncludeResource {
				resServiceBinding, diags := globalaccountTrustConfigurationFromValue(ctx, trustConfig)

				result.Diagnostics.Append(diags...)

				// Set the resource information on the result
				if !result.Diagnostics.HasError() {
					result.Diagnostics.Append(result.Resource.Set(ctx, resServiceBinding)...)
				}
			}

			if !push(result) {
				return
			}
		}
	}
}
