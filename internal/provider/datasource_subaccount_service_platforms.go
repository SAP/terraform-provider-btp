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

func newSubaccountServicePlatformsDataSource() datasource.DataSource {
	return &subaccountServicePlatformsDataSource{}
}

type subaccountServicePlatformsValue struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Ready        types.Bool   `tfsdk:"ready"`
	PlatformType types.String `tfsdk:"type"`
	Description  types.String `tfsdk:"description"`
	CreatedDate  types.String `tfsdk:"created_date"`
	LastModified types.String `tfsdk:"last_modified"`
	Labels       types.Map    `tfsdk:"labels"`
}

type subaccountServicePlatformsDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	FieldsFilter types.String `tfsdk:"fields_filter"`
	LabelsFilter types.String `tfsdk:"labels_filter"`
	/* OUTPUT */
	Values []subaccountServicePlatformsValue `tfsdk:"values"`
}

type subaccountServicePlatformsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServicePlatformsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_platforms", req.ProviderTypeName)
}

func (ds *subaccountServicePlatformsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServicePlatformsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Lists all platforms in a subaccount that are registered for service consumption.

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
			"fields_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the platforms based on their fields. For example, to display all 'kubernetes' platforms, use \"type eq 'kubernetes'\".",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"labels_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the platforms based on the label query. For example, to list all platforms whose purpose is 'dev', use \"purpose eq 'dev'\".",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the platform.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the platform.",
							Optional:            true,
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the platform.",
							Computed:            true,
						},
						"ready": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the platform is ready for consumption.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type of the platform.",
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
							MarkdownDescription: "Set of words or phrases assigned to the platform.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountServicePlatformsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServicePlatformsDataSourceConfig

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

	cliRes, _, err := ds.cli.Services.Platform.List(ctx, data.SubaccountId.ValueString(), fieldsFilter, labelsFilter)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Platforms (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Values = []subaccountServicePlatformsValue{}

	for _, platform := range cliRes {
		platformValue := subaccountServicePlatformsValue{
			Id:           types.StringValue(platform.Id),
			Ready:        types.BoolValue(platform.Ready),
			PlatformType: types.StringValue(platform.Type_),
			Name:         types.StringValue(platform.Name),
			Description:  types.StringValue(platform.Description),
			CreatedDate:  timeToValue(platform.CreatedAt),
			LastModified: timeToValue(platform.UpdatedAt),
		}

		platformValue.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, platform.Labels)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, platformValue)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
