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

var _ list.ListResourceWithConfigure = &subaccountRoleCollectionBaseListType{}

type subaccountRoleCollectionBaseListType struct {
	client *btpcli.ClientFacade
}

type subaccountRoleCollectionBaseFilterType struct {
	SubaccountID types.String `tfsdk:"subaccount_id"`
}

func NewSubaccountRoleCollectionBaseListResource() list.ListResource {
	return &subaccountRoleCollectionBaseListType{}
}

func (r *subaccountRoleCollectionBaseListType) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_role_collection_base" // must match managed resource
}

func (r *subaccountRoleCollectionBaseListType) Configure(_ context.Context,
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

func (r *subaccountRoleCollectionBaseListType) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	// This list resource takes the directory ID as input to filter directory role collections.
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all role collections base available within a specific BTP Directory. It requires the subaccount ID as input.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
		},
	}
}

// List streams directory role collections from the API
func (r *subaccountRoleCollectionBaseListType) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	var filter subaccountRoleCollectionBaseFilterType

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	roleCollections, _, err := r.client.Security.RoleCollection.ListBySubaccount(ctx, filter.SubaccountID.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Role Collection (Subaccount)",
			fmt.Sprintf("Failed to list directory role collections: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, value := range roleCollections {

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("name"), types.StringValue(value.Name))
			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountID)

			if req.IncludeResource {

				resRoleCollection := &subaccountRoleCollectionBaseType{
					Name:         types.StringValue(value.Name),
					ID:           types.StringValue(fmt.Sprintf("%s,%s", filter.SubaccountID.ValueString(), value.Name)),
					SubaccountId: filter.SubaccountID,
					Description:  types.StringValue(value.Description),
				}

				// Set the resource information on the result
				result.Diagnostics.Append(result.Resource.Set(ctx, resRoleCollection)...)
			}

			if !push(result) {
				return
			}
		}
	}
}
