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
}

func entitledServiceType() map[string]attr.Type {
	return map[string]attr.Type{
		"service_name":         types.StringType,
		"service_display_name": types.StringType,
		"plan_name":            types.StringType,
		"plan_display_name":    types.StringType,
		"plan_description":     types.StringType,
		"quota_assigned":       types.NumberType,
		"quota_remaining":      types.NumberType,
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
		MarkdownDescription: `Get all the entitlements and quota assignments for a global account.

To view all the resources a global account:
* Target only the global account in the command line.
* You must be assigned to either the global account admin or global account viewers role.

__Tips__
You must be assigned to one of these roles: global account admin, global account viewer.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
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
						"quota_assigned": schema.NumberAttribute{
							MarkdownDescription: "The overall quota assigned.",
							Computed:            true,
						},
						"quota_remaining": schema.NumberAttribute{
							MarkdownDescription: "The quota which is not used.",
							Computed:            true,
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
			}
		}
	}

	data.Id = types.StringValue(ds.cli.GetGlobalAccountSubdomain())

	data.Values, diags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: entitledServiceType()}, values)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
