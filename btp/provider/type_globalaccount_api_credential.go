package provider

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_api"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type globalaccountApiCredentialType struct {
	GlobalaccountId   types.String `tfsdk:"globalaccount_id"`
	Name              types.String `tfsdk:"name"`
	ClientId          types.String `tfsdk:"client_id"`
	CredentialType    types.String `tfsdk:"credential_type"`
	ClientSecret      types.String `tfsdk:"client_secret"`
	CertificatePassed types.String `tfsdk:"certificate_passed"`
	Certificate       types.String `tfsdk:"certificate_received"`
	Key               types.String `tfsdk:"key"`
	ReadOnly          types.Bool   `tfsdk:"read_only"`
	TokenUrl          types.String `tfsdk:"token_url"`
	ApiUrl            types.String `tfsdk:"api_url"`
}

func globalaccountApiCredentialFromValue(_ context.Context, cliRes xsuaa_api.ApiCredential) (globalaccountApiCredentialType, diag.Diagnostics) {

	res := globalaccountApiCredentialType{
		GlobalaccountId: types.StringValue(cliRes.SubaccountId),
		Name:            types.StringValue(cliRes.Name),
		ClientId:        types.StringValue(cliRes.ClientId),
		CredentialType:  types.StringValue(cliRes.CredentialType),
		ReadOnly:        types.BoolValue(cliRes.ReadOnly),
		TokenUrl:        types.StringValue(cliRes.TokenUrl),
		ApiUrl:          types.StringValue(cliRes.ApiUrl),
	}

	if len(cliRes.ClientSecret) > 0 {
		res.ClientSecret = types.StringValue(cliRes.ClientSecret)
	} else if len(cliRes.Certificate) > 0 {
		res.Certificate = types.StringValue(cliRes.Certificate)
		res.Key = types.StringValue(cliRes.Key)
	}

	return res, diag.Diagnostics{}
}
