package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ list.ListResourceWithConfigure = &subaccountListResource{}

type subaccountListResource struct {
	client *btpcli.ClientFacade
}

type subaccountListResourceFilter struct {
	Region       types.String `tfsdk:"region"`
	LabelsFilter types.String `tfsdk:"labels_filter"`
}

func NewSubaccountListResource() list.ListResource {
	return &subaccountListResource{}
}

func (r *subaccountListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount" // must match managed resource
}

func (r *subaccountListResource) Configure(_ context.Context,
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

func (r *subaccountListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all subaccounts. The results can be filtered using `region` or `labels_filter`.",
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				MarkdownDescription: "The region of the subaccount. For example, `eu12`.",
				Optional:            true,
			},
			"labels_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the subaccount based on the label query.  For example, to list all subaccount that are available in a production landscape, use \"landscape=production\".",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

// List streams all subaccount from the API
func (r *subaccountListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	var (
		filter       subaccountListResourceFilter
		labelsFilter string
	)

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	if !filter.LabelsFilter.IsNull() {
		labelsFilter = filter.LabelsFilter.ValueString()
	}

	cliRes, _, err := r.client.Accounts.Subaccount.List(ctx, labelsFilter)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Subaccount",
			fmt.Sprintf("Failed to list subaccounts: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, subaccount := range cliRes.Value {

			if !filter.Region.IsNull() && subaccount.Region != filter.Region.ValueString() {
				continue
			}

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), subaccount.Guid)

			if req.IncludeResource {
				resSubaccount, diags := subaccountListValueFrom(ctx, subaccount)

				result.Diagnostics.Append(diags...)

				// Set the resource information on the result
				if !result.Diagnostics.HasError() {
					result.Diagnostics.Append(result.Resource.Set(ctx, resSubaccount)...)
				}
			}

			if !push(result) {
				return
			}
		}
	}
}
