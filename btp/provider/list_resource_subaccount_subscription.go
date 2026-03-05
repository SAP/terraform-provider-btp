package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ list.ListResourceWithConfigure = &subaccountSubscriptionListResource{}

type subaccountSubscriptionListResource struct {
	client *btpcli.ClientFacade
}

type subaccountSubscriptionListResourceFilter struct {
	SubaccountID types.String `tfsdk:"subaccount_id"`
}

func NewsubaccountSubscriptionListResource() list.ListResource {
	return &subaccountSubscriptionListResource{}
}

func (r *subaccountSubscriptionListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_subscription" // must match managed resource
}

func (r *subaccountSubscriptionListResource) Configure(_ context.Context,
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

func (r *subaccountSubscriptionListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all subscriptions available for given subaccount.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
		},
	}
}

// List streams all subscriptions available for given subaccount from the API
func (r *subaccountSubscriptionListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {
	var (
		filter subaccountSubscriptionListResourceFilter
	)

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cliRes, _, err := r.client.Accounts.Subscription.List(ctx, filter.SubaccountID.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Subscription (SubAccount)",
			fmt.Sprintf("Failed to list subscriptions: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, sub := range cliRes {

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountID)
			result.Identity.SetAttribute(ctx, path.Root("app_name"), types.StringValue(sub.AppName))
			result.Identity.SetAttribute(ctx, path.Root("plan_name"), types.StringValue(sub.PlanName))

			if req.IncludeResource {
				resSub, diags := subaccountSubscriptionValueFrom(ctx, sub)
				resSub.Timeouts = newNullTimeouts()

				result.Diagnostics.Append(diags...)

				// Set the resource information on the result
				if !result.Diagnostics.HasError() {
					result.Diagnostics.Append(result.Resource.Set(ctx, resSub)...)
				}
			}

			if !push(result) {
				return
			}
		}
	}
}

func newNullTimeouts() timeouts.Value {
	timeoutAttrTypes := map[string]attr.Type{
		"create": types.StringType,
		"delete": types.StringType,
		"update": types.StringType,
	}

	return timeouts.Value{
		Object: types.ObjectNull(timeoutAttrTypes),
	}
}
