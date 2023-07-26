package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountEntitlementsDataSource() datasource.DataSource {
	return &globalaccountEntitlementsDataSource{}
}

type entitledService struct {
	ServiceName        types.String  `tfsdk:"service_name"`
	ServiceDisplayName types.String  `tfsdk:"service_display_name"`
	PlanName           types.String  `tfsdk:"plan_name"`
	PlanDisplayName    types.String  `tfsdk:"plan_display_name"`
	PlanDescription    types.String  `tfsdk:"plan_description"`
	QuotaAssigned      types.Float64 `tfsdk:"quota_assigned"`
	QuotaRemaining     types.Float64 `tfsdk:"quota_remaining"`
	Category           types.String  `tfsdk:"category"`
}

func entitledServiceType() map[string]attr.Type {
	return map[string]attr.Type{
		"service_name":         types.StringType,
		"service_display_name": types.StringType,
		"plan_name":            types.StringType,
		"plan_display_name":    types.StringType,
		"plan_description":     types.StringType,
		"quota_assigned":       types.Float64Type,
		"quota_remaining":      types.Float64Type,
		"category":             types.StringType,
	}
}

type globalaccountEntitlementsDataSourceConfig struct {
	/* INPUT */
	Id types.String `tfsdk:"id"`
	/* OUTPUT */
	Values types.Map `tfsdk:"values"`
}

type globalaccountEntitlementsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountEntitlementsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_entitlements", req.ProviderTypeName)
}

func (ds *globalaccountEntitlementsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountEntitlementsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets all the entitlements and quota assignments for a global account.

To view all the resources a global account:
* Target only the global account in the command line.
* You must be assigned to either the global account admin or global account viewers role.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `btp_globalaccount` datasource instead",
				MarkdownDescription: "The ID of the global account.",
				Computed:            true,
			},
			"values": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"service_name": schema.StringAttribute{
							MarkdownDescription: "The name of the entitled service.",
							Computed:            true,
						},
						"service_display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the entitled service.",
							Computed:            true,
						},
						"plan_name": schema.StringAttribute{
							MarkdownDescription: "The name of the entitled service plan.",
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
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *globalaccountEntitlementsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountEntitlementsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.Entitlement.ListByGlobalAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Entitlements (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	values := map[string]entitledService{}

	for _, service := range cliRes.EntitledServices {
		for _, servicePlan := range service.ServicePlans {
			values[fmt.Sprintf("%s:%s", service.Name, servicePlan.Name)] = entitledService{
				ServiceName:        types.StringValue(service.Name),
				ServiceDisplayName: types.StringValue(service.DisplayName),
				PlanName:           types.StringValue(servicePlan.Name),
				PlanDisplayName:    types.StringValue(servicePlan.DisplayName),
				PlanDescription:    types.StringValue(servicePlan.Description),
				QuotaAssigned:      types.Float64Value(servicePlan.Amount),
				QuotaRemaining:     types.Float64Value(servicePlan.RemainingAmount),
				Category:           types.StringValue(servicePlan.Category),
			}
		}
	}

	data.Id = types.StringValue(ds.cli.GetGlobalAccountSubdomain())

	data.Values, diags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: entitledServiceType()}, values)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
