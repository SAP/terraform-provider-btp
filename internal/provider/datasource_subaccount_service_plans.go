package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountServicePlansDataSource() datasource.DataSource {
	return &subaccountServicePlansDataSource{}
}

type subaccountServicePlanValueConfig struct {
	/* INPUT */
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	/* OUTPUT */
	Ready             types.Bool   `tfsdk:"ready"`
	Description       types.String `tfsdk:"description"`
	CatalogId         types.String `tfsdk:"catalog_id"`
	CatalogName       types.String `tfsdk:"catalog_name"`
	Free              types.Bool   `tfsdk:"free"`
	Bindable          types.Bool   `tfsdk:"bindable"`
	ServiceOfferingId types.String `tfsdk:"serviceoffering_id"`
	CreatedDate       types.String `tfsdk:"created_date"`
	LastModified      types.String `tfsdk:"last_modified"`
	Shareable         types.Bool   `tfsdk:"shareable"`
}

type subaccountServicePlansDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	Environment  types.String `tfsdk:"environment"`
	FieldsFilter types.String `tfsdk:"fields_filter"`
	LabelsFilter types.String `tfsdk:"labels_filter"`
	/* OUTPUT */
	Values []subaccountServicePlanValueConfig `tfsdk:"values"`
}

type subaccountServicePlansDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServicePlansDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_plans", req.ProviderTypeName)
}

func (ds *subaccountServicePlansDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServicePlansDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists the plans of services that your subaccount is entitled to use in your environment.`,
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
				MarkdownDescription: "Filter the response on the environment (sapbtp, kubernetes, cloudfoundry).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"fields_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the response based on the field query.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"labels_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the response based on the labels query.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service plan.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the service plan.",
							Computed:            true,
						},
						"ready": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the service plan is ready.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the service plan.",
							Computed:            true,
						},
						"catalog_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service plan in the service broker catalog.",
							Computed:            true,
						},
						"catalog_name": schema.StringAttribute{
							MarkdownDescription: "The name of the associated service broker catalog.",
							Computed:            true,
						},
						"free": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the service plan is free.",
							Computed:            true,
						},
						"bindable": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the service plan is bindable.",
							Computed:            true,
						},
						"serviceoffering_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service offering.",
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
						"shareable": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the service plan supports instance sharing.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountServicePlansDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServicePlansDataSourceConfig

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

	cliRes, _, err := ds.cli.Services.Plan.List(ctx, data.SubaccountId.ValueString(), fieldsFilter, labelsFilter, environment)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Plans (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.SubaccountId
	data.Values = []subaccountServicePlanValueConfig{}

	for _, item := range cliRes {

		shareable := false

		if item.Metadata != nil {
			shareable = item.Metadata.SupportsInstanceSharing
		}

		value := subaccountServicePlanValueConfig{
			Id:                types.StringValue(item.Id),
			Name:              types.StringValue(item.Name),
			Ready:             types.BoolValue(item.Ready),
			Description:       types.StringValue(item.Description),
			CatalogId:         types.StringValue(item.CatalogId),
			CatalogName:       types.StringValue(item.CatalogName),
			Free:              types.BoolValue(item.Free),
			Bindable:          types.BoolValue(item.Bindable),
			ServiceOfferingId: types.StringValue(item.ServiceOfferingId),
			CreatedDate:       timeToValue(item.CreatedAt),
			LastModified:      timeToValue(item.UpdatedAt),
			Shareable:         types.BoolValue(shareable),
		}

		data.Values = append(data.Values, value)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
