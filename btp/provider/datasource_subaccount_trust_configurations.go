package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountTrustConfigurationsDataSource() datasource.DataSource {
	return &subaccountTrustConfigurationsDataSource{}
}

type subaccountTrustConfigurationsDataSourceConfig struct {
	Id           types.String                                `tfsdk:"id"`
	SubaccountId types.String                                `tfsdk:"subaccount_id"`
	Values       []subaccountTrustConfigurationListEntryType `tfsdk:"values"`
}

type subaccountTrustConfigurationsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountTrustConfigurationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_trust_configurations", req.ProviderTypeName)
}

func (ds *subaccountTrustConfigurationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountTrustConfigurationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets all trust configurations that are configured for your subaccount.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/trust-and-federation-with-identity-providers>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` attribute instead",
				MarkdownDescription: "The ID of the subaccount.",
				Computed:            true,
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"origin": schema.StringAttribute{
							MarkdownDescription: "The origin of the identity provider.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							DeprecationMessage:  "Use the `origin` attribute instead",
							MarkdownDescription: "The origin of the identity provider.",
							Computed:            true,
						},
						"identity_provider": schema.StringAttribute{
							MarkdownDescription: "The name of the Identity Authentication tenant the subaccount is connected to.",
							Computed:            true,
						},
						"domain": schema.StringAttribute{
							MarkdownDescription: "The tenant's domain which should be used for user logon.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The display name of the trust configuration.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the trust configuration.",
							Computed:            true,
						},
						"link_text": schema.StringAttribute{
							MarkdownDescription: "Short string that helps users to identify the link for login.",
							Computed:            true,
						},
						"available_for_user_logon": schema.BoolAttribute{
							MarkdownDescription: "Shows whether end users can choose the trust configuration for login. If not set, the trust configuration can remain active, however only application users that explicitly specify the origin key can use if for login.",
							Computed:            true,
						},
						"auto_create_shadow_users": schema.BoolAttribute{
							MarkdownDescription: "Shows whether any user from the tenant can log in. If not set, only the ones who already have a shadow user can log in.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "Shows whether the identity provider is currently 'active' or 'inactive'.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The trust type.",
							Computed:            true,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "The protocol used to establish trust with the identity provider.",
							Computed:            true,
						},
						"read_only": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the trust configuration can be modified.",
							Computed:            true,
						},
					},
				},
				MarkdownDescription: "The trust configurations associated with the subaccount.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountTrustConfigurationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountTrustConfigurationsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Trust.ListBySubaccount(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Trust Configurations (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.SubaccountId
	data.Values = []subaccountTrustConfigurationListEntryType{}

	for _, trustConfig := range cliRes {
		trustConfigValue, diags := subaccountTrustConfigurationFromListEntry(ctx, trustConfig)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, trustConfigValue)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
