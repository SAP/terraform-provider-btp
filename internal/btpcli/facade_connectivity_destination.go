package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/connectivity"
)

func newConnectivityDestinationFacade(cliClient *v2Client) connectivityDestinationFacade {
	return connectivityDestinationFacade{cliClient: cliClient}
}

type connectivityDestinationFacade struct {
	cliClient *v2Client
}

func (f *connectivityDestinationFacade) getCommand() string {
	return "connectivity/destination"
}

func (f *connectivityDestinationFacade) GetBySubaccount(ctx context.Context, subaccountId string, name string, serviceInstance string) (connectivity.DestinationResponse, CommandResponse, error) {
	args := map[string]string{
		"name":       name,
		"subaccount": subaccountId,
	}

	if len(serviceInstance) > 0 {
		args["serviceInstance"] = serviceInstance
	}
	return doExecute[connectivity.DestinationResponse](f.cliClient, ctx, NewGetRequest(f.getCommand(), args))
}
