package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
)

func newSecurityAppFacade(cliClient *v2Client) securityAppFacade {
	return securityAppFacade{cliClient: cliClient}
}

type securityAppFacade struct {
	cliClient *v2Client
}

func (f *securityAppFacade) getCommand() string {
	return "security/app"
}

func (f *securityAppFacade) ListByGlobalAccount(ctx context.Context) ([]xsuaa_authz.App, *CommandResponse, error) {
	return doExecute[[]xsuaa_authz.App](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *securityAppFacade) GetByGlobalAccount(ctx context.Context, appId string) (xsuaa_authz.App, *CommandResponse, error) {
	return doExecute[xsuaa_authz.App](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"appId":         appId,
	}))
}

func (f *securityAppFacade) ListBySubaccount(ctx context.Context, subaccountId string) ([]xsuaa_authz.App, *CommandResponse, error) {
	return doExecute[[]xsuaa_authz.App](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

func (f *securityAppFacade) GetBySubaccount(ctx context.Context, subaccountId string, appId string) (xsuaa_authz.App, *CommandResponse, error) {
	return doExecute[xsuaa_authz.App](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"appId":      appId,
	}))
}

func (f *securityAppFacade) ListByDirectory(ctx context.Context, directoryId string) ([]xsuaa_authz.App, *CommandResponse, error) {
	return doExecute[[]xsuaa_authz.App](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"directory": directoryId,
	}))
}

func (f *securityAppFacade) GetByDirectory(ctx context.Context, directoryId string, appId string) (xsuaa_authz.App, *CommandResponse, error) {
	return doExecute[xsuaa_authz.App](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"directory": directoryId,
		"appId":     appId,
	}))
}
