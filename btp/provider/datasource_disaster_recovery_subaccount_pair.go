package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newDisasterRecoverySubaccountPairDataSource() datasource.DataSource {
	return &DisasterRecoverySubaccountPairDataSource{}
}

type DisasterRecoverySubaccountPairDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *DisasterRecoverySubaccountPairDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_disaster_recovery_subaccount_pair", req.ProviderTypeName)
}

func (ds *DisasterRecoverySubaccountPairDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *DisasterRecoverySubaccountPairDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets the subaccount pair details for the specified subaccount.

__Tip:__
You must be assigned to the central .`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of a paired subaccount.",
				Required:            true,
			},
			"pair_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount pair.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the subaccount pair was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The user who created the subaccount pair.",
				Computed:            true,
			},
			"global_account_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the global account.",
				Computed:            true,
			},

			"subaccounts": schema.ListNestedAttribute{
				MarkdownDescription: "The list of subaccounts in the disaster recovery pair.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The ID of the subaccount.",
							Computed:            true,
						},
						"region": schema.StringAttribute{
							MarkdownDescription: "The region of the subaccount.",
							Computed:            true,
						},
						"subdomain": schema.StringAttribute{
							MarkdownDescription: "The subdomain of the subaccount.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (ds *DisasterRecoverySubaccountPairDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DisasterRecoverySubaccountPairDataSourceType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	subaccountPair, _, err := ds.cli.DisasterRecovery.SubaccountPair.Get(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Subaccount Pair",
			fmt.Sprintf("An error was encountered when reading the subaccount pair with ID %q: %v", data.SubaccountId.ValueString(), err),
		)
		return
	}

	dataSourceValue, diags := SubaccountPairDataSourceValueFrom(ctx, data.SubaccountId, subaccountPair)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, dataSourceValue)
	resp.Diagnostics.Append(diags...)
}
