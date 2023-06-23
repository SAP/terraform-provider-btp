package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
)

func newServicesOfferingFacade(cliClient *v2Client) servicesOfferingFacade {
	return servicesOfferingFacade{cliClient: cliClient}
}

type servicesOfferingFacade struct {
	cliClient *v2Client
}

func (f servicesOfferingFacade) getCommand() string {
	return "services/offering"
}

func (f servicesOfferingFacade) List(ctx context.Context, subaccountId string, fieldsFilter string, labelsFilter string, environment string) ([]servicemanager.ServiceOfferingResponseObject, CommandResponse, error) {
	params := map[string]string{
		"subaccount": subaccountId,
	}

	if len(fieldsFilter) > 0 {
		params["fieldsFilter"] = fieldsFilter
	}

	if len(labelsFilter) > 0 {
		params["labelsFilter"] = labelsFilter
	}

	if len(environment) > 0 {
		params["environment"] = environment
	}

	return doExecute[[]servicemanager.ServiceOfferingResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), params))
}

func (f servicesOfferingFacade) GetById(ctx context.Context, subaccountId string, offeringId string) (servicemanager.ServiceOfferingResponseObject, CommandResponse, error) {
	return doExecute[servicemanager.ServiceOfferingResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"id":         offeringId,
	}))
}

func (f servicesOfferingFacade) GetByName(ctx context.Context, subaccountId string, offeringName string) (servicemanager.ServiceOfferingResponseObject, CommandResponse, error) {
	return doExecute[servicemanager.ServiceOfferingResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"name":       offeringName,
	}))
}
