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

func (f *connectivityDestinationFacade) ListBySubaccount(ctx context.Context, subaccountId string, serviceInstance string) ([]connectivity.DestinationResponse, CommandResponse, error) {
	args := map[string]any{
		"subaccount": subaccountId,
		"namesOnly":  false,
	}

	if len(serviceInstance) > 0 {
		args["serviceInstance"] = serviceInstance
	}
	return doExecute[[]connectivity.DestinationResponse](f.cliClient, ctx, NewListRequest(f.getCommand(), args))
}

func (f *connectivityDestinationFacade) ListNamesBySubaccount(ctx context.Context, subaccountId string, serviceInstance string) ([]map[string]string, CommandResponse, error) {
	args := map[string]any{
		"subaccount": subaccountId,
		"namesOnly":  true,
	}

	if len(serviceInstance) > 0 {
		args["serviceInstance"] = serviceInstance
	}
	return doExecute[[]map[string]string](f.cliClient, ctx, NewListRequest(f.getCommand(), args))
}
func (f *connectivityDestinationFacade) CreateBySubaccount(ctx context.Context, subaccountId string, jsonarg string, serviceInstance string) (connectivity.DestinationResponse, CommandResponse, error) {
	args := map[string]string{
		"destinationConfiguration": jsonarg,
		"subaccount":               subaccountId,
	}

	if len(serviceInstance) > 0 {
		args["serviceInstance"] = serviceInstance
	}
	return doExecute[connectivity.DestinationResponse](f.cliClient, ctx, NewCreateRequest(f.getCommand(), args))
}
func (f *connectivityDestinationFacade) UpdateBySubaccount(ctx context.Context, subaccountId string, jsonarg string, serviceInstance string) (connectivity.DestinationResponse, CommandResponse, error) {
	args := map[string]string{
		"destinationConfiguration": jsonarg,
		"subaccount":               subaccountId,
	}

	if len(serviceInstance) > 0 {
		args["serviceInstance"] = serviceInstance
	}
	return doExecute[connectivity.DestinationResponse](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), args))
}

func (f *connectivityDestinationFacade) DeleteBySubaccount(ctx context.Context, subaccountId string, name string, serviceInstance string) (connectivity.DestinationResponse, CommandResponse, error) {
	args := map[string]string{
		"name":       name,
		"subaccount": subaccountId,
	}

	if len(serviceInstance) > 0 {
		args["serviceInstance"] = serviceInstance
	}
	return doExecute[connectivity.DestinationResponse](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), args))
}
