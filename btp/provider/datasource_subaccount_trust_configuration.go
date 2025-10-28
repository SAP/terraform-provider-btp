package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountTrustConfigurationDataSource() datasource.DataSource {
	return &subaccountTrustConfigurationDataSource{}
}

type subaccountTrustConfigurationDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountTrustConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_trust_configuration", req.ProviderTypeName)
}

func (ds *subaccountTrustConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountTrustConfigurationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a trust configuration.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-btp-neo-environment/platform-identity-provider>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The origin of the identity provider.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
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
	}
}

func (ds *subaccountTrustConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountTrustConfigurationType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Trust.GetBySubaccount(ctx, data.SubaccountId.ValueString(), data.Origin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Trust Configuration (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	state, diags := subaccountTrustConfigurationFromValue(ctx, cliRes)
	state.SubaccountId = data.SubaccountId
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
