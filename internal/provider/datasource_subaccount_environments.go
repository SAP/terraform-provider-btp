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

func newSubaccountEnvironmentsDataSource() datasource.DataSource {
	return &subaccountEnvironmentsDataSource{}
}

type subaccountEnvironment struct {
	AvailabilityLevel  types.String `tfsdk:"availability_level"`
	CreateSchema       types.String `tfsdk:"schema_create"`
	Description        types.String `tfsdk:"description"`
	EnvironmentType    types.String `tfsdk:"environment_type"`
	LandscapeLabel     types.String `tfsdk:"landscape_label"`
	PlanName           types.String `tfsdk:"plan_name"`
	PlanUpdatable      types.Bool   `tfsdk:"plan_updateable"`
	ServiceDescription types.String `tfsdk:"service_description"`
	ServiceDisplayName types.String `tfsdk:"service_display_name"`
	ServiceName        types.String `tfsdk:"service_name"`
	TechnicalKey       types.String `tfsdk:"technical_key"`
	UpdateSchema       types.String `tfsdk:"schema_update"`
}

type subaccountEnvironmentsDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	/* OUTPUT */
	Values []subaccountEnvironment `tfsdk:"values"`
}

type subaccountEnvironmentsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountEnvironmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_environments", req.ProviderTypeName)
}

func (ds *subaccountEnvironmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountEnvironmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `List all the available environments for a subaccount.

This includes the environments, such as Cloud Foundry, which are available by default to all subaccounts, and those restricted environments, such as Kyma, which are offered in the product catalog as service entitlements and whose plans have already been assigned by a global account admin to the subaccount.

__Tips__
You must be assigned to the subaccount admin or viewer role.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{
				DeprecationMessage:  "Use the `subaccount_id` attribute instead",
				MarkdownDescription: "The ID of the subaccount.",
				Computed:            true,
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"availability_level": schema.StringAttribute{
							MarkdownDescription: "The availability level of the environment broker.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the environment.",
							Computed:            true,
						},
						"environment_type": schema.StringAttribute{
							MarkdownDescription: "The type of environment that is available (for example: cloudfoundry).",
							Computed:            true,
						},
						"landscape_label": schema.StringAttribute{
							MarkdownDescription: "The landscape label of the environment broker.",
							Computed:            true,
						},
						"plan_name": schema.StringAttribute{
							MarkdownDescription: "Name of the service plan for the available environment.",
							Computed:            true,
						},
						"plan_updateable": schema.BoolAttribute{
							MarkdownDescription: "Specifies if the consumer can change the plan of an existing instance of the environment.",
							Computed:            true,
						},
						"service_description": schema.StringAttribute{
							MarkdownDescription: "The short description of the service.",
							Computed:            true,
						},
						"service_display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the service.",
							Computed:            true,
						},
						"service_name": schema.StringAttribute{
							MarkdownDescription: "Name of the service offered in the catalog of the corresponding environment broker (for example: cloudfoundry).",
							Computed:            true,
						},
						"schema_create": schema.StringAttribute{
							MarkdownDescription: "The create schema of the environment broker.",
							Computed:            true,
						},
						"schema_update": schema.StringAttribute{
							MarkdownDescription: "The update schema of the environment broker.",
							Computed:            true,
						},
						"technical_key": schema.StringAttribute{
							MarkdownDescription: "Technical key of the corresponding environment broker.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountEnvironmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountEnvironmentsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.AvailableEnvironment.List(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Environments (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Values = []subaccountEnvironment{}

	for _, availableEnvironment := range cliRes.AvailableEnvironments {
		data.Values = append(data.Values, subaccountEnvironment{
			AvailabilityLevel:  types.StringValue(availableEnvironment.AvailabilityLevel),
			CreateSchema:       types.StringValue(availableEnvironment.CreateSchema),
			Description:        types.StringValue(availableEnvironment.Description),
			EnvironmentType:    types.StringValue(availableEnvironment.EnvironmentType),
			LandscapeLabel:     types.StringValue(availableEnvironment.LandscapeLabel),
			PlanName:           types.StringValue(availableEnvironment.PlanName),
			PlanUpdatable:      types.BoolValue(availableEnvironment.PlanUpdatable),
			ServiceDescription: types.StringValue(availableEnvironment.ServiceDescription),
			ServiceDisplayName: types.StringValue(availableEnvironment.ServiceDisplayName),
			ServiceName:        types.StringValue(availableEnvironment.ServiceName),
			TechnicalKey:       types.StringValue(availableEnvironment.TechnicalKey),
			UpdateSchema:       types.StringValue(availableEnvironment.UpdateSchema),
		})
	}

	data.Id = data.SubaccountId

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
