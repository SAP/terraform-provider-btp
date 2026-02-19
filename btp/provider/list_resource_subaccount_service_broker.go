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

var _ list.ListResourceWithConfigure = &subaccountServiceBrokerListResource{}

type subaccountServiceBrokerListResource struct {
	client *btpcli.ClientFacade
}

func NewSubaccountServiceBrokerListResource() list.ListResource {
	return &subaccountServiceBrokerListResource{}
}

type subaccountServiceBrokerListResourceFilter struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	FieldsFilter types.String `tfsdk:"fields_filter"`
	LabelsFilter types.String `tfsdk:"labels_filter"`
}

func (r *subaccountServiceBrokerListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_service_broker" // must match managed resource
}

func (r *subaccountServiceBrokerListResource) Configure(_ context.Context,
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

func (r *subaccountServiceBrokerListResource) ListResourceConfigSchema(
	_ context.Context,
	req list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover all service brokers available within the configured subaccount. The results can be filtered using `fields_filter` or `labels_filter`.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"fields_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the service brokers based on their fields. For example, to display a service broker with the name 'my-service-broker2', use \"name eq 'my-service-broker2'\".",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"labels_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the service brokers based on the label query. For example, to display a service broker with the label 'country', whose value is 'France', use \"country eq 'France'\".",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

// List streams all service brokers for given subaccount from the API
func (r *subaccountServiceBrokerListResource) List(
	ctx context.Context,
	req list.ListRequest,
	stream *list.ListResultsStream,
) {

	var (
		filter                     subaccountServiceBrokerListResourceFilter
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

	cliRes, _, err := r.client.Services.Broker.List(ctx, filter.SubaccountId.ValueString(), fieldsFilter, labelsFilter)
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Service Broker (Subaccount)",
			fmt.Sprintf("Failed to list service brokers: %s", err),
		)

		stream.Results = list.ListResultsStreamDiagnostics(diags)

		return
	}

	stream.Results = func(push func(list.ListResult) bool) {

		for _, broker := range cliRes {

			result := req.NewListResult(ctx)

			result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), filter.SubaccountId)
			result.Identity.SetAttribute(ctx, path.Root("id"), types.StringValue(broker.Id))

			if req.IncludeResource {
				resServiceBroker := &subaccountServiceBrokerResourceType{
					Name:         types.StringValue(broker.Name),
					Description:  types.StringValue(broker.Description),
					Id:           types.StringValue(broker.Id),
					Url:          types.StringValue(broker.BrokerUrl),
					SubaccountId: filter.SubaccountId,
					Ready:        types.BoolValue(broker.Ready),
					CreatedDate:  timeToValue(broker.CreatedAt),
					LastModified: timeToValue(broker.UpdatedAt),
				}

				// Set the resource information on the result
				result.Diagnostics.Append(result.Resource.Set(ctx, resServiceBroker)...)
			}

			if !push(result) {
				return
			}
		}
	}
}
