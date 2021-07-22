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

func newSubaccountAppsDataSource() datasource.DataSource {
	return &subaccountAppsDataSource{}
}

type subaccountAppsValue struct {
	Id                     types.String                        `tfsdk:"id"`
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

type subaccountAppsDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	/* OUTPUT */
	Values []subaccountAppsValue `tfsdk:"values"`
}

type subaccountAppsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountAppsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_apps", req.ProviderTypeName)
}

func (ds *subaccountAppsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountAppsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `List all apps.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The application ID is the xsappname plus the identifier, which consists of an exclamation mark (!), an identifier for the plan underwhich the application is deployed, and an index number.",
							Computed:            true,
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
				},
				Computed: true,
			},
		},
	}
}

func (ds *subaccountAppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountAppsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Security.App.ListBySubaccount(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Apps (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Values = []subaccountAppsValue{}
	for _, app := range cliRes {
		appVal := subaccountAppsValue{}

		appVal.Id = types.StringValue(app.Appid)

		appVal.Description = types.StringValue(app.Description)

		if app.MasterAppId == nil {
			appVal.MasterAppId = types.StringNull()
		} else {
			appVal.MasterAppId = types.StringValue(*app.MasterAppId)
		}

		appVal.OrgId = types.StringValue(app.OrgId)
		appVal.PlanId = types.StringValue(app.PlanId)
		appVal.PlanName = types.StringValue(app.PlanName)
		appVal.ServiceinstanceId = types.StringValue(app.Serviceinstanceid)
		appVal.Xsappname = types.StringValue(app.Xsappname)

		if app.SpaceId == nil {
			appVal.SpaceId = types.StringNull()
		} else {
			appVal.SpaceId = types.StringValue(*app.SpaceId)
		}
		appVal.TenantMode = types.StringValue(app.TenantMode)

		if app.UserName == nil {
			appVal.Username = types.StringNull()
		} else {
			appVal.Username = types.StringValue(*app.UserName)
		}

		if app.Oauth2Configuration != nil {
			appVal.Oauth2Configuration = &globalaccountAppOauthConfiguration{
				Autoapprove:          types.BoolValue(app.Oauth2Configuration.Autoapprove),
				RefreshTokenValidity: types.Int64Value(int64(app.Oauth2Configuration.RefreshTokenValidity)),
				TokenValidity:        types.Int64Value(int64(app.Oauth2Configuration.TokenValidity)),
			}

			appVal.Oauth2Configuration.Allowedproviders, diags = types.SetValueFrom(ctx, types.StringType, app.Oauth2Configuration.Allowedproviders)
			resp.Diagnostics.Append(diags...)

			appVal.Oauth2Configuration.GrantTypes, diags = types.SetValueFrom(ctx, types.StringType, app.Oauth2Configuration.GrantTypes)
			resp.Diagnostics.Append(diags...)

			appVal.Oauth2Configuration.RedirectUris, diags = types.SetValueFrom(ctx, types.StringType, app.Oauth2Configuration.RedirectUris)
			resp.Diagnostics.Append(diags...)

			appVal.Oauth2Configuration.SystemAttributes, diags = types.SetValueFrom(ctx, types.StringType, app.Oauth2Configuration.SystemAttributes)
			resp.Diagnostics.Append(diags...)
		}

		appVal.Authorities, diags = types.SetValueFrom(ctx, types.StringType, app.Authorities)
		resp.Diagnostics.Append(diags...)

		appVal.ForeignScopeReferences, diags = types.SetValueFrom(ctx, types.StringType, app.ForeignScopeReferences)
		resp.Diagnostics.Append(diags...)

		data.Values = append(data.Values, appVal)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
