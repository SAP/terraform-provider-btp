package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountSecuritySettingsDataSource() datasource.DataSource {
	return &globalaccountSecuritySettingsDataSource{}
}

type globalaccountSecuritySettingsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountSecuritySettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_security_settings", req.ProviderTypeName)
}

func (ds *globalaccountSecuritySettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountSecuritySettingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets the security settings of a global account.

__Tip:__
You must be viewer or administrator of the global account.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/configure-trusted-domains-for-sap-authorization-and-trust-management-service>
<https://help.sap.com/docs/btp/sap-business-technology-platform/configure-token-policy-for-sap-authorization-and-trust-management-service>`,
		Attributes: map[string]schema.Attribute{
			"custom_email_domains": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Set of domains that are allowed to be used for user authentication.",
				Computed:            true,
			},
			"default_identity_provider": schema.StringAttribute{
				MarkdownDescription: "The global account's default identity provider for platform users. Used to log on to platform tools such as SAP BTP cockpit or the btp CLI.",
				Computed:            true,
			},
			"treat_users_with_same_email_as_same_user": schema.BoolAttribute{
				MarkdownDescription: "If set to true, users with the same email are treated as same users.",
				Computed:            true,
			},
			"access_token_validity": schema.Int64Attribute{
				MarkdownDescription: "The validity of the access token.",
				Computed:            true,
			},
			"refresh_token_validity": schema.Int64Attribute{
				MarkdownDescription: "The validity of the refresh token.",
				Computed:            true,
			},
		},
	}
}

func (ds *globalaccountSecuritySettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountSecuritySettingsType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Settings.ListByGlobalAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Security Settings (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	data, diags = globalaccountSecuritySettingsFromValue(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
