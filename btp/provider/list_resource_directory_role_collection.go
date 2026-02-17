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

var _ list.ListResourceWithConfigure = &directoryRoleCollectionListType{}

type directoryRoleCollectionListType struct {
	client *btpcli.ClientFacade
}

type directoryRoleCollectionFilterType struct {
	DirectoryID types.String `tfsdk:"directory_id"`
}

func NewDirectoryRoleCollectionListResource() list.ListResource {
	return &directoryRoleCollectionListType{}
}

func (r *directoryRoleCollectionListType) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_directory_role_collection" // must match managed resource
}

func (r *directoryRoleCollectionListType) Configure(_ context.Context,
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

func (r *directoryRoleCollectionListType) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	// This list resource takes the directory ID as input to filter entitlements.
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all role collections available within a specific BTP Directory. It requires the directory ID as input.",
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
			},
		},
	}
}

// List streams all global account roles from the API
func (r *directoryRoleCollectionListType) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	var filter directoryRoleCollectionFilterType

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	cliRes, _, err := r.client.Security.RoleCollection.ListByDirectory(ctx, filter.DirectoryID.ValueString())
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
			result.Identity.SetAttribute(ctx, path.Root("directory_id"), filter.DirectoryID.ValueString())

			if req.IncludeResource {
				roles := []directoryRoleCollectionRoleRefType{}
				for _, role := range value.RoleReferences {
					roles = append(roles, directoryRoleCollectionRoleRefType{
						RoleTemplateName:  types.StringValue(role.RoleTemplateName),
						RoleTemplateAppId: types.StringValue(role.RoleTemplateAppId),
						Name:              types.StringValue(role.Name),
					})
				}

				resDir := &directoryRoleCollectionTypeConfig{
					Name:        types.StringValue(value.Name),
					Id:          types.StringValue(fmt.Sprintf("%s,%s", filter.DirectoryID.ValueString(), value.Name)),
					DirectoryId: filter.DirectoryID,
					Description: types.StringValue(value.Description),
					Roles:       roles,
				}

				// Set the resource information on the result
				result.Diagnostics.Append(result.Resource.Set(ctx, resDir)...)
			}

			if !push(result) {
				return
			}
		}
	}
}
