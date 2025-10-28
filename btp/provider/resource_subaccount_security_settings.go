package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
You must be assigned to the admin role of the subaccount.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/configure-trusted-domains-for-sap-authorization-and-trust-management-service>
<https://help.sap.com/docs/btp/sap-business-technology-platform/configure-token-policy-for-sap-authorization-and-trust-management-service>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework for imports
				DeprecationMessage:  "Use the `subaccount_id`attribute instead",
				MarkdownDescription: "The ID of the security settings used for import operations.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_email_domains": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Set of domains that are allowed to be used for user authentication.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"default_identity_provider": schema.StringAttribute{
				MarkdownDescription: "The subaccount's default identity provider for business application users.",
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
				DeprecationMessage:  "Use the `iframe_domains_list` attribute instead",
				MarkdownDescription: "The new domains of the iframe. Enter as string. To provide multiple domains, separate them by spaces.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^(|.{4,})$`), "The attribute iframe_domains must be empty or contain domains."),
					stringvalidator.ConflictsWith(path.MatchRoot("iframe_domains"), path.MatchRoot("iframe_domains_list")),
				},
			},
			"iframe_domains_list": schema.ListAttribute{
				MarkdownDescription: "The new domains of the iframe. Enter as list. It is recommended to use in place of iframe_domains as list of iframes is better managed by this parameter.",
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Validators: []validator.List{
					listvalidator.ConflictsWith(
						path.MatchRoot("iframe_domains"),
						path.MatchRoot("iframe_domains_list"),
					),
					listvalidator.ValueStringsAre(
						stringvalidator.RegexMatches(regexp.MustCompile(`^([^ ]{4,})$`), "the attribute iframe_domains_list must contain valid domains."),
					),
				},
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
		handleReadErrors(ctx, rawRes, cliRes, resp, err, "Resource Security Settings (Subaccount)")
		return
	}

	var transferIframeString bool
	if isIFrameDomainsSet(state.IframeDomains) {
		transferIframeString = true
	} else {
		// During IMPORT we must make sure that only one iframe attribute is filled.
		// Precedence is given to the list attribute, so we clear the computed iframe string attribute
		// This causes errors when the configuration contains the deprecated value for the iframe string attribute which is intended.
		transferIframeString = false
	}

	updatedState, diags := subaccountSecuritySettingsValueFrom(ctx, cliRes, transferIframeString)
	updatedState.SubaccountId = state.SubaccountId

	if updatedState.Id.IsNull() || updatedState.Id.IsUnknown() {
		// Setting ID of state - required by hashicorps terraform plugin testing framework for Import . See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
		updatedState.Id = state.SubaccountId
	}
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

	var customEmailDomains []string

	diags = plan.CustomEmailDomains.ElementsAs(ctx, &customEmailDomains, false)
	resp.Diagnostics.Append(diags...)

	iFrameDomains := plan.IframeDomains.ValueString()
	if !plan.IframeDomainsList.IsUnknown() {
		if !plan.IframeDomainsList.IsNull() {
			var domains []string
			diags := plan.IframeDomainsList.ElementsAs(ctx, &domains, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			iFrameDomains = strings.Join(domains, " ")
		}
	}

	res, _, err := rs.cli.Security.Settings.UpdateBySubaccount(ctx, plan.SubaccountId.ValueString(), btpcli.SecuritySettingsUpdateInput{
		CustomEmail:                       customEmailDomains,
		DefaultIDPForNonInteractiveLogon:  plan.DefaultIdentityProvider.ValueString(),
		TreatUsersWithSameEmailAsSameUser: plan.TreatUsersWithSameEmailAsSameUser.ValueBool(),
		AccessTokenValidity:               int(plan.AccessTokenValidity.ValueInt64()),
		RefreshTokenValidity:              int(plan.RefreshTokenValidity.ValueInt64()),
		IFrame:                            iFrameDomains,
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Security Settings (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	var transferIframeString bool
	if isIFrameDomainsSet(plan.IframeDomains) {
		// the plan contains a value for the iframe string attribute, so we must transfer it
		transferIframeString = true
	} else {
		transferIframeString = false
	}

	state, diags := subaccountSecuritySettingsValueFrom(ctx, res, transferIframeString)
	state.SubaccountId = plan.SubaccountId
	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	state.Id = plan.SubaccountId

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

	var customEmailDomains []string
	diags = plan.CustomEmailDomains.ElementsAs(ctx, &customEmailDomains, false)
	resp.Diagnostics.Append(diags...)

	planIFrameDomains := plan.IframeDomains.ValueString()
	stateIFrameDomains := state.IframeDomains.ValueString()
	if !plan.IframeDomainsList.IsUnknown() {
		if !plan.IframeDomainsList.IsNull() {
			var planDomains []string
			var stateDomains []string
			diags := plan.IframeDomainsList.ElementsAs(ctx, &planDomains, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			planIFrameDomains = strings.Join(planDomains, " ")
			if len(planIFrameDomains) == 0 {
				planIFrameDomains = ""
			}
			diags = state.IframeDomainsList.ElementsAs(ctx, &stateDomains, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
			stateIFrameDomains = strings.Join(stateDomains, " ")
			if len(stateIFrameDomains) == 0 {
				stateIFrameDomains = ""
			}
		}
	}

	iFrameDomain := transformIframeDomain(planIFrameDomains, stateIFrameDomains)

	res, _, err := rs.cli.Security.Settings.UpdateBySubaccount(ctx, plan.SubaccountId.ValueString(), btpcli.SecuritySettingsUpdateInput{
		CustomEmail:                       customEmailDomains,
		DefaultIDPForNonInteractiveLogon:  plan.DefaultIdentityProvider.ValueString(),
		TreatUsersWithSameEmailAsSameUser: plan.TreatUsersWithSameEmailAsSameUser.ValueBool(),
		AccessTokenValidity:               int(plan.AccessTokenValidity.ValueInt64()),
		RefreshTokenValidity:              int(plan.RefreshTokenValidity.ValueInt64()),
		IFrame:                            iFrameDomain,
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Security Settings (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	var transferIframeString bool
	if isIFrameDomainsSet(plan.IframeDomains) || isIFrameDomainsSet(state.IframeDomains) {
		// a value for the iframe string attribute is present in either plan or state, so we transfer it
		transferIframeString = true
	} else {
		// default is that we do not transfer the iframe string attribute
		transferIframeString = false
	}

	state, diags = subaccountSecuritySettingsValueFrom(ctx, res, transferIframeString)
	state.SubaccountId = plan.SubaccountId
	state.Id = plan.SubaccountId
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
		CustomEmail:                       []string{},
		DefaultIDPForNonInteractiveLogon:  "sap.default",
		TreatUsersWithSameEmailAsSameUser: false,
		AccessTokenValidity:               -1,
		RefreshTokenValidity:              -1,
		IFrame:                            " ", // The string should be empty, however to do the update the value must be " " (space)
	})

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Security Settings (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountSecuritySettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), req.ID)...)
}
