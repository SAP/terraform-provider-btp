package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
)

func newSecurityIdentityProviderFacade(cliClient *v2Client) securityIdentityProviderFacade {
	return securityIdentityProviderFacade{cliClient: cliClient}
}

type securityIdentityProviderFacade struct {
	cliClient *v2Client
}

func (f *securityIdentityProviderFacade) getCommand() string {
	return "security/available-idp"
}

func (f *securityIdentityProviderFacade) ListByGlobalAccount(ctx context.Context) ([]xsuaa_authz.Idp, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.Idp](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *securityIdentityProviderFacade) GetByGlobalAccount(ctx context.Context, host string) (xsuaa_authz.Idp, CommandResponse, error) {
	return doExecute[xsuaa_authz.Idp](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"iasTenantUrl":  host,
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *securityIdentityProviderFacade) ListBySubaccount(ctx context.Context, subaccountId string) ([]xsuaa_authz.Idp, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.Idp](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

func (f *securityIdentityProviderFacade) GetBySubaccount(ctx context.Context, subaccountId string, host string) (xsuaa_authz.Idp, CommandResponse, error) {
	return doExecute[xsuaa_authz.Idp](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount":   subaccountId,
		"iasTenantUrl": host,
	}))
}
