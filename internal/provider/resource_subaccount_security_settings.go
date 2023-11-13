package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountSecuritySettingsResource() resource.Resource {
	return &subaccountSecuritySettingsResource{}
}

type subaccountSecuritySettingsResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountSecuritySettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_security_settings", req.ProviderTypeName)
}

func (rs *subaccountSecuritySettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountSecuritySettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Sets the security settings of a subaccount.

__Tip:__
You must be administrator of the subaccount.

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
				MarkdownDescription: "Set of domains which are allowed to be used for user authentication.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Set{
					setvalidator.IsRequired(),
				},
			},
			"default_identity_provider": schema.StringAttribute{
				MarkdownDescription: "The default identity provider which is used for noninteractive logon.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("sap.default"),
			},
			"treat_users_with_same_email_as_same_user": schema.BoolAttribute{
				MarkdownDescription: "If set to true, users with the same email are treated as same users.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"access_token_validity": schema.Int64Attribute{
				MarkdownDescription: "The validity of the access token.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(int64(-1)),
			},
			"refresh_token_validity": schema.Int64Attribute{
				MarkdownDescription: "The validity of the refresh token.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(int64(-1)),
			},
		},
	}
}

func (rs *subaccountSecuritySettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountSecuritySettingsType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Security.Settings.ListBySubaccount(ctx, state.SubaccountId.ValueString())

	if err != nil {
		handleReadErrors(ctx, rawRes, resp, err, "Resource Security Settings (Subaccount)")
		return
	}

	updatedState, diags := subaccountSecuritySettingsFromValue(ctx, cliRes)
	updatedState.SubaccountId = state.SubaccountId
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountSecuritySettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountSecuritySettingsType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := rs.cli.Security.Settings.UpdateBySubaccount(ctx, plan.SubaccountId.ValueString(), btpcli.SecuritySettingsUpdateInput{
		//CustomEmail: "[]",
		DefaultIDPForNonInteractiveLogon:  plan.DefaultIdentityProvider.ValueString(),
		TreatUsersWithSameEmailAsSameUser: plan.TreatUsersWithSameEmailAsSameUser.ValueBool(),
		AccessTokenValidity:               int(plan.AccessTokenValidity.ValueInt64()),
		RefreshTokenValidity:              int(plan.RefreshTokenValidity.ValueInt64()),
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Security Settings (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	state, diags := subaccountSecuritySettingsFromValue(ctx, res)
	state.SubaccountId = plan.SubaccountId
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountSecuritySettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountSecuritySettingsType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	var state subaccountSecuritySettingsType
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	res, _, err := rs.cli.Security.Settings.UpdateBySubaccount(ctx, plan.SubaccountId.ValueString(), btpcli.SecuritySettingsUpdateInput{
		//CustomEmail: "[]",
		DefaultIDPForNonInteractiveLogon:  plan.DefaultIdentityProvider.ValueString(),
		TreatUsersWithSameEmailAsSameUser: plan.TreatUsersWithSameEmailAsSameUser.ValueBool(),
		AccessTokenValidity:               int(plan.AccessTokenValidity.ValueInt64()),
		RefreshTokenValidity:              int(plan.RefreshTokenValidity.ValueInt64()),
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Security Settings (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	state, diags = subaccountSecuritySettingsFromValue(ctx, res)
	state.SubaccountId = plan.SubaccountId
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountSecuritySettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountSecuritySettingsType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.Settings.UpdateBySubaccount(ctx, state.SubaccountId.ValueString(), btpcli.SecuritySettingsUpdateInput{
		CustomEmail:                       "[]",
		DefaultIDPForNonInteractiveLogon:  "sap.default",
		TreatUsersWithSameEmailAsSameUser: false,
		AccessTokenValidity:               -1,
		RefreshTokenValidity:              -1,
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Security Settings (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}
