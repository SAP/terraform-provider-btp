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

func newGlobalaccountTrustConfigurationDataSource() datasource.DataSource {
	return &globalaccountTrustConfigurationDataSource{}
}

type globalaccountTrustConfigurationDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountTrustConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_trust_configuration", req.ProviderTypeName)
}

func (ds *globalaccountTrustConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountTrustConfigurationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Get details about a trust configuration.

__Tip:__
You must be viewer or administrator of the global account.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-btp-neo-environment/platform-identity-provider>`,
		Attributes: map[string]schema.Attribute{
			"origin": schema.StringAttribute{
				MarkdownDescription: "The origin of the identity provider.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the trust configuration.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the trust configuration.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the trust configuration.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The trust type.",
				Computed:            true,
			},
			"identity_provider": schema.StringAttribute{
				MarkdownDescription: "The name of the identity provider.",
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol used to establish trust with the identity provider.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Whether the identity provider is currently active or not.",
				Computed:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Whether the trust configuration can be modified.",
				Computed:            true,
			},
		},
	}
}

func (ds *globalaccountTrustConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountTrustConfigurationType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.Trust.GetByGlobalAccount(ctx, data.Origin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Trust Configuration (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	data, diags = globalaccountTrustConfigurationFromValue(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
