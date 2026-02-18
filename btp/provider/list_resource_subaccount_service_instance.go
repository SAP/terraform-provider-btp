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

var _ list.ListResourceWithConfigure = &subaccountServiceInstanceListResource{}

type subaccountServiceInstanceListResource struct {
	client *btpcli.ClientFacade
}

func NewSubaccountServiceInstanceListResource() list.ListResource {
	return &GlobalaccountRoleCollectionListResource{}
}

func (r *subaccountServiceInstanceListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_service_instance" // must match managed resource
}

func (r *subaccountServiceInstanceListResource) Configure(_ context.Context,
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

func (r *subaccountServiceInstanceListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	// This list resource takes no input configurations (e.g. filters)
	// from the HCL, so the schema remains empty.
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all role collections available within the configured BTP Global Account. It does not require any input configuration filters.",
	}
}

// List streams all global account role collections from the API
func (r *subaccountServiceInstanceListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	cliRes, _, err := r.client.Security.RoleCollection.ListByGlobalAccount(ctx)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Role Collection (Global Account)",
			fmt.Sprintf("Failed to list role collections: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, value := range cliRes {

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("name"), types.StringValue(value.Name))

			if req.IncludeResource {
				resDm := &globalaccountRoleCollectionType{
					Name:        types.StringValue(value.Name),
					Description: types.StringValue(value.Description),
					Id:          types.StringValue(value.Name),
				}

				roles := []globalaccountRoleCollectionRoleRefType{}
				for _, role := range value.RoleReferences {
					roles = append(roles, globalaccountRoleCollectionRoleRefType{
						RoleTemplateName:  types.StringValue(role.RoleTemplateName),
						RoleTemplateAppId: types.StringValue(role.RoleTemplateAppId),
						Name:              types.StringValue(role.Name),
					})
				}

				resDm.Roles = roles

				// Set the resource information on the result
				result.Diagnostics.Append(result.Resource.Set(ctx, resDm)...)
			}

			if !push(result) {
				return
			}
		}
	}
}
