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

var _ list.ListResourceWithConfigure = &subaccountRoleCollectionListType{}

type subaccountRoleCollectionListType struct {
	client *btpcli.ClientFacade
}

type subaccountRoleCollectionFilterType struct {
	SubaccountID types.String `tfsdk:"subaccount_id"`
}

func NewSubaccountRoleCollectionListResource() list.ListResource {
	return &subaccountRoleCollectionListType{}
}

func (r *subaccountRoleCollectionListType) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_role_collection" // must match managed resource
}

func (r *subaccountRoleCollectionListType) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *subaccountRoleCollectionListType) ListResourceConfigSchema(_ context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all role collections available within a specific BTP subaccount. It requires the subaccount ID as input.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
		},
	}
}

// List streams subaccount role collections from the API
func (r *subaccountRoleCollectionListType) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {

	var filter subaccountRoleCollectionFilterType

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cliRes, _, err := r.client.Security.RoleCollection.ListBySubaccount(ctx, filter.SubaccountID.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Subaccount Role Collection",
			fmt.Sprintf("Failed to list subaccount role collections: %s", err),
		)
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	stream.Results = func(push func(list.ListResult) bool) {

		for _, value := range cliRes {
			result := req.NewListResult(ctx)
			result.Identity.SetAttribute(ctx, path.Root("name"), types.StringValue(value.Name))
			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountID)

			if req.IncludeResource {
				roles := []subaccountRoleCollectionRoleRefType{}
				for _, role := range value.RoleReferences {
					roles = append(roles, subaccountRoleCollectionRoleRefType{
						RoleTemplateName:  types.StringValue(role.RoleTemplateName),
						RoleTemplateAppId: types.StringValue(role.RoleTemplateAppId),
						Name:              types.StringValue(role.Name),
					})
				}

				resDir := &subaccountRoleCollectionType{
					Name:         types.StringValue(value.Name),
					Id:           types.StringValue(fmt.Sprintf("%s,%s", filter.SubaccountID.ValueString(), value.Name)),
					SubaccountId: filter.SubaccountID,
					Description:  types.StringValue(value.Description),
					Roles:        roles,
				}
				result.Diagnostics.Append(result.Resource.Set(ctx, resDir)...)
			}

			if !push(result) {
				return
			}
		}
	}
}
