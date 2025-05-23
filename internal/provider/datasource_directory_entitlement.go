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

func newDirectoryEntitlementDataSource() datasource.DataSource {
	return &directoryEntitlementDataSource{}
}

type directoryEntitlementDataSource struct {
	cli *btpcli.ClientFacade
}

type directoryEntitlementDataSourceModel struct {
	DirectoryId          types.String  `tfsdk:"directory_id"`
	ServiceName          types.String  `tfsdk:"service_name"`
	PlanName             types.String  `tfsdk:"plan_name"`
	PlanUniqueIdentifier types.String  `tfsdk:"plan_unique_identifier"`
	PlanId               types.String  `tfsdk:"plan_id"`
	QuotaAssigned        types.Float64 `tfsdk:"quota_assigned"`
	QuotaRemaining       types.Float64 `tfsdk:"quota_remaining"`
	Category             types.String  `tfsdk:"category"`
	Id                   types.String  `tfsdk:"id"`
}

func (ds *directoryEntitlementDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_entitlement", req.ProviderTypeName)
}

func (ds *directoryEntitlementDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		ds.cli = req.ProviderData.(*btpcli.ClientFacade)
	}
}

func (ds *directoryEntitlementDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Gets a specific entitlement assigned to a directory.",
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the directory.",
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"service_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the entitled service.",
			},
			"plan_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the entitled service plan.",
			},
			"plan_unique_identifier": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the entitled service plan.",
			},
			"plan_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the entitled service plan.",
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
				Computed:            true,
				MarkdownDescription: "The category of the entitlement.",
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Synthetic ID combining directory ID, service name, and plan name.",
			},
		},
	}
}

func (ds *directoryEntitlementDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryEntitlementDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entitlement, _, err := ds.cli.Accounts.Entitlement.GetEntitledByDirectory(
		ctx,
		data.DirectoryId.ValueString(),
		data.ServiceName.ValueString(),
		data.PlanName.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Directory Entitlement", err.Error())
		return
	}
	if entitlement == nil {
		resp.Diagnostics.AddError("Directory Entitlement Not Found", "The specified entitlement could not be found.")
		return
	}

	data.PlanUniqueIdentifier = types.StringValue(entitlement.Plan.UniqueIdentifier)
	data.PlanId = types.StringValue(entitlement.Plan.UniqueIdentifier)
	data.QuotaAssigned = types.Float64Value(entitlement.Plan.Amount)
	data.QuotaRemaining = types.Float64Value(entitlement.Plan.RemainingAmount)
	data.Category = types.StringValue(entitlement.Plan.Category)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
