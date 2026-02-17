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

var _ list.ListResourceWithConfigure = &directoryEntitlementListResource{}

type directoryEntitlementListResource struct {
	client *btpcli.ClientFacade
}

type directoryEntitlementFilterType struct {
	DirectoryID types.String `tfsdk:"directory_id"`
}

func NewDirectoryEntitlementListResource() list.ListResource {
	return &directoryEntitlementListResource{}
}

func (r *directoryEntitlementListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directory_entitlement" // must match managed resource
}

func (r *directoryEntitlementListResource) Configure(_ context.Context,
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

func (r *directoryEntitlementListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	// This list resource takes the directory ID as input to filter entitlements.
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all entitlements available within a specific BTP Directory. It requires the directory ID as input.",
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
			},
		},
	}
}

// List streams all directory entitlements from the API
func (r *directoryEntitlementListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	var filter directoryEntitlementFilterType

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cliRes, _, err := r.client.Accounts.Entitlement.ListByDirectory(ctx, filter.DirectoryID.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Role (Global Account)",
			fmt.Sprintf("Failed to list roles: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, service := range cliRes.EntitledServices {
			for _, servicePlan := range service.ServicePlans {
				result := req.NewListResult(ctx)

				result.Identity.SetAttribute(ctx, path.Root("directory_id"), filter.DirectoryID.ValueString())
				result.Identity.SetAttribute(ctx, path.Root("service_name"), types.StringValue(service.Name))
				result.Identity.SetAttribute(ctx, path.Root("plan_name"), types.StringValue(servicePlan.Name))

				if req.IncludeResource {
					resDm := &directoryEntitlementType{
						DirectoryId:          filter.DirectoryID,
						Id:                   types.StringValue(servicePlan.UniqueIdentifier),
						ServiceName:          types.StringValue(service.Name),
						PlanName:             types.StringValue(servicePlan.Name),
						PlanUniqueIdentifier: types.StringValue(servicePlan.UniqueIdentifier),
						Amount:               types.Int64Value(int64(servicePlan.Amount)),
						AutoAssign:           types.BoolValue(servicePlan.AutoAssign),
						AutoDistributeAmount: types.Int64Value(int64(servicePlan.AutoDistributeAmount)),
						Category:             types.StringValue(servicePlan.Category),
						PlanId:               types.StringValue(servicePlan.UniqueIdentifier),
						Distribute:           types.BoolValue(false),
					}

					// Set the resource information on the result
					result.Diagnostics.Append(result.Resource.Set(ctx, resDm)...)
				}

				if !push(result) {
					return
				}
			}

		}
	}
}
