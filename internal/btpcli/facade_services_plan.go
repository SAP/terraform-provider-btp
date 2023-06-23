package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
)

func newServicesPlanFacade(cliClient *v2Client) servicesPlanFacade {
	return servicesPlanFacade{cliClient: cliClient}
}

type servicesPlanFacade struct {
	cliClient *v2Client
}

func (f servicesPlanFacade) getCommand() string {
	return "services/plan"
}

func (f servicesPlanFacade) List(ctx context.Context, subaccountId string, fieldsFilter string, labelsFilter string, environment string) ([]servicemanager.ServicePlanResponseObject, CommandResponse, error) {
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

	return doExecute[[]servicemanager.ServicePlanResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), params))
}

func (f servicesPlanFacade) GetById(ctx context.Context, subaccountId string, planId string) (servicemanager.ServicePlanResponseObject, CommandResponse, error) {
	return doExecute[servicemanager.ServicePlanResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"id":         planId,
	}))
}

func (f servicesPlanFacade) GetByName(ctx context.Context, subaccountId string, planName string, offeringName string) (servicemanager.ServicePlanResponseObject, CommandResponse, error) {
	return doExecute[servicemanager.ServicePlanResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount":   subaccountId,
		"name":         planName,
		"offeringName": offeringName,
	}))
}
