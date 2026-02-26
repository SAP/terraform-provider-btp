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

var _ list.ListResourceWithConfigure = &globalaccountSecuritySettingsListResource{}

type globalaccountSecuritySettingsListResource struct {
	client *btpcli.ClientFacade
}

func NewGlobalaccountSecuritySettingsListResource() list.ListResource {
	return &globalaccountSecuritySettingsListResource{}
}

func (r *globalaccountSecuritySettingsListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_globalaccount_security_settings" // must match managed resource
}

func (r *globalaccountSecuritySettingsListResource) Configure(_ context.Context,
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

func (r *globalaccountSecuritySettingsListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all security settings available for given global account.",
	}
}

// List streams all security settings for given global account from the API
func (r *globalaccountSecuritySettingsListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	cliRes, _, err := r.client.Security.Settings.ListByGlobalAccount(ctx)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Security Settings (Global Account)",
			fmt.Sprintf("Failed to list security settings: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		result := req.NewListResult(ctx)

		result.Identity.SetAttribute(ctx, path.Root("globalaccount_subdomain"), types.StringValue(r.client.GetGlobalAccountSubdomain()))

		if req.IncludeResource {
			resSecuritySt, diags := globalaccountSecuritySettingsValueFrom(ctx, cliRes)
			resSecuritySt.Id = types.StringValue(r.client.GetGlobalAccountSubdomain())

			result.Diagnostics.Append(diags...)

			// Set the resource information on the result
			if !result.Diagnostics.HasError() {
				result.Diagnostics.Append(result.Resource.Set(ctx, resSecuritySt)...)
			}
		}

		if !push(result) {
			return
		}
	}
}
