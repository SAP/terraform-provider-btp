package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
)

func newSecurityUserFacade(cliClient *v2Client) securityUserFacade {
	return securityUserFacade{cliClient: cliClient}
}

type securityUserFacade struct {
	cliClient *v2Client
}

func (f *securityUserFacade) getCommand() string {
	return "security/user"
}

func (f *securityUserFacade) ListByGlobalAccount(ctx context.Context, origin string) ([]string, *CommandResponse, error) {
	return doExecute[[]string](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"origin":        origin,
	}))
}

func (f *securityUserFacade) GetByGlobalAccount(ctx context.Context, username string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"userName":      username,
		"origin":        origin,
	}))
}

func (f *securityUserFacade) ListBySubaccount(ctx context.Context, subaccountId string, origin string) ([]string, *CommandResponse, error) {
	return doExecute[[]string](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"origin":     origin,
	}))
}

func (f *securityUserFacade) GetBySubaccount(ctx context.Context, subaccountId string, username string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"userName":   username,
		"origin":     origin,
	}))
}

func (f *securityUserFacade) ListByDirectory(ctx context.Context, directoryId string, origin string) ([]string, *CommandResponse, error) {
	return doExecute[[]string](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"directory": directoryId,
		"origin":    origin,
	}))
}

func (f *securityUserFacade) GetByDirectory(ctx context.Context, directoryId string, username string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"directory": directoryId,
		"userName":  username,
		"origin":    origin,
	}))
}
