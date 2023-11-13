package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_settings"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newSecuritySettingsFacade(cliClient *v2Client) securitySettingsFacade {
	return securitySettingsFacade{cliClient: cliClient}
}

type securitySettingsFacade struct {
	cliClient *v2Client
}

func (f *securitySettingsFacade) getCommand() string {
	return "security/settings"
}

func (f *securitySettingsFacade) ListByGlobalAccount(ctx context.Context) (xsuaa_settings.TenantSettingsResp, CommandResponse, error) {
	return doExecute[xsuaa_settings.TenantSettingsResp](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *securitySettingsFacade) ListBySubaccount(ctx context.Context, subaccountId string) (xsuaa_settings.TenantSettingsResp, CommandResponse, error) {
	return doExecute[xsuaa_settings.TenantSettingsResp](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

type SecuritySettingsUpdateInput struct {
	IFrame                            string `btpcli:"iFrameDomain"`
	CustomEmail                       string `btpcli:"customEmailDomains"`
	DefaultIDPForNonInteractiveLogon  string `btpcli:"defaultIdp"`
	TreatUsersWithSameEmailAsSameUser bool   `btpcli:"treatUsersWithSameEmailAsSameUser"`
	HomeRedirect                      string `btpcli:"homeRedirect"`
	AccessTokenValidity               int    `btpcli:"accessTokenValidity"`
	RefreshTokenValidity              int    `btpcli:"refreshTokenValidity"`
}

func (f *securitySettingsFacade) UpdateByGlobalAccount(ctx context.Context, args SecuritySettingsUpdateInput) (xsuaa_settings.TenantSettingsResp, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_settings.TenantSettingsResp{}, CommandResponse{}, err
	}

	params["globalAccount"] = f.cliClient.GetGlobalAccountSubdomain()

	return doExecute[xsuaa_settings.TenantSettingsResp](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
}

func (f *securitySettingsFacade) UpdateBySubaccount(ctx context.Context, subaccountId string, args SecuritySettingsUpdateInput) (xsuaa_settings.TenantSettingsResp, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_settings.TenantSettingsResp{}, CommandResponse{}, err
	}

	params["subaccount"] = subaccountId

	return doExecute[xsuaa_settings.TenantSettingsResp](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
}
