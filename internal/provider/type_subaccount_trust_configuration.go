package provider

import (
	"context"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_trust"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type subaccountTrustConfigurationType struct {
	SubaccountId     types.String `tfsdk:"subaccount_id"`
	Origin           types.String `tfsdk:"origin"`
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Type             types.String `tfsdk:"type"`
	IdentityProvider types.String `tfsdk:"identity_provider"`
	Protocol         types.String `tfsdk:"protocol"`
	Status           types.String `tfsdk:"status"`
	ReadOnly         types.Bool   `tfsdk:"read_only"`
}

func subaccountTrustConfigurationFromValue(ctx context.Context, value xsuaa_trust.TrustConfigurationResponseObject) (subaccountTrustConfigurationType, diag.Diagnostics) {
	return subaccountTrustConfigurationType{
		SubaccountId:     types.StringNull(),
		Origin:           types.StringValue(value.OriginKey),
		Id:               types.StringValue(value.OriginKey),
		Name:             types.StringValue(value.Name),
		Description:      types.StringValue(value.Description),
		Type:             types.StringValue(value.TypeOfTrust),
		IdentityProvider: types.StringValue(value.IdentityProvider),
		Protocol:         types.StringValue(value.Protocol),
		Status:           types.StringValue(value.Status),
		ReadOnly:         types.BoolValue(value.ReadOnly),
	}, diag.Diagnostics{}
}
