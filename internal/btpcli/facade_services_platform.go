package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
)

func newServicesPlatformFacade(cliClient *v2Client) servicesPlatformFacade {
	return servicesPlatformFacade{cliClient: cliClient}
}

type servicesPlatformFacade struct {
	cliClient *v2Client
}

func (f servicesPlatformFacade) getCommand() string {
	return "services/platform"
}

func (f servicesPlatformFacade) List(ctx context.Context, subaccountId string, fieldsFilter string, labelsFilter string) ([]servicemanager.PlatformResponseObject, *CommandResponse, error) {
	params := map[string]string{
		"subaccount": subaccountId,
	}

	if len(fieldsFilter) > 0 {
		params["fieldsFilter"] = fieldsFilter
	}

	if len(labelsFilter) > 0 {
		params["labelsFilter"] = labelsFilter
	}

	return doExecute[[]servicemanager.PlatformResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), params))
}

func (f servicesPlatformFacade) GetById(ctx context.Context, subaccountId string, platformId string) (servicemanager.PlatformResponseObject, *CommandResponse, error) {
	return doExecute[servicemanager.PlatformResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"id":         platformId,
	}))
}

func (f servicesPlatformFacade) GetByName(ctx context.Context, subaccountId string, platformName string) (servicemanager.PlatformResponseObject, *CommandResponse, error) {
	return doExecute[servicemanager.PlatformResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"name":       platformName,
	}))
}
