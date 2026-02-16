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

var _ list.ListResourceWithConfigure = &GlobalaccountRoleListResource{}

type GlobalaccountRoleListResource struct {
	client *btpcli.ClientFacade
}

func NewGlobalaccountRoleListResource() list.ListResource {
	return &GlobalaccountRoleListResource{}
}

func (r *GlobalaccountRoleListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_globalaccount_role" // must match managed resource
}

func (r *GlobalaccountRoleListResource) Configure(_ context.Context,
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

func (r *GlobalaccountRoleListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	// This list resource takes no input configurations (e.g. filters)
	// from the HCL, so the schema remains empty.
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all roles available within the configured BTP Global Account. It does not require any input configuration filters.",
	}
}

// List streams all global account roles from the API
func (r *GlobalaccountRoleListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	cliRes, _, err := r.client.Security.Role.ListByGlobalAccount(ctx)
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

		for _, value := range cliRes {

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("name"), types.StringValue(value.Name))
			result.Identity.SetAttribute(ctx, path.Root("role_template_name"), types.StringValue(value.RoleTemplateName))
			result.Identity.SetAttribute(ctx, path.Root("app_id"), types.StringValue(value.RoleTemplateAppId))

			if req.IncludeResource {
				resDm := &globalaccountRoleType{
					Name:              types.StringValue(value.Name),
					RoleTemplateAppId: types.StringValue(value.RoleTemplateAppId),
					RoleTemplateName:  types.StringValue(value.RoleTemplateName),
					IsReadOnly:        types.BoolValue(value.IsReadOnly),
					Description:       types.StringValue(value.Description),
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
