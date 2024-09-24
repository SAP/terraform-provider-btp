package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountSecuritySettingsResource() resource.Resource {
	return &globalaccountSecuritySettingsResource{}
}

type globalaccountSecuritySettingsResource struct {
	cli *btpcli.ClientFacade
}

func (rs *globalaccountSecuritySettingsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_security_settings", req.ProviderTypeName)
}

func (rs *globalaccountSecuritySettingsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *globalaccountSecuritySettingsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Sets the security settings of a global account.

__Tip:__
You must be assigned to the admin role of the global account.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/configure-trusted-domains-for-sap-authorization-and-trust-management-service>
<https://help.sap.com/docs/btp/sap-business-technology-platform/configure-token-policy-for-sap-authorization-and-trust-management-service>`,
		Attributes: map[string]schema.Attribute{
			"custom_email_domains": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Set of domains that are allowed to be used for user authentication.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"default_identity_provider": schema.StringAttribute{
				MarkdownDescription: "The global account's default identity provider for platform users. Used to log on to platform tools such as SAP BTP cockpit or the btp CLI.",
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
			"iframe_domains": schema.StringAttribute{
				MarkdownDescription: "The new domains of the iframe. Enter as string. To provide multiple domains, separate them by spaces.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^(|.{4,})$`), "The attribute iframe_domains must be empty or contain domains."),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework for imports
				DeprecationMessage:  "Automatically filled with the subdomain of the global account",
				MarkdownDescription: "The ID of the security settings used for import operations.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (rs *globalaccountSecuritySettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state globalaccountSecuritySettingsType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Security.Settings.ListByGlobalAccount(ctx)
	if err != nil {
		handleReadErrors(ctx, rawRes, resp, err, "Resource Security Settings (Global Account)")
		return
	}

	state, diags = globalaccountSecuritySettingsValueFrom(ctx, cliRes)

	if state.Id.IsNull() || state.Id.IsUnknown() {
		// Setting ID of state - required by hashicorps terraform plugin testing framework for Import . See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
		state.Id = types.StringValue(rs.cli.GetGlobalAccountSubdomain())
	}

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountSecuritySettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan globalaccountSecuritySettingsType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var customEmailDomains []string

	diags = plan.CustomEmailDomains.ElementsAs(ctx, &customEmailDomains, false)
	resp.Diagnostics.Append(diags...)

	res, _, err := rs.cli.Security.Settings.UpdateByGlobalAccount(ctx, btpcli.SecuritySettingsUpdateInput{
		CustomEmail:                       customEmailDomains,
		DefaultIDPForNonInteractiveLogon:  plan.DefaultIdentityProvider.ValueString(),
		TreatUsersWithSameEmailAsSameUser: plan.TreatUsersWithSameEmailAsSameUser.ValueBool(),
		AccessTokenValidity:               int(plan.AccessTokenValidity.ValueInt64()),
		RefreshTokenValidity:              int(plan.RefreshTokenValidity.ValueInt64()),
		IFrame:                            plan.IframeDomains.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Security Settings (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	state, diags := globalaccountSecuritySettingsValueFrom(ctx, res)
	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	state.Id = types.StringValue(rs.cli.GetGlobalAccountSubdomain())

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountSecuritySettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan globalaccountSecuritySettingsType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var currentState globalaccountSecuritySettingsType
	diags = req.State.Get(ctx, &currentState)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var customEmailDomains []string

	diags = plan.CustomEmailDomains.ElementsAs(ctx, &customEmailDomains, false)
	resp.Diagnostics.Append(diags...)

	iFrameDomain := transformIframeDomain(plan.IframeDomains.ValueString(), currentState.IframeDomains.ValueString())

	res, _, err := rs.cli.Security.Settings.UpdateByGlobalAccount(ctx, btpcli.SecuritySettingsUpdateInput{
		CustomEmail:                       customEmailDomains,
		DefaultIDPForNonInteractiveLogon:  plan.DefaultIdentityProvider.ValueString(),
		TreatUsersWithSameEmailAsSameUser: plan.TreatUsersWithSameEmailAsSameUser.ValueBool(),
		AccessTokenValidity:               int(plan.AccessTokenValidity.ValueInt64()),
		RefreshTokenValidity:              int(plan.RefreshTokenValidity.ValueInt64()),
		IFrame:                            iFrameDomain,
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Security Settings (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	state, diags := globalaccountSecuritySettingsValueFrom(ctx, res)
	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	state.Id = currentState.Id

	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountSecuritySettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state globalaccountSecuritySettingsType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.Settings.UpdateByGlobalAccount(ctx, btpcli.SecuritySettingsUpdateInput{
		CustomEmail:                       []string{},
		DefaultIDPForNonInteractiveLogon:  "sap.default",
		TreatUsersWithSameEmailAsSameUser: false,
		AccessTokenValidity:               -1,
		RefreshTokenValidity:              -1,
		IFrame:                            " ", // The string should be empty, however to do the update the value must be " " (space)
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Security Settings (Global Account)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *globalaccountSecuritySettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
