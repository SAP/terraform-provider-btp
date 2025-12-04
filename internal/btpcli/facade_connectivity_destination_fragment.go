package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/connectivity"
)

type connectivityDestinationFragmentFacade struct {
	cliClient *v2Client
}

func newConnectivityDestinationFragmentFacade(cliClient *v2Client) connectivityDestinationFragmentFacade {
	return connectivityDestinationFragmentFacade{cliClient: cliClient}
}

func (f *connectivityDestinationFragmentFacade) getCommand() string {
	return "connectivity/destination-fragment"
}

func (f *connectivityDestinationFragmentFacade) GetBySubaccount(ctx context.Context, subaccountID string, name string) (connectivity.DestinationFragment, CommandResponse, error) {
	return doExecute[connectivity.DestinationFragment](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"name":       name,
		"subaccount": subaccountID,
	}))
}

func (f *connectivityDestinationFragmentFacade) GetByServiceInstance(ctx context.Context, subaccountID string, name string, serviceInstanceID string) (connectivity.DestinationFragment, CommandResponse, error) {
	return doExecute[connectivity.DestinationFragment](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"name":            name,
		"subaccount":      subaccountID,
		"serviceInstance": serviceInstanceID,
	}))
}

func (f *connectivityDestinationFragmentFacade) ListBySubaccount(ctx context.Context, subaccountID string) ([]connectivity.DestinationFragment, CommandResponse, error) {
	return doExecute[[]connectivity.DestinationFragment](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountID,
	}))
}

func (f *connectivityDestinationFragmentFacade) ListByServiceInstance(ctx context.Context, subaccountID string, serviceInstanceID string) ([]connectivity.DestinationFragment, CommandResponse, error) {
	return doExecute[[]connectivity.DestinationFragment](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount":      subaccountID,
		"serviceInstance": serviceInstanceID,
	}))
}
