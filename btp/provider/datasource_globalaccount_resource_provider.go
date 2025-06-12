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

func newGlobalaccountResourceProviderDataSource() datasource.DataSource {
	return &globalaccountResourceProviderDataSource{}
}

type globalaccountResourceProviderDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountResourceProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_resource_provider", req.ProviderTypeName)
}

func (ds *globalaccountResourceProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountResourceProviderDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a resource provider instance.

__Tip:__
You must be assigned to the global account admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/managing-resource-providers>`,
		Attributes: map[string]schema.Attribute{
			"provider_type": schema.StringAttribute{
				MarkdownDescription: "The cloud vendor from which to consume services through your subscribed account. Possible values are: \n" +
					getFormattedValueAsTableRow("value", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`AWS`", "Amazon Web Services") +
					getFormattedValueAsTableRow("`AZURE`", "Microsoft Azure"),
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"technical_name": schema.StringAttribute{
				MarkdownDescription: "The unique technical name of the resource provider.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				DeprecationMessage:  "Use the `technical_name` attribute instead",
				MarkdownDescription: "The unique technical name of the resource provider.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The descriptive name of the resource provider.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the resource provider.",
				Computed:            true,
			},
			"configuration": schema.StringAttribute{
				MarkdownDescription: "The configuration properties for the resource provider as required by the vendor.",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (ds *globalaccountResourceProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountResourceProviderType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.ResourceProvider.Get(ctx, data.Provider.ValueString(), data.TechnicalName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Resource Provider (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	data, diags = globalaccountResourceProviderValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
