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

var _ list.ListResourceWithConfigure = &directoryRoleListResource{}

type directoryRoleListResource struct {
	client *btpcli.ClientFacade
}

type directoryRoleListResourceFilter struct {
	DirectoryID types.String `tfsdk:"directory_id"`
}

func NewDirectoryRoleListResource() list.ListResource {
	return &directoryRoleListResource{}
}

func (r *directoryRoleListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directory_role" // must match managed resource
}

func (r *directoryRoleListResource) Configure(_ context.Context,
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

func (r *directoryRoleListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all directory roles.",
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
			},
		},
	}
}

// List streams all directory roles from the API
func (r *directoryRoleListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {
	var (
		filter directoryRoleListResourceFilter
	)

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cliRes, _, err := r.client.Security.Role.ListByDirectory(ctx, filter.DirectoryID.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Directory Roles",
			fmt.Sprintf("Failed to list directory roles: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, dirRole := range cliRes {
			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("directory_id"), types.StringValue(dirRole.Name))
			result.Identity.SetAttribute(ctx, path.Root("name"), types.StringValue(dirRole.RoleTemplateName))
			result.Identity.SetAttribute(ctx, path.Root("role_template_name"), types.StringValue(dirRole.RoleTemplateAppId))
			result.Identity.SetAttribute(ctx, path.Root("app_id"), types.StringValue(dirRole.RoleTemplateAppId))

			if req.IncludeResource {
				resDirRole, diags := directoryRoleFromValue(ctx, dirRole)

				result.Diagnostics.Append(diags...)

				// Set the resource information on the result
				if !result.Diagnostics.HasError() {
					result.Diagnostics.Append(result.Resource.Set(ctx, resDirRole)...)
				}
			}

			if !push(result) {
				return
			}
		}
	}
}
