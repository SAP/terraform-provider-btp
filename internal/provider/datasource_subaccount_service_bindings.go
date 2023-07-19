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

func newSubaccountServiceBindingsDataSource() datasource.DataSource {
	return &subaccountServiceBindingsDataSource{}
}

type subaccountServiceBindingValue struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Ready             types.Bool   `tfsdk:"ready"`
	ServiceInstanceId types.String `tfsdk:"service_instance_id"`
	Context           types.Map    `tfsdk:"context"`
	BindResource      types.Map    `tfsdk:"bind_resource"`
	Credentials       types.String `tfsdk:"credentials"`
	CreatedDate       types.String `tfsdk:"created_date"`
	LastModified      types.String `tfsdk:"last_modified"`
	Labels            types.Map    `tfsdk:"labels"`
}

type subaccountServiceBindingsDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	FieldsFilter types.String `tfsdk:"fields_filter"`
	LabelsFilter types.String `tfsdk:"labels_filter"`
	/* OUTPUT */
	Values []subaccountServiceBindingValue `tfsdk:"values"`
}

type subaccountServiceBindingsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServiceBindingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_bindings", req.ProviderTypeName)
}

func (ds *subaccountServiceBindingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServiceBindingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists all service bindings in a subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` attribute instead",
				MarkdownDescription: "The ID of the subaccount.",
				Computed:            true,
			},
			"fields_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the service bindings based on the field query.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"labels_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the service binding based on the label query.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service binding.",
							Optional:            true,
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the service binding.",
							Optional:            true,
							Computed:            true,
						},
						"ready": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the service binding is ready.",
							Computed:            true,
						},
						"service_instance_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the service instance associated with the binding.",
							Computed:            true,
						},
						"context": schema.MapAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "Contextual data for the resource.",
							Computed:            true,
						},
						"bind_resource": schema.MapAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "Contains the resources associated with the binding.",
							Computed:            true,
						},
						"credentials": schema.StringAttribute{
							MarkdownDescription: "The credentials to access the binding.",
							Computed:            true,
							Sensitive:           true,
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
							MarkdownDescription: "Set of words or phrases assigned to the binding.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountServiceBindingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServiceBindingsDataSourceConfig

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

	cliRes, _, err := ds.cli.Services.Binding.List(ctx, data.SubaccountId.ValueString(), fieldsFilter, labelsFilter)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Bindings (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.SubaccountId
	data.Values = []subaccountServiceBindingValue{}

	for _, binding := range cliRes {
		bindingValue := subaccountServiceBindingValue{
			Id:                types.StringValue(binding.Id),
			Name:              types.StringValue(binding.Name),
			Ready:             types.BoolValue(binding.Ready),
			ServiceInstanceId: types.StringValue(binding.ServiceInstanceId),
			Credentials:       types.StringValue(string(binding.Credentials)),
			CreatedDate:       timeToValue(binding.CreatedAt),
			LastModified:      timeToValue(binding.UpdatedAt),
		}
		bindingValue.Context, diags = types.MapValueFrom(ctx, types.StringType, binding.Context)
		resp.Diagnostics.Append(diags...)

		bindingValue.BindResource, diags = types.MapValueFrom(ctx, types.StringType, binding.BindResource)
		resp.Diagnostics.Append(diags...)

		bindingValue.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, binding.Labels)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, bindingValue)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
