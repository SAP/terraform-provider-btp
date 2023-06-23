package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
)

func newAccountsLabelFacade(cliClient *v2Client) accountsLabelFacade {
	return accountsLabelFacade{cliClient: cliClient}
}

type accountsLabelFacade struct {
	cliClient *v2Client
}

func (f *accountsLabelFacade) getCommand() string {
	return "accounts/label"
}

func (f *accountsLabelFacade) ListBySubaccount(ctx context.Context, subaccountId string) (cis.LabelsResponseObject, CommandResponse, error) {
	return doExecute[cis.LabelsResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"subaccountID":  subaccountId,
	}))
}

func (f *accountsLabelFacade) ListByDirectory(ctx context.Context, directoryId string) (cis.LabelsResponseObject, CommandResponse, error) {
	return doExecute[cis.LabelsResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"directoryID":   directoryId,
	}))
}
