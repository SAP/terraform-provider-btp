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
		MarkdownDescription: `Get details about a resource provider instance.

__Tips__
You must be assigned to the global account admin or viewer role.

__Further documentation__
https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/e2c250dc5abd468a81f4f619206157a2.html`,
		Attributes: map[string]schema.Attribute{
			"resource_provider": schema.StringAttribute{
				MarkdownDescription: "Provider of the requested resource. For example: AWS, AZURE.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique technical name of the resource provider.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Descriptive name of the resource provider.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the resource provider.",
				Computed:            true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "Any relevant information about the resource provider that is not provided by other parameter values.",
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

	cliRes, _, err := ds.cli.Accounts.ResourceProvider.Get(ctx, data.ResourceProvider.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Resource Provider (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	data, diags = globalaccountResourceProviderValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
