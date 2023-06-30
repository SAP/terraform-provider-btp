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

func newDirectoryAppDataSource() datasource.DataSource {
	return &directoryAppDataSource{}
}

type directoryAppDataSourceConfig struct {
	/* INPUT */
	DirectoryId types.String `tfsdk:"directory_id"`
	Id          types.String `tfsdk:"id"`
	/* OUTPUT */
	Authorities            types.Set                           `tfsdk:"authorities"`
	Description            types.String                        `tfsdk:"description"`
	ForeignScopeReferences types.Set                           `tfsdk:"foreign_scope_references"`
	MasterAppId            types.String                        `tfsdk:"master_app_id"`
	Oauth2Configuration    *globalaccountAppOauthConfiguration `tfsdk:"oauth2_configuration"`
	OrgId                  types.String                        `tfsdk:"org_id"`
	PlanId                 types.String                        `tfsdk:"plan_id"`
	PlanName               types.String                        `tfsdk:"plan_name"`
	ServiceinstanceId      types.String                        `tfsdk:"serviceinstance_id"`
	SpaceId                types.String                        `tfsdk:"space_id"`
	TenantMode             types.String                        `tfsdk:"tenant_mode"`
	Username               types.String                        `tfsdk:"username"`
	Xsappname              types.String                        `tfsdk:"xsappname"`
}

type directoryAppDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *directoryAppDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_app", req.ProviderTypeName)
}

func (ds *directoryAppDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *directoryAppDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific app.`,
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The application ID is the xsappname plus the identifier, which consists of an exclamation mark (!), an identifier for the plan underwhich the application is deployed, and an index number.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"authorities": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the app.",
				Computed:            true,
			},
			"foreign_scope_references": schema.SetAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"master_app_id": schema.StringAttribute{
				Computed: true,
			},
			"oauth2_configuration": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"allowedproviders": schema.SetAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"autoapprove": schema.BoolAttribute{
						Computed: true,
					},
					"grant_types": schema.SetAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"redirect_uris": schema.SetAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"refresh_token_validity": schema.Int64Attribute{
						Computed: true,
					},
					"system_attributes": schema.SetAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
					"token_validity": schema.Int64Attribute{
						Computed: true,
					},
				},
				Computed: true,
			},
			"org_id": schema.StringAttribute{
				Computed: true,
			},
			"plan_id": schema.StringAttribute{
				Computed: true,
			},
			"plan_name": schema.StringAttribute{
				Computed: true,
			},
			"serviceinstance_id": schema.StringAttribute{
				Computed: true,
			},
			"space_id": schema.StringAttribute{
				Computed: true,
			},
			"tenant_mode": schema.StringAttribute{
				Computed: true,
			},
			"username": schema.StringAttribute{
				Computed: true,
			},
			"xsappname": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (ds *directoryAppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryAppDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.App.GetByDirectory(ctx, data.DirectoryId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource App (Directory)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = types.StringValue(cliRes.Appid)

	data.Description = types.StringValue(cliRes.Description)

	if cliRes.MasterAppId == nil {
		data.MasterAppId = types.StringNull()
	} else {
		data.MasterAppId = types.StringValue(*cliRes.MasterAppId)
	}

	data.OrgId = types.StringValue(cliRes.OrgId)
	data.PlanId = types.StringValue(cliRes.PlanId)
	data.PlanName = types.StringValue(cliRes.PlanName)
	data.ServiceinstanceId = types.StringValue(cliRes.Serviceinstanceid)
	data.Xsappname = types.StringValue(cliRes.Xsappname)

	if cliRes.SpaceId == nil {
		data.SpaceId = types.StringNull()
	} else {
		data.SpaceId = types.StringValue(*cliRes.SpaceId)
	}
	data.TenantMode = types.StringValue(cliRes.TenantMode)

	if cliRes.UserName == nil {
		data.Username = types.StringNull()
	} else {
		data.Username = types.StringValue(*cliRes.UserName)
	}

	if cliRes.Oauth2Configuration != nil {
		data.Oauth2Configuration = &globalaccountAppOauthConfiguration{
			Autoapprove:          types.BoolValue(cliRes.Oauth2Configuration.Autoapprove),
			RefreshTokenValidity: types.Int64Value(int64(cliRes.Oauth2Configuration.RefreshTokenValidity)),
			TokenValidity:        types.Int64Value(int64(cliRes.Oauth2Configuration.TokenValidity)),
		}

		data.Oauth2Configuration.Allowedproviders, diags = types.SetValueFrom(ctx, types.StringType, cliRes.Oauth2Configuration.Allowedproviders)
		resp.Diagnostics.Append(diags...)

		data.Oauth2Configuration.GrantTypes, diags = types.SetValueFrom(ctx, types.StringType, cliRes.Oauth2Configuration.GrantTypes)
		resp.Diagnostics.Append(diags...)

		data.Oauth2Configuration.RedirectUris, diags = types.SetValueFrom(ctx, types.StringType, cliRes.Oauth2Configuration.RedirectUris)
		resp.Diagnostics.Append(diags...)

		data.Oauth2Configuration.SystemAttributes, diags = types.SetValueFrom(ctx, types.StringType, cliRes.Oauth2Configuration.SystemAttributes)
		resp.Diagnostics.Append(diags...)
	}

	data.Authorities, diags = types.SetValueFrom(ctx, types.StringType, cliRes.Authorities)
	resp.Diagnostics.Append(diags...)

	data.ForeignScopeReferences, diags = types.SetValueFrom(ctx, types.StringType, cliRes.ForeignScopeReferences)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
