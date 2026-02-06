package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountEntitlementWithDcDataSource() datasource.DataSource {
	return &globalaccountEntitlementWithDcDataSource{}
}

type globalaccountEntitlementWithDcDataSource struct {
	cli *btpcli.ClientFacade
}

type globalaccountEntitlementWithDcDataSourceModel struct {
	ServiceName          types.String  `tfsdk:"service_name"`
	ServiceDisplayName   types.String  `tfsdk:"service_display_name"`
	PlanName             types.String  `tfsdk:"plan_name"`
	PlanDisplayName      types.String  `tfsdk:"plan_display_name"`
	PlanDescription      types.String  `tfsdk:"plan_description"`
	PlanUniqueIdentifier types.String  `tfsdk:"plan_unique_identifier"`
	QuotaAssigned        types.Float64 `tfsdk:"quota_assigned"`
	QuotaRemaining       types.Float64 `tfsdk:"quota_remaining"`
	Category             types.String  `tfsdk:"category"`
	DataCenterInfo       types.Map     `tfsdk:"datacenter_information"`
}

func (ds *globalaccountEntitlementWithDcDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_entitlement_with_data_centers", req.ProviderTypeName)
}

func (ds *globalaccountEntitlementWithDcDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		ds.cli = req.ProviderData.(*btpcli.ClientFacade)
	}
}

func (ds *globalaccountEntitlementWithDcDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Returns entitlement details for a specific plan of a given service for a global account including the data center availability.

__Tip:__
You must be assigned to either the global account admin or global account viewers role.`,
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				MarkdownDescription: "The name of the entitled service.",
				Required:            true,
			},
			"plan_name": schema.StringAttribute{
				MarkdownDescription: "The name of the entitled service plan.",
				Required:            true,
			},
			"service_display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the entitled service.",
				Computed:            true,
			},
			"plan_display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the entitled service plan.",
				Computed:            true,
			},
			"plan_description": schema.StringAttribute{
				MarkdownDescription: "The description of the entitled service plan.",
				Computed:            true,
			},
			"plan_unique_identifier": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the entitled service plan.",
				Computed:            true,
			},
			"quota_assigned": schema.Float64Attribute{
				MarkdownDescription: "The overall quota assigned.",
				Computed:            true,
			},
			"quota_remaining": schema.Float64Attribute{
				MarkdownDescription: "The quota, which is not used.",
				Computed:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "The current state of the entitlement. Possible values are: \n " +
					getFormattedValueAsTableRow("value", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`PLATFORM`", " A service required for using a specific platform; for example, Application Runtime is required for the Cloud Foundry platform.") +
					getFormattedValueAsTableRow("`SERVICE`", "A commercial or technical service. that has a numeric quota (amount) when entitled or assigned to a resource. When assigning entitlements of this type, use the 'amount' option.") +
					getFormattedValueAsTableRow("`ELASTIC_SERVICE`", "A commercial or technical service that has no numeric quota (amount) when entitled or assigned to a resource. Generally this type of service can be as many times as needed when enabled, but may in some cases be restricted by the service owner.") +
					getFormattedValueAsTableRow("`ELASTIC_LIMITED`", "An elastic service that can be enabled for only one subaccount per global account.") +
					getFormattedValueAsTableRow("`APPLICATION`", "A multitenant application to which consumers can subscribe. As opposed to applications defined as a 'QUOTA_BASED_APPLICATION', these applications do not have a numeric quota and are simply enabled or disabled as entitlements per subaccount.") +
					getFormattedValueAsTableRow("`QUOTA_BASED_APPLICATION`", "A multitenant application to which consumers can subscribe. As opposed to applications defined as 'APPLICATION', these applications have an numeric quota that limits consumer usage of the subscribed application per subaccount.") +
					getFormattedValueAsTableRow("`ENVIRONMENT`", " An environment service; for example, Cloud Foundry."),
				Computed: true,
			},
			"datacenter_information": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"dc_region": schema.StringAttribute{
							MarkdownDescription: "The region in which the data center is located.",
							Computed:            true,
						},
						"dc_name": schema.StringAttribute{
							MarkdownDescription: "The technical name of the data center.",
							Computed:            true,
						},
						"dc_display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the data center.",
							Computed:            true,
						},
						"dc_iaas_provider": schema.StringAttribute{
							MarkdownDescription: "The infrastructure provider for the data center.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *globalaccountEntitlementWithDcDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Note: plan_name must be specified. Only exact matches will be returned.
	var data globalaccountEntitlementWithDcDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	planName := data.PlanName.ValueString()

	cliRes, _, err := ds.cli.Accounts.Entitlement.ListByGlobalAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Entitlements (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	found := false
	for _, service := range cliRes.EntitledServices {
		if service.Name != serviceName {
			continue
		}

		for _, servicePlan := range service.ServicePlans {
			if servicePlan.Name != planName {
				continue
			}

			// Build datacenter information map
			dataCenterMap := make(map[string]datacenterInformation)
			for _, dc := range servicePlan.DataCenters {
				dataCenterMap[dc.Region] = datacenterInformation{
					Region:       types.StringValue(dc.Region),
					Name:         types.StringValue(dc.Name),
					DisplayName:  types.StringValue(dc.DisplayName),
					IaasProvider: types.StringValue(dc.IaasProvider),
				}
			}

			// Convert datacenter map to types.Map
			dcMapValue, dcDiags := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: datacenterInformationType()}, dataCenterMap)
			resp.Diagnostics.Append(dcDiags...)
			if resp.Diagnostics.HasError() {
				return
			}

			// Populate the data model with matching entry
			data.ServiceDisplayName = types.StringValue(service.DisplayName)
			data.PlanDisplayName = types.StringValue(servicePlan.DisplayName)
			data.PlanDescription = types.StringValue(servicePlan.Description)
			data.PlanUniqueIdentifier = types.StringValue(servicePlan.UniqueIdentifier)
			data.QuotaAssigned = types.Float64Value(servicePlan.Amount)
			data.QuotaRemaining = types.Float64Value(servicePlan.RemainingAmount)
			data.Category = types.StringValue(servicePlan.Category)
			data.DataCenterInfo = dcMapValue

			found = true
			break
		}

		if found {
			break
		}
	}

	// If no matching entry was found, set computed fields to null/empty
	if !found {
		data.ServiceDisplayName = types.StringNull()
		data.PlanDisplayName = types.StringNull()
		data.PlanDescription = types.StringNull()
		data.PlanUniqueIdentifier = types.StringNull()
		data.QuotaAssigned = types.Float64Null()
		data.QuotaRemaining = types.Float64Null()
		data.Category = types.StringNull()
		data.DataCenterInfo = types.MapNull(types.ObjectType{AttrTypes: datacenterInformationType()})
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
