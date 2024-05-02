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

func newSubaccountSecuritySettingsDataSource() datasource.DataSource {
	return &subaccountSecuritySettingsDataSource{}
}

type subaccountSecuritySettingsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountSecuritySettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_security_settings", req.ProviderTypeName)
}

func (ds *subaccountSecuritySettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountSecuritySettingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets the security settings of a subaccount.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/configure-trusted-domains-for-sap-authorization-and-trust-management-service>
<https://help.sap.com/docs/btp/sap-business-technology-platform/configure-token-policy-for-sap-authorization-and-trust-management-service>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"custom_email_domains": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Set of domains that are allowed to be used for user authentication.",
				Computed:            true,
			},
			"default_identity_provider": schema.StringAttribute{
				MarkdownDescription: "The subaccount's default identity provider for business application users.",
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

func (ds *subaccountSecuritySettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountSecuritySettingsType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Settings.ListBySubaccount(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Security Settings (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	state, diags := subaccountSecuritySettingsFromValue(ctx, cliRes)
	state.SubaccountId = data.SubaccountId
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
