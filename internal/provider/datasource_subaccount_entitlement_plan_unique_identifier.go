package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis_entitlements"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func NewSubaccountEntitlementUniqueIdentifierDataSource() datasource.DataSource {
	return &subaccountEntitlementUniqueIdentifierDataSource{}
}

type subaccountEntitlementUniqueIdentifierDataSource struct {
	cli *btpcli.ClientFacade
}

type subaccountEntitlementUniqueIdentifierDataSourceModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	ServiceName  types.String `tfsdk:"service_name"`
	PlanName     types.String `tfsdk:"plan_name"`
	Entitlements types.List   `tfsdk:"entitlements"`
	Id           types.String `tfsdk:"id"`
}

func (ds *subaccountEntitlementUniqueIdentifierDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_entitlement_unique_identifier", req.ProviderTypeName)
}

func (ds *subaccountEntitlementUniqueIdentifierDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		ds.cli = req.ProviderData.(*btpcli.ClientFacade)
	}
}

func (ds *subaccountEntitlementUniqueIdentifierDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Returns the `plan_unique_identifier` for a given service and plan name from a subaccount's entitlements.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the subaccount.",
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"service_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The technical name of the service.",
			},
			"plan_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the service plan.",
			},
			"entitlements": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of entitlements that match the service and plan names.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"service_name": schema.StringAttribute{
							Computed: true,
						},
						"plan_name": schema.StringAttribute{
							Computed: true,
						},
						"plan_unique_identifier": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal Terraform tracking ID.",
			},
		},
	}
}

func (ds *subaccountEntitlementUniqueIdentifierDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config subaccountEntitlementUniqueIdentifierDataSourceModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	subaccountID := config.SubaccountId.ValueString()
	serviceName := config.ServiceName.ValueString()
	planName := config.PlanName.ValueString()

	subaccountData, _, err := ds.cli.Accounts.Subaccount.Get(ctx, subaccountID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get subaccount", err.Error())
		return
	}

	parentId, isParentGlobalAccount := determineParentIdForEntitlement(ds.cli, ctx, subaccountData.ParentGUID)

	var cliRes cis_entitlements.EntitledAndAssignedServicesResponseObject
	if isParentGlobalAccount {
		cliRes, _, err = ds.cli.Accounts.Entitlement.FilterBySubaccount(ctx, subaccountID)
	} else {
		cliRes, _, err = ds.cli.Accounts.Entitlement.ListBySubaccountWithDirectoryParent(ctx, subaccountID, parentId)
	}
	if err != nil {
		resp.Diagnostics.AddError("Failed to list entitlements", err.Error())
		return
	}

	var entitlementsList []attr.Value

	for _, service := range cliRes.EntitledServices {
		if service.Name == serviceName {
			for _, plan := range service.ServicePlans {
				if plan.Name == planName {
					entitlement, diag := types.ObjectValue(map[string]attr.Type{
						"service_name":           types.StringType,
						"plan_name":              types.StringType,
						"plan_unique_identifier": types.StringType,
					}, map[string]attr.Value{
						"service_name":           types.StringValue(service.Name),
						"plan_name":              types.StringValue(plan.Name),
						"plan_unique_identifier": types.StringValue(plan.UniqueIdentifier),
					})
					if diag.HasError() {
						resp.Diagnostics.Append(diag...)
						return
					}
					entitlementsList = append(entitlementsList, entitlement)
				}
			}
		}
	}

	if len(entitlementsList) == 0 {
		resp.Diagnostics.AddError("No Entitlements Found", fmt.Sprintf("No entitlements found for service '%s' with plan '%s' in subaccount '%s'", serviceName, planName, subaccountID))
		return
	}

	entitlements, diag := types.ListValue(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"service_name":           types.StringType,
			"plan_name":              types.StringType,
			"plan_unique_identifier": types.StringType,
		},
	}, entitlementsList)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}

	// Set final state
	config.Entitlements = entitlements
	config.Id = types.StringValue(fmt.Sprintf("%s:%s:%s", subaccountID, serviceName, planName))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
