package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/provisioning"
)

func newAccountsAvailableEnvironmentFacade(cliClient *v2Client) accountsAvailableEnvironmentFacade {
	return accountsAvailableEnvironmentFacade{cliClient: cliClient}
}

type accountsAvailableEnvironmentFacade struct {
	cliClient *v2Client
}

func (f *accountsAvailableEnvironmentFacade) getCommand() string {
	return "accounts/available-environment"
}

func (f *accountsAvailableEnvironmentFacade) List(ctx context.Context, subaccountId string) (provisioning.AvailableEnvironmentResponseCollection, CommandResponse, error) {
	return doExecute[provisioning.AvailableEnvironmentResponseCollection](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}
