package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
)

func newServicesBrokerFacade(cliClient *v2Client) servicesBrokerFacade {
	return servicesBrokerFacade{cliClient: cliClient}
}

type servicesBrokerFacade struct {
	cliClient *v2Client
}

func (f servicesBrokerFacade) getCommand() string {
	return "services/broker"
}

func (f servicesBrokerFacade) List(ctx context.Context, subaccountId string, fieldsFilter string, labelsFilter string) ([]servicemanager.ServiceBrokerResponseObject, *CommandResponse, error) {
	params := map[string]string{
		"subaccount": subaccountId,
	}

	if len(fieldsFilter) > 0 {
		params["fieldsFilter"] = fieldsFilter
	}

	if len(labelsFilter) > 0 {
		params["labelsFilter"] = labelsFilter
	}

	return doExecute[[]servicemanager.ServiceBrokerResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), params))
}

func (f servicesBrokerFacade) GetById(ctx context.Context, subaccountId string, brokerId string) (servicemanager.ServiceBrokerResponseObject, *CommandResponse, error) {
	return doExecute[servicemanager.ServiceBrokerResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"id":         brokerId,
	}))
}

func (f servicesBrokerFacade) GetByName(ctx context.Context, subaccountId string, brokerName string) (servicemanager.ServiceBrokerResponseObject, *CommandResponse, error) {
	return doExecute[servicemanager.ServiceBrokerResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"name":       brokerName,
	}))
}
