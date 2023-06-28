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

func newSubaccountServiceBrokersDataSource() datasource.DataSource {
	return &subaccountServiceBrokersDataSource{}
}

type subaccountServiceBrokerValue struct {
	Id           types.String `tfsdk:"id" btpcli:"id,get"`
	Name         types.String `tfsdk:"name" btpcli:"name,get"`
	Ready        types.Bool   `tfsdk:"ready"`
	Description  types.String `tfsdk:"description"`
	BrokerUrl    types.String `tfsdk:"broker_url"`
	CreatedDate  types.String `tfsdk:"created_date"`
	LastModified types.String `tfsdk:"last_modified"`
	Labels       types.Map    `tfsdk:"labels"`
}

type subaccountServiceBrokersDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	FieldsFilter types.String `tfsdk:"fields_filter"`
	LabelsFilter types.String `tfsdk:"labels_filter"`
	/* OUTPUT */
	Values []subaccountServiceBrokerValue `tfsdk:"values"`
}

type subaccountServiceBrokersDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServiceBrokersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_brokers", req.ProviderTypeName)
}

func (ds *subaccountServiceBrokersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServiceBrokersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists all service brokers in a subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
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
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service broker.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the service broker.",
							Computed:            true,
						},
						"ready": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the service broker is ready.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the service broker.",
							Computed:            true,
						},
						"broker_url": schema.StringAttribute{
							MarkdownDescription: "The URL of the service broker.",
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
						"labels": schema.MapAttribute{
							ElementType: types.SetType{
								ElemType: types.StringType,
							},
							MarkdownDescription: "Set of words or phrases assigned to the service broker.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountServiceBrokersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServiceBrokersDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var fieldsFilter, labelsFilter string
	if !data.FieldsFilter.IsNull() {
		fieldsFilter = data.FieldsFilter.ValueString()
	}
	if !data.LabelsFilter.IsNull() {
		labelsFilter = data.LabelsFilter.ValueString()
	}

	cliRes, _, err := ds.cli.Services.Broker.List(ctx, data.SubaccountId.ValueString(), fieldsFilter, labelsFilter)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Brokers (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Values = []subaccountServiceBrokerValue{}

	for _, broker := range cliRes {
		brokerValue := subaccountServiceBrokerValue{
			Id:           types.StringValue(broker.Id),
			Name:         types.StringValue(broker.Name),
			Ready:        types.BoolValue(broker.Ready),
			Description:  types.StringValue(broker.Description),
			BrokerUrl:    types.StringValue(broker.BrokerUrl),
			CreatedDate:  timeToValue(broker.CreatedAt),
			LastModified: timeToValue(broker.UpdatedAt),
		}

		brokerValue.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, broker.Labels)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, brokerValue)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
