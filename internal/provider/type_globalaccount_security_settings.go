package provider

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_settings"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type globalaccountSecuritySettingsType struct {
	CustomEmailDomains                types.Set    `tfsdk:"custom_email_domains"`
	DefaultIdentityProvider           types.String `tfsdk:"default_identity_provider"`
	TreatUsersWithSameEmailAsSameUser types.Bool   `tfsdk:"treat_users_with_same_email_as_same_user"`
	AccessTokenValidity               types.Int64  `tfsdk:"access_token_validity"`
	RefreshTokenValidity              types.Int64  `tfsdk:"refresh_token_validity"`
}

func globalaccountSecuritySettingsFromValue(ctx context.Context, value xsuaa_settings.TenantSettingsResp) (tenantSettings globalaccountSecuritySettingsType, diags diag.Diagnostics) {
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

	return
}
