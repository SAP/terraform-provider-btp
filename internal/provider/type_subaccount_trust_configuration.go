package provider

import (
	"context"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_trust"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

type subaccountTrustConfigurationType struct {
	SubaccountId          types.String `tfsdk:"subaccount_id"`
	IdentityProvider      types.String `tfsdk:"identity_provider"`
	Domain                types.String `tfsdk:"domain"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	LinkText              types.String `tfsdk:"link_text"`
	AvailableForUserLogon types.Bool   `tfsdk:"available_for_user_logon"`
	AutoCreateShadowUsers types.Bool   `tfsdk:"auto_create_shadow_users"`
	Origin                types.String `tfsdk:"origin"`
	Id                    types.String `tfsdk:"id"`
	Type                  types.String `tfsdk:"type"`
	Protocol              types.String `tfsdk:"protocol"`
	Status                types.String `tfsdk:"status"`
	ReadOnly              types.Bool   `tfsdk:"read_only"`
}

func subaccountTrustConfigurationFromValue(ctx context.Context, value xsuaa_trust.TrustConfigurationResponseObject) (subaccountTrustConfigurationType, diag.Diagnostics) {
	availableForUserLogon, _ := strconv.ParseBool(value.AvailableForUserLogon)
	autoCreateShadowUsers, _ := strconv.ParseBool(value.CreateShadowUsersDuringLogon)
	domain := types.StringNull()
	if len(value.Domain) > 0 {
		domain = types.StringValue(value.Domain)
	}
	return subaccountTrustConfigurationType{
		SubaccountId:          types.StringNull(),
		IdentityProvider:      types.StringValue(value.IdentityProvider),
		Domain:                domain,
		Name:                  types.StringValue(value.Name),
		Description:           types.StringValue(value.Description),
		LinkText:              types.StringValue(value.LinkTextForUserLogon),
		AvailableForUserLogon: types.BoolValue(availableForUserLogon),
		AutoCreateShadowUsers: types.BoolValue(autoCreateShadowUsers),
		Origin:                types.StringValue(value.OriginKey),
		Id:                    types.StringValue(value.OriginKey),
		Type:                  types.StringValue(value.TypeOfTrust),
		Protocol:              types.StringValue(value.Protocol),
		Status:                types.StringValue(value.Status),
		ReadOnly:              types.BoolValue(value.ReadOnly),
	}, diag.Diagnostics{}
}
