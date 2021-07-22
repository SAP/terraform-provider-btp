package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
)

func newAccountsAvailableRegionFacade(cliClient *v2Client) accountsAvailableRegionFacade {
	return accountsAvailableRegionFacade{cliClient: cliClient}
}

type accountsAvailableRegionFacade struct {
	cliClient *v2Client
}

func (f *accountsAvailableRegionFacade) getCommand() string {
	return "accounts/available-region"
}

func (f *accountsAvailableRegionFacade) List(ctx context.Context) (cis.DataCenterResponseCollection, *CommandResponse, error) {
	return doExecute[cis.DataCenterResponseCollection](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}
