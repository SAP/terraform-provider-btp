package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ list.ListResourceWithConfigure = &subaccountServiceInstanceListResource{}

type subaccountServiceInstanceListResource struct {
	client *btpcli.ClientFacade
}

type subaccountServiceInstanceListResourceFilter struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	FieldsFilter types.String `tfsdk:"fields_filter"`
	LabelsFilter types.String `tfsdk:"labels_filter"`
}

func NewSubaccountServiceInstanceListResource() list.ListResource {
	return &subaccountServiceInstanceListResource{}
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
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all service instances available for given subaccount_id. The results can be filtered using `fields_filter` or `labels_filter`.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"fields_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the instances based on their fields. For example, to list all instances that are usable, use \"usable eq 'true'\".",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"labels_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the instances based on the label query.  For example, to list all instances that are available in a production landscape, use \"landscape eq 'production'\".",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

// List streams all service instances for given subaccount from the API
func (r *subaccountServiceInstanceListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	var (
		filter                     subaccountServiceInstanceListResourceFilter
		fieldsFilter, labelsFilter string
	)

	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	if !filter.FieldsFilter.IsNull() {
		fieldsFilter = filter.FieldsFilter.ValueString()
	}

	if !filter.LabelsFilter.IsNull() {
		labelsFilter = filter.LabelsFilter.ValueString()
	}

	cliRes, _, err := r.client.Services.Instance.List(ctx, filter.SubaccountId.ValueString(), fieldsFilter, labelsFilter)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Service Instance (Subaccount)",
			fmt.Sprintf("Failed to list service instances: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, serviceInstance := range cliRes {

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountId)
			result.Identity.SetAttribute(ctx, path.Root("id"), types.StringValue(serviceInstance.Id))

			if req.IncludeResource {
				serviceInstance, diags := subaccountServiceInstanceListValueFrom(ctx, serviceInstance)

				result.Diagnostics.Append(diags...)

				// Set the resource information on the result
				if !result.Diagnostics.HasError() {
					result.Diagnostics.Append(result.Resource.Set(ctx, serviceInstance)...)
				}
			}

			if !push(result) {
				return
			}
		}
	}
}
