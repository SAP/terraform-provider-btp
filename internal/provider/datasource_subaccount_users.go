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
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountUsersDataSource() datasource.DataSource {
	return &subaccountUsersDataSource{}
}

type subaccountUsersDataSourceConfig struct {
	/* INPUT */
	Id           types.String `tfsdk:"id"`
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Origin       types.String `tfsdk:"origin"`
	Values       types.Set    `tfsdk:"values"`
}

type subaccountUsersDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountUsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_users", req.ProviderTypeName)
}

func (ds *subaccountUsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountUsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets all users.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/user-and-member-management>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` attribute instead",
				MarkdownDescription: "The ID of the subaccount.",
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
				MarkdownDescription: "The list of users assigned to the subaccount.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountUsersDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Origin.IsNull() {
		data.Origin = types.StringValue("ldap")
	}

	cliRes, _, err := ds.cli.Security.User.ListBySubaccount(ctx, data.SubaccountId.ValueString(), data.Origin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Users (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.SubaccountId
	data.Values, diags = types.SetValueFrom(ctx, types.StringType, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
