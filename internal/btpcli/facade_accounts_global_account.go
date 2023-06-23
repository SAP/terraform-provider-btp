package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
)

func newAccountsGlobalAccountFacade(cliClient *v2Client) accountsGlobalAccountFacade {
	return accountsGlobalAccountFacade{cliClient: cliClient}
}

type accountsGlobalAccountFacade struct {
	cliClient *v2Client
}

func (f *accountsGlobalAccountFacade) getCommand() string {
	return "accounts/global-account"
}

func (f *accountsGlobalAccountFacade) Get(ctx context.Context) (cis.GlobalAccountResponseObject, CommandResponse, error) {
	return doExecute[cis.GlobalAccountResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}
