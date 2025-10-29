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

func newSubaccountUserDataSource() datasource.DataSource {
	return &subaccountUserDataSource{}
}

type subaccountUserDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Origin       types.String `tfsdk:"origin"`
	UserName     types.String `tfsdk:"user_name"`
	/* OUTPUT */
	Id              types.String `tfsdk:"id"`
	Email           types.String `tfsdk:"email"`
	GivenName       types.String `tfsdk:"given_name"`
	FamilyName      types.String `tfsdk:"family_name"`
	Verified        types.Bool   `tfsdk:"verified"`
	Active          types.Bool   `tfsdk:"active"`
	RoleCollections types.Set    `tfsdk:"role_collections"`
}

type subaccountUserDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_user", req.ProviderTypeName)
}

func (ds *subaccountUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Shows registered users in a subaccount. Users belong to one of the identity providers (IdPs) of the subaccount.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The identity provider that hosts the user. Only needed for custom identity provider.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"user_name": schema.StringAttribute{
				MarkdownDescription: "The username of the user.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user.",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The e-mail address of the user.",
				Computed:            true,
			},
			"given_name": schema.StringAttribute{
				MarkdownDescription: "The given name of the user.",
				Computed:            true,
			},
			"family_name": schema.StringAttribute{
				MarkdownDescription: "The last name of the user.",
				Computed:            true,
			},
			"verified": schema.BoolAttribute{
				MarkdownDescription: "The verification status of the user.",
				Computed:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Shows if the account is still in use.",
				Computed:            true,
			},
			"role_collections": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The set of role collections, which are assigned to the user.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountUserDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Origin.IsNull() {
		data.Origin = types.StringValue("ldap")
	}

	cliRes, _, err := ds.cli.Security.User.GetBySubaccount(ctx, data.SubaccountId.ValueString(), data.UserName.ValueString(), data.Origin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource User (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = types.StringValue(cliRes.Id)
	data.Email = types.StringValue(cliRes.Email)
	data.GivenName = types.StringValue(cliRes.GivenName)
	data.FamilyName = types.StringValue(cliRes.FamilyName)
	data.Verified = types.BoolValue(cliRes.Verified)
	data.Active = types.BoolValue(cliRes.Active)

	data.RoleCollections, diags = types.SetValueFrom(ctx, types.StringType, cliRes.RoleCollections)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
