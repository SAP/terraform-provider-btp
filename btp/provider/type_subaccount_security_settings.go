package provider

import (
	"context"
	"strings"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_settings"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type subaccountSecuritySettingsType struct {
	SubaccountId                      types.String `tfsdk:"subaccount_id"`
	Id                                types.String `tfsdk:"id"`
	CustomEmailDomains                types.Set    `tfsdk:"custom_email_domains"`
	DefaultIdentityProvider           types.String `tfsdk:"default_identity_provider"`
	TreatUsersWithSameEmailAsSameUser types.Bool   `tfsdk:"treat_users_with_same_email_as_same_user"`
	AccessTokenValidity               types.Int64  `tfsdk:"access_token_validity"`
	RefreshTokenValidity              types.Int64  `tfsdk:"refresh_token_validity"`
	IframeDomains                     types.String `tfsdk:"iframe_domains"`
	IframeDomainsList                 types.List   `tfsdk:"iframe_domains_list"`
}

func subaccountSecuritySettingsValueFrom(ctx context.Context, value xsuaa_settings.TenantSettingsResp, transferIframestring bool) (tenantSettings subaccountSecuritySettingsType, diags diag.Diagnostics) {
	tenantSettings.TreatUsersWithSameEmailAsSameUser = types.BoolValue(value.TreatUsersWithSameEmailAsSameUser)

	if len(value.DefaultIdp) > 0 {
		tenantSettings.DefaultIdentityProvider = types.StringValue(value.DefaultIdp)
	} else {
		tenantSettings.DefaultIdentityProvider = types.StringNull()
	}

	if value.TokenPolicySettings != nil {
		tenantSettings.AccessTokenValidity = types.Int64Value(int64(value.TokenPolicySettings.AccessTokenValidity))
		tenantSettings.RefreshTokenValidity = types.Int64Value(int64(value.TokenPolicySettings.RefreshTokenValidity))
	}

	if len(value.CustomEmailDomains) > 0 {
		tenantSettings.CustomEmailDomains, diags = types.SetValueFrom(ctx, types.StringType, value.CustomEmailDomains)
	} else {
		tenantSettings.CustomEmailDomains, diags = types.SetValueFrom(ctx, types.StringType, []string{})
	}

	if transferIframestring {
		// Transfer of deprecated attribute to the structure only if needed
		tenantSettings.IframeDomains = types.StringValue(value.IframeDomains)
	}

	iframeDomainsList := []string{}
	if value.IframeDomains != "" {
		iframeDomainsList = strings.Fields(value.IframeDomains)
	}
	tenantSettings.IframeDomainsList, _ = types.ListValueFrom(ctx, types.StringType, iframeDomainsList)

	return
}

type subaccountSecuritySettingsDataSourceType struct {
	SubaccountId                      types.String `tfsdk:"subaccount_id"`
	CustomEmailDomains                types.Set    `tfsdk:"custom_email_domains"`
	DefaultIdentityProvider           types.String `tfsdk:"default_identity_provider"`
	TreatUsersWithSameEmailAsSameUser types.Bool   `tfsdk:"treat_users_with_same_email_as_same_user"`
	AccessTokenValidity               types.Int64  `tfsdk:"access_token_validity"`
	RefreshTokenValidity              types.Int64  `tfsdk:"refresh_token_validity"`
	IframeDomains                     types.String `tfsdk:"iframe_domains"`
	IframDomainsList                  types.List   `tfsdk:"iframe_domains_list"`
}

func subaccountSecuritySettingsDataSourceValueFrom(ctx context.Context, value xsuaa_settings.TenantSettingsResp) (tenantSettings subaccountSecuritySettingsDataSourceType, diags diag.Diagnostics) {
	tenantSettings.TreatUsersWithSameEmailAsSameUser = types.BoolValue(value.TreatUsersWithSameEmailAsSameUser)

	if len(value.DefaultIdp) > 0 {
		tenantSettings.DefaultIdentityProvider = types.StringValue(value.DefaultIdp)
	} else {
		tenantSettings.DefaultIdentityProvider = types.StringNull()
	}

	if value.TokenPolicySettings != nil {
		tenantSettings.AccessTokenValidity = types.Int64Value(int64(value.TokenPolicySettings.AccessTokenValidity))
		tenantSettings.RefreshTokenValidity = types.Int64Value(int64(value.TokenPolicySettings.RefreshTokenValidity))
	}

	if len(value.CustomEmailDomains) > 0 {
		tenantSettings.CustomEmailDomains, diags = types.SetValueFrom(ctx, types.StringType, value.CustomEmailDomains)
	} else {
		tenantSettings.CustomEmailDomains, diags = types.SetValueFrom(ctx, types.StringType, []string{})
	}

	tenantSettings.IframeDomains = types.StringValue(value.IframeDomains)

	iframeDomainsList := []string{}
	if value.IframeDomains != "" {
		iframeDomainsList = strings.Fields(value.IframeDomains)
	}

	tenantSettings.IframDomainsList, _ = types.ListValueFrom(ctx, types.StringType, iframeDomainsList)

	return
}

func isIFrameDomainsSet(iFrameDomains types.String) bool {
	if iFrameDomains.IsUnknown() {
		// The value is  unkown, so no value is available
		return false
	}

	if iFrameDomains.IsNull() {
		// The value is null, so no value is available
		return false
	}

	return true
}
