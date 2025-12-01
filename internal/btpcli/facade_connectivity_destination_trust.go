package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/connectivity"
)

type connectivityDestinationTrustFacade struct {
	cliClient *v2Client
}

func newConnectivityDestinationTrustFacade(cliClient *v2Client) connectivityDestinationTrustFacade {
	return connectivityDestinationTrustFacade{cliClient: cliClient}
}

func (f *connectivityDestinationTrustFacade) getCommand() string {
	return "connectivity/destination-trust"
}

func (f *connectivityDestinationTrustFacade) GetBySubaccount(ctx context.Context, subaccountID string, trustType bool) (connectivity.DestinationTrust, CommandResponse, error) {
	passive := "true"
	if trustType {
		passive = "false"
	}
	return doExecute[connectivity.DestinationTrust](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountID,
		"passive":    passive,
	}))
}
