package provider

import (
	"context"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_trust"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type globalaccountTrustConfigurationType struct {
	IdentityProvider types.String `tfsdk:"identity_provider"`
	Domain           types.String `tfsdk:"domain"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Origin           types.String `tfsdk:"origin"`
	Id               types.String `tfsdk:"id"`
	Type             types.String `tfsdk:"type"`
	Protocol         types.String `tfsdk:"protocol"`
	Status           types.String `tfsdk:"status"`
	ReadOnly         types.Bool   `tfsdk:"read_only"`
}

func globalaccountTrustConfigurationFromValue(ctx context.Context, value xsuaa_trust.TrustConfigurationResponseObject) (globalaccountTrustConfigurationType, diag.Diagnostics) {
	return globalaccountTrustConfigurationType{
		IdentityProvider: types.StringValue(value.IdentityProvider),
		Domain:           types.StringValue(value.Domain),
		Name:             types.StringValue(value.Name),
		Description:      types.StringValue(value.Description),
		Origin:           types.StringValue(value.OriginKey),
		Id:               types.StringValue(value.OriginKey),
		Type:             types.StringValue(value.TypeOfTrust),
		Protocol:         types.StringValue(value.Protocol),
		Status:           types.StringValue(value.Status),
		ReadOnly:         types.BoolValue(value.ReadOnly),
	}, diag.Diagnostics{}
}
