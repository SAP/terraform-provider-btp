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

var _ list.ListResourceWithConfigure = &subaccountRoleListResource{}

type subaccountRoleListResource struct {
	client *btpcli.ClientFacade
}

type subaccountRoleListResourceFilter struct {
	SubaccountID types.String `tfsdk:"subaccount_id"`
}

func NewSubaccountRoleListResource() list.ListResource {
	return &subaccountRoleListResource{}
}

func (r *subaccountRoleListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_role" // must match managed resource
}

func (r *subaccountRoleListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *subaccountRoleListResource) ListResourceConfigSchema(_ context.Context, req list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all subaccount roles.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
		},
	}
}

// List streams all subaccount roles from the API
func (r *subaccountRoleListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var filter subaccountRoleListResourceFilter

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cliRes, _, err := r.client.Security.Role.ListBySubaccount(ctx, filter.SubaccountID.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Subaccount Roles",
			fmt.Sprintf("Failed to list subaccount roles: %s", err),
		)
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, subaccountRole := range cliRes {
			result := req.NewListResult(ctx)
			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountID)
			result.Identity.SetAttribute(ctx, path.Root("name"), types.StringValue(subaccountRole.Name))
			result.Identity.SetAttribute(ctx, path.Root("role_template_name"), types.StringValue(subaccountRole.RoleTemplateName))
			result.Identity.SetAttribute(ctx, path.Root("app_id"), types.StringValue(subaccountRole.RoleTemplateAppId))
			if req.IncludeResource {
				resRole, diags := subaccountRoleFromValue(ctx, subaccountRole)
				result.Diagnostics.Append(diags...)

				if !result.Diagnostics.HasError() {
					resRole.SubaccountId = filter.SubaccountID
					compositeID := fmt.Sprintf("%s,%s,%s,%s",
						filter.SubaccountID.ValueString(),
						subaccountRole.Name,
						subaccountRole.RoleTemplateName,
						subaccountRole.RoleTemplateAppId,
					)
					resRole.Id = types.StringValue(compositeID)
					result.Diagnostics.Append(result.Resource.Set(ctx, resRole)...)
				}
			}
			if !push(result) {
				return
			}
		}
	}
}
