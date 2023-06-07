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

func newSubaccountEnvironmentInstancesDataSource() datasource.DataSource {
	return &subaccountEnvironmentInstancesDataSource{}
}

type subaccountEnvironmentInstanceValue struct {
	Id              types.String `tfsdk:"id"`
	BrokerId        types.String `tfsdk:"broker_id"`
	CreatedDate     types.String `tfsdk:"created_date"`
	CustomLabels    types.Map    `tfsdk:"custom_labels"`
	DashboardUrl    types.String `tfsdk:"dashboard_url"`
	Description     types.String `tfsdk:"description"`
	EnvironmentType types.String `tfsdk:"environment_type"`
	Labels          types.String `tfsdk:"labels"`
	LandscapeLabel  types.String `tfsdk:"landscape_label"`
	LastModified    types.String `tfsdk:"last_modified"`
	Name            types.String `tfsdk:"name"`
	Operation       types.String `tfsdk:"operation"`
	Parameters      types.String `tfsdk:"parameters"`
	PlanId          types.String `tfsdk:"plan_id"`
	PlanName        types.String `tfsdk:"plan_name"`
	PlatformId      types.String `tfsdk:"platform_id"`
	ServiceId       types.String `tfsdk:"service_id"`
	ServiceName     types.String `tfsdk:"service_name"`
	State           types.String `tfsdk:"state"`
	TenantId        types.String `tfsdk:"tenant_id"`
	Type_           types.String `tfsdk:"type"`
}

type subaccountEnvironmentInstancesDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	/* OUTPUT */
	Values []subaccountEnvironmentInstanceValue `tfsdk:"values"`
}

type subaccountEnvironmentInstancesDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountEnvironmentInstancesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_environment_instances", req.ProviderTypeName)
}

func (ds *subaccountEnvironmentInstancesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountEnvironmentInstancesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `List all the environment instances in a subaccount.

__Tip__
You must be assigned to the subaccount admin or viewer role.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the environment instance.",
							Computed:            true,
						},
						"broker_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the associated environment broker.",
							Computed:            true,
						},
						"created_date": schema.StringAttribute{
							MarkdownDescription: "The date and time the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
							Computed:            true,
						},
						"custom_labels": schema.MapAttribute{
							ElementType: types.SetType{
								ElemType: types.StringType,
							},
							Computed: true,
						},
						"dashboard_url": schema.StringAttribute{
							MarkdownDescription: "The URL of the service dashboard, which is a web-based management user interface for the service instances.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the environment instance.",
							Computed:            true,
						},
						"environment_type": schema.StringAttribute{
							MarkdownDescription: "Type of the environment instance that is used.",
							Computed:            true,
						},
						"labels": schema.StringAttribute{
							MarkdownDescription: "Broker-specified key-value pairs that specify attributes of an environment instance.",
							Computed:            true,
						},
						"landscape_label": schema.StringAttribute{
							MarkdownDescription: "The name of the landscape within the logged-in region on which the environment instance is created.",
							Computed:            true,
						},
						"last_modified": schema.StringAttribute{
							MarkdownDescription: "The date and time the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the environment instance.",
							Computed:            true,
						},
						"operation": schema.StringAttribute{
							MarkdownDescription: "An identifier that represents the last operation. This ID is returned by the environment brokers.",
							Computed:            true,
						},
						"parameters": schema.StringAttribute{
							MarkdownDescription: "Configuration parameters for the environment instance.",
							Computed:            true,
						},
						"plan_id": schema.StringAttribute{
							MarkdownDescription: "ID of the service plan for the environment instance in the corresponding service broker's catalog.",
							Computed:            true,
						},
						"plan_name": schema.StringAttribute{
							MarkdownDescription: "Name of the service plan for the environment instance in the corresponding service broker's catalog.",
							Computed:            true,
						},
						"platform_id": schema.StringAttribute{
							MarkdownDescription: "ID of the platform for the environment instance in the corresponding service broker's catalog.",
							Computed:            true,
						},
						"service_id": schema.StringAttribute{
							MarkdownDescription: "ID of the service for the environment instance in the corresponding service broker's catalog.",
							Computed:            true,
						},
						"service_name": schema.StringAttribute{
							MarkdownDescription: "Name of the service for the environment instance in the corresponding service broker's catalog.",
							Computed:            true,
						},
						"state": schema.StringAttribute{
							MarkdownDescription: "Current state of the environment instance. Possible values are: " + // TODO describe states listed below
								"\n\t - `OK`" +
								"\n\t - `CREATING`" +
								"\n\t - `CREATION_FAILED`" +
								"\n\t - `DELETING`" +
								"\n\t - `DELETION_FAILED`" +
								"\n\t - `UPDATE_FAILED`" +
								"\n\t - `UPDATING`",
							Computed: true,
						},
						"tenant_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the tenant that owns the environment instance.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The last provisioning operation on the environment instance. Possible values are: " +
								"\n\t - `Provision` Environment instance created." +
								"\n\t - `Update` Environment instance changed." +
								"\n\t - `Deprovision` Environment instance deleted.",
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountEnvironmentInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountEnvironmentInstancesDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.EnvironmentInstance.List(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Environment Instances (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Values = []subaccountEnvironmentInstanceValue{}

	for _, instance := range cliRes.EnvironmentInstances {
		instanceValue := subaccountEnvironmentInstanceValue{
			Id:              types.StringValue(instance.Id),
			BrokerId:        types.StringValue(instance.BrokerId),
			CreatedDate:     timeToValue(instance.CreatedDate.Time()),
			DashboardUrl:    types.StringValue(instance.DashboardUrl),
			Description:     types.StringValue(instance.Description),
			EnvironmentType: types.StringValue(instance.EnvironmentType),
			Labels:          types.StringValue(instance.Labels),
			LandscapeLabel:  types.StringValue(instance.LandscapeLabel),
			LastModified:    timeToValue(instance.ModifiedDate.Time()),
			Name:            types.StringValue(instance.Name),
			Operation:       types.StringValue(instance.Operation),
			Parameters:      types.StringValue(instance.Parameters),
			PlanId:          types.StringValue(instance.PlanId),
			PlanName:        types.StringValue(instance.PlanName),
			PlatformId:      types.StringValue(instance.PlatformId),
			ServiceId:       types.StringValue(instance.ServiceId),
			ServiceName:     types.StringValue(instance.ServiceName),
			State:           types.StringValue(instance.State),
			TenantId:        types.StringValue(instance.TenantId),
			Type_:           types.StringValue(instance.Type_),
		}

		instanceValue.CustomLabels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, instance.CustomLabels)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, instanceValue)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
