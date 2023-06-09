package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountUsersDataSource() datasource.DataSource {
	return &globalaccountUsersDataSource{}
}

type globalaccountUsersDataSourceConfig struct {
	/* INPUT */
	Origin types.String `tfsdk:"origin"`
	Id     types.String `tfsdk:"id"`
	/* OUTPUT */
	Values types.Set `tfsdk:"values"`
}

type globalaccountUsersDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountUsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_users", req.ProviderTypeName)
}

func (ds *globalaccountUsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountUsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `List all users.

__Further documentation:__
https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/cc1c676b43904066abb2a4838cbd0c37.html`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				DeprecationMessage:  "Use the `btp_globalaccount` datasource instead",
				MarkdownDescription: "The ID of the global account.",
				Computed:            true,
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The identity provider that hosts the user. The default value is 'ldap'.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"values": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The list of users assigned to the global account.",
				Computed:            true,
			},
		},
	}
}

func (ds *globalaccountUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountUsersDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Origin.IsNull() {
		data.Origin = types.StringValue("ldap")
	}

	cliRes, _, err := ds.cli.Security.User.ListByGlobalAccount(ctx, data.Origin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Users (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = types.StringValue(ds.cli.GetGlobalAccountSubdomain())
	data.Values, diags = types.SetValueFrom(ctx, types.StringType, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
