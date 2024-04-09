package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountServiceOfferingsDataSource() datasource.DataSource {
	return &subaccountServiceOfferingsDataSource{}
}

type subaccountServiceOfferingValue struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Ready                types.Bool   `tfsdk:"ready"`
	Description          types.String `tfsdk:"description"`
	Bindable             types.Bool   `tfsdk:"bindable"`
	InstancesRetrievable types.Bool   `tfsdk:"instances_retrievable"`
	BindingsRetrievable  types.Bool   `tfsdk:"bindings_retrievable"`
	PlanUpdateable       types.Bool   `tfsdk:"plan_updateable"`
	AllowContextUpdates  types.Bool   `tfsdk:"allow_context_updates"`
	Tags                 types.Set    `tfsdk:"tags"`
	BrokerId             types.String `tfsdk:"broker_id"`
	CatalogId            types.String `tfsdk:"catalog_id"`
	CatalogName          types.String `tfsdk:"catalog_name"`
	CreatedDate          types.String `tfsdk:"created_date"`
	LastModified         types.String `tfsdk:"last_modified"`
}

type subaccountServiceOfferingsDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	Environment  types.String `tfsdk:"environment"`
	FieldsFilter types.String `tfsdk:"fields_filter"`
	LabelsFilter types.String `tfsdk:"labels_filter"`
	/* OUTPUT */
	Values []subaccountServiceOfferingValue `tfsdk:"values"`
}

type subaccountServiceOfferingsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServiceOfferingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_offerings", req.ProviderTypeName)
}

func (ds *subaccountServiceOfferingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServiceOfferingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists the services your subaccount is entitled to use in your runtime environment.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.`,		
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` attribute instead",
				MarkdownDescription: "The ID of the subaccount.",
				Computed:            true,
			},
			"environment": schema.StringAttribute{
				MarkdownDescription: "Lists services to be consumed in a Cloud Foundry or Kubernetes-native way. Valid values are: \n " +
					getFormattedValueAsTableRow("value", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`cloudfoundry`", "Cloud Foundry") +
					getFormattedValueAsTableRow("`kubernetes`", "Kubernetes"),
				Optional: true,
			},
			"fields_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the response based on the field query. For example, use \"name eq 'my service offering name'\".",
				Optional:            true,
			},
			"labels_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the response based on the label query.  For example, to list all the service offerings associated with the testing environment, use \"environment eq 'test'\".",
				Optional:            true,
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service offering.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the service offering.",
							Computed:            true,
						},
						"ready": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the service offering is ready to be advertised.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the service offering.",
							Computed:            true,
						},
						"bindable": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the service offering is bindable.",
							Computed:            true,
						},
						"instances_retrievable": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the service instances associated with the service offering can be retrieved.",
							Computed:            true,
						},
						"bindings_retrievable": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the bindings associated with the service offering can be retrieved.",
							Computed:            true,
						},
						"plan_updateable": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the offered plan can be updated.",
							Computed:            true,
						},
						"allow_context_updates": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the context for the service offering can be updated.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "The list of tags for the service offering.",
							Computed:            true,
						},
						"broker_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the broker that provides the service plan.",
							Computed:            true,
						},
						"catalog_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service offering as provided by the catalog.",
							Computed:            true,
						},
						"catalog_name": schema.StringAttribute{
							MarkdownDescription: "The catalog name of the service offering.",
							Computed:            true,
						},
						"created_date": schema.StringAttribute{
							MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
							Computed:            true,
						},
						"last_modified": schema.StringAttribute{
							MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountServiceOfferingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServiceOfferingsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var fieldsFilter, labelsFilter, environment string
	if !data.FieldsFilter.IsNull() {
		fieldsFilter = data.FieldsFilter.ValueString()
	}
	if !data.LabelsFilter.IsNull() {
		labelsFilter = data.LabelsFilter.ValueString()
	}
	if !data.Environment.IsNull() {
		environment = data.Environment.ValueString()
	}

	cliRes, _, err := ds.cli.Services.Offering.List(ctx, data.SubaccountId.ValueString(), fieldsFilter, labelsFilter, environment)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Offerings (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.SubaccountId
	data.Values = []subaccountServiceOfferingValue{}

	for _, offering := range cliRes {
		offeringValue := subaccountServiceOfferingValue{
			Id:                   types.StringValue(offering.Id),
			Name:                 types.StringValue(offering.Name),
			Ready:                types.BoolValue(offering.Ready),
			Description:          types.StringValue(offering.Description),
			Bindable:             types.BoolValue(offering.Bindable),
			InstancesRetrievable: types.BoolValue(offering.InstancesRetrievable),
			BindingsRetrievable:  types.BoolValue(offering.BindingsRetrievable),
			PlanUpdateable:       types.BoolValue(offering.PlanUpdateable),
			AllowContextUpdates:  types.BoolValue(offering.AllowContextUpdates),
			BrokerId:             types.StringValue(offering.BrokerId),
			CatalogId:            types.StringValue(offering.CatalogId),
			CatalogName:          types.StringValue(offering.CatalogName),
			CreatedDate:          timeToValue(offering.CreatedAt),
			LastModified:         timeToValue(offering.UpdatedAt),
		}

		offeringValue.Tags, diags = types.SetValueFrom(ctx, types.StringType, offering.Tags)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, offeringValue)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
