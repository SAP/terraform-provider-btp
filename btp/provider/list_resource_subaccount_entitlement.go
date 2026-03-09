package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis_entitlements"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ list.ListResourceWithConfigure = &subaccountEntitlementListResource{}

type subaccountEntitlementListResource struct {
	client *btpcli.ClientFacade
}

type subaccountEntitlementFilterType struct {
	SubaccountID types.String `tfsdk:"subaccount_id"`
}

func NewSubaccountEntitlementListResource() list.ListResource {
	return &subaccountEntitlementListResource{}
}

func (r *subaccountEntitlementListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_entitlement" // must match managed resource
}

func (r *subaccountEntitlementListResource) Configure(_ context.Context,
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

func (r *subaccountEntitlementListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	// This list resource takes the subaccount ID as input to filter entitlements.
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all entitlements available within a specific BTP subaccount. It requires the subaaccount ID as input.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
		},
	}
}

// List streams subaccount entitlements from the API
func (r *subaccountEntitlementListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	var (
		filter subaccountEntitlementFilterType
		cliRes cis_entitlements.EntitledAndAssignedServicesResponseObject
	)

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	// Determine the parent of the subaccount
	// In case of a directory with feature "ENTITLEMENTS" enabled we must hand over the ID in the "List" call
	subaccountData, _, err := r.client.Accounts.Subaccount.Get(ctx, filter.SubaccountID.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Entitlement (Subaccount)",
			fmt.Sprintf("Failed to retrieve subaccount details: %s", err),
		)
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	parentId, isParentGlobalAccount, err := determineParentIdForEntitlement(r.client, ctx, subaccountData.ParentGUID)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Entitlement (Subaccount)",
			fmt.Sprintf("Failed to list subaccount entitlements: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	if isParentGlobalAccount {
		cliRes, _, err = r.client.Accounts.Entitlement.ListBySubaccount(ctx, filter.SubaccountID.ValueString())
	} else {
		cliRes, _, err = r.client.Accounts.Entitlement.ListBySubaccountWithDirectoryParent(ctx, filter.SubaccountID.ValueString(), parentId)
	}

	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Entitlement (Subaccount)",
			fmt.Sprintf("Failed to list subaccount entitlements: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, service := range cliRes.EntitledServices {
			for _, servicePlan := range service.ServicePlans {
				result := req.NewListResult(ctx)

				result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountID)
				result.Identity.SetAttribute(ctx, path.Root("service_name"), types.StringValue(service.Name))
				result.Identity.SetAttribute(ctx, path.Root("plan_name"), types.StringValue(servicePlan.Name))

				if req.IncludeResource {
					resEnt := subaccountEntitlementListValueFrom(service, servicePlan, filter.SubaccountID.ValueString())

					// Set the resource information on the result
					result.Diagnostics.Append(result.Resource.Set(ctx, resEnt)...)
				}

				if !push(result) {
					return
				}
			}

		}
	}
}
