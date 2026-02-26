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

var _ list.ListResourceWithConfigure = &subaccountSecuritySettingsListResource{}

type subaccountSecuritySettingsListResource struct {
	client *btpcli.ClientFacade
}

func NewSubaccountSecuritySettingsListResource() list.ListResource {
	return &subaccountSecuritySettingsListResource{}
}

type subaccountSecuritySettingsListResourceFilter struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
}

func (r *subaccountSecuritySettingsListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_security_settings" // must match managed resource
}

func (r *subaccountSecuritySettingsListResource) Configure(_ context.Context,
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

func (r *subaccountSecuritySettingsListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all security settings available for given subaccount_id.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
		},
	}
}

// List streams all security settings for given subaccount from the API
func (r *subaccountSecuritySettingsListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	var (
		filter subaccountSecuritySettingsListResourceFilter
	)

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cliRes, _, err := r.client.Security.Settings.ListBySubaccount(ctx, filter.SubaccountId.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Security Settings (Subaccount)",
			fmt.Sprintf("Failed to list security settings: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		result := req.NewListResult(ctx)

		result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountId)

		if req.IncludeResource {
			resSecuritySt, diags := subaccountSecuritySettingsValueFrom(ctx, cliRes, true)
			resSecuritySt.SubaccountId = filter.SubaccountId
			resSecuritySt.Id = types.StringValue(filter.SubaccountId.ValueString())

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
