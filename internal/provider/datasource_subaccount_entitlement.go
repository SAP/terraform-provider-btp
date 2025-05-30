package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis_entitlements"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountEntitlementDataSource() datasource.DataSource {
	return &subaccountEntitlementDataSource{}
}

type subaccountEntitlementDataSource struct {
	cli *btpcli.ClientFacade
}

type subaccountEntitlementDataSourceModel struct {
	SubaccountId         types.String  `tfsdk:"subaccount_id"`
	ServiceName          types.String  `tfsdk:"service_name"`
	PlanName             types.String  `tfsdk:"plan_name"`
	PlanUniqueIdentifier types.String  `tfsdk:"plan_unique_identifier"`
	PlanId               types.String  `tfsdk:"plan_id"`
	QuotaAssigned        types.Float64 `tfsdk:"quota_assigned"`
	QuotaRemaining       types.Float64 `tfsdk:"quota_remaining"`
	Category             types.String  `tfsdk:"category"`
	Id                   types.String  `tfsdk:"id"`
}

func (ds *subaccountEntitlementDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_entitlement", req.ProviderTypeName)
}

func (ds *subaccountEntitlementDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		ds.cli = req.ProviderData.(*btpcli.ClientFacade)
	}
}

func (ds *subaccountEntitlementDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Returns entitlement details for a specific plan of a given service in a subaccount.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount. Additionally you must be a Global Account Administrator or Viewer.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"service_name": schema.StringAttribute{
				MarkdownDescription: "The name of the entitled service.",
				Required:            true,
			},
			"plan_name": schema.StringAttribute{
				MarkdownDescription: "The name of the entitled service plan.",
				Required:            true,
			},
			"plan_id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of the entitled service plan. Alias for plan_unique_identifier.",
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
			"id": schema.StringAttribute{
				MarkdownDescription: "Synthetic ID combining subaccount ID, service name, and plan name.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountEntitlementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Note: plan_name must be specified. Only exact matches will be returned.
	var data subaccountEntitlementDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	subaccountId := data.SubaccountId.ValueString()
	serviceName := data.ServiceName.ValueString()
	planFilter := data.PlanName.ValueString()

	subaccountData, _, err := ds.cli.Accounts.Subaccount.Get(ctx, subaccountId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get subaccount", err.Error())
		return
	}

	parentId, isParentGlobalAccount, err := determineParentIdForEntitlement(ds.cli, ctx, subaccountData.ParentGUID)
	if err != nil {
		resp.Diagnostics.AddError("API Error determining parent features for entitlement", fmt.Sprintf("%s", err))
		return
	}

	var cliRes cis_entitlements.EntitledAndAssignedServicesResponseObject
	if isParentGlobalAccount {
		cliRes, _, err = ds.cli.Accounts.Entitlement.ListBySubaccount(ctx, subaccountId)
	} else {
		cliRes, _, err = ds.cli.Accounts.Entitlement.ListBySubaccountWithDirectoryParent(ctx, subaccountId, parentId)
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to list entitlements", err.Error())
		return
	}

	for _, service := range cliRes.EntitledServices {
		if service.Name != serviceName {
			continue
		}

		for _, plan := range service.ServicePlans {
			if planFilter != "" && plan.Name != planFilter {
				continue
			}

			data.PlanName = types.StringValue(plan.Name)
			data.PlanId = types.StringValue(plan.UniqueIdentifier)
			data.PlanUniqueIdentifier = types.StringValue(plan.UniqueIdentifier)
			data.QuotaAssigned = types.Float64Value(plan.Amount)
			data.QuotaRemaining = types.Float64Value(plan.RemainingAmount)
			data.Category = types.StringValue(plan.Category)
			data.Id = types.StringValue(fmt.Sprintf("%s:%s:%s", subaccountId, serviceName, plan.Name))
			diags = resp.State.Set(ctx, &data)
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	resp.Diagnostics.AddError("Plan Not Found", fmt.Sprintf(
		"No plan found for service '%s' with plan_name = '%s' in subaccount '%s'.",
		serviceName, planFilter, subaccountId,
	))
}
