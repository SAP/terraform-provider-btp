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

func newSubaccountServiceInstancesDataSource() datasource.DataSource {
	return &subaccountServiceInstancesDataSource{}
}

type subaccountServiceInstancesValueConfig struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Ready         types.Bool   `tfsdk:"ready"`
	ServicePlanId types.String `tfsdk:"serviceplan_id"`
	PlatformId    types.String `tfsdk:"platform_id"`
	Context       types.String `tfsdk:"context"`
	Usable        types.Bool   `tfsdk:"usable"`
	CreatedDate   types.String `tfsdk:"created_date"`
	LastModified  types.String `tfsdk:"last_modified"`
	Labels        types.Map    `tfsdk:"labels"`
}

type subaccountServiceInstancesDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	FieldsFilter types.String `tfsdk:"fields_filter"`
	LabelsFilter types.String `tfsdk:"labels_filter"`
	/* OUTPUT */
	Values []subaccountServiceInstancesValueConfig `tfsdk:"values"`
}

type subaccountServiceInstancesDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServiceInstancesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_instances", req.ProviderTypeName)
}

func (ds *subaccountServiceInstancesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServiceInstancesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists all service instances in a subaccount.`,
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
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service instance.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the service instance.",
							Computed:            true,
						},
						"ready": schema.BoolAttribute{
							MarkdownDescription: "",
							Computed:            true,
						},
						"serviceplan_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service plan.",
							Computed:            true,
						},
						"platform_id": schema.StringAttribute{
							MarkdownDescription: "The platform ID.",
							Computed:            true,
						},
						"context": schema.StringAttribute{
							MarkdownDescription: "Contextual data for the resource.",
							Computed:            true,
						},
						"usable": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the resource can be used.",
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
							MarkdownDescription: "The set of words or phrases assigned to the service instance.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountServiceInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServiceInstancesDataSourceConfig

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

	cliRes, _, err := ds.cli.Services.Instance.List(ctx, data.SubaccountId.ValueString(), fieldsFilter, labelsFilter)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Instances (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.SubaccountId
	data.Values = []subaccountServiceInstancesValueConfig{}
	for _, serviceInstance := range cliRes {
		val := subaccountServiceInstancesValueConfig{
			Id:            types.StringValue(serviceInstance.Id),
			Name:          types.StringValue(serviceInstance.Name),
			Ready:         types.BoolValue(serviceInstance.Ready),
			ServicePlanId: types.StringValue(serviceInstance.ServicePlanId),
			PlatformId:    types.StringValue(serviceInstance.PlatformId),
			Usable:        types.BoolValue(serviceInstance.Usable),
			Context:       types.StringValue(string(serviceInstance.Context)),
			CreatedDate:   timeToValue(serviceInstance.CreatedAt),
			LastModified:  timeToValue(serviceInstance.UpdatedAt),
		}

		val.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, serviceInstance.Labels)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, val)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
