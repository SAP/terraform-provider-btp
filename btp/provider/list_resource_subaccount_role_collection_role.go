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

var _ list.ListResourceWithConfigure = &subaccountRoleCollectionRoleListResource{}

type subaccountRoleCollectionRoleListResource struct {
	client *btpcli.ClientFacade
}

type subaccountRoleCollectionRoleFilterType struct {
	Name         types.String `tfsdk:"name"`
	SubaccountID types.String `tfsdk:"subaccount_id"`
}

func NewSubaccountRoleCollectionRoleListResource() list.ListResource {
	return &subaccountRoleCollectionRoleListResource{}
}

func (r *subaccountRoleCollectionRoleListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_role_collection_role" // must match managed resource
}

func (r *subaccountRoleCollectionRoleListResource) Configure(_ context.Context,
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

func (r *subaccountRoleCollectionRoleListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all role collection roles available within the configured BTP subaccount.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The Name of the role collection.",
				Required:            true,
			},
		},
	}
}

// List streams all subaccount role collection roles from the API
func (r *subaccountRoleCollectionRoleListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	var filter subaccountRoleCollectionRoleFilterType

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cliRes, _, err := r.client.Security.RoleCollection.GetBySubaccount(ctx, filter.SubaccountID.ValueString(), filter.Name.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Role Collection Role (Subaccount)",
			fmt.Sprintf("Failed to list role collection roles: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, value := range cliRes.RoleReferences {
			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountID)
			result.Identity.SetAttribute(ctx, path.Root("name"), cliRes.Name)
			result.Identity.SetAttribute(ctx, path.Root("role_name"), types.StringValue(value.Name))
			result.Identity.SetAttribute(ctx, path.Root("role_template_name"), types.StringValue(value.RoleTemplateName))
			result.Identity.SetAttribute(ctx, path.Root("role_template_app_id"), types.StringValue(value.RoleTemplateAppId))

			if req.IncludeResource {
				res := &subaccountRoleAssignmentType{
					Name:         types.StringValue(cliRes.Name),
					SubaccountId: filter.SubaccountID,
					ID: types.StringValue(fmt.Sprintf("%s,%s,%s,%s,%s",
						filter.SubaccountID.ValueString(),
						cliRes.Name,
						value.Name,
						value.RoleTemplateAppId,
						value.RoleTemplateName,
					)),
					RoleTemplateName:  types.StringValue(value.RoleTemplateName),
					RoleTemplateAppID: types.StringValue(value.RoleTemplateAppId),
					RoleName:          types.StringValue(value.Name),
				}

				// Set the resource information on the result
				result.Diagnostics.Append(result.Resource.Set(ctx, res)...)
			}

			if !push(result) {
				return
			}
		}
	}
}
