package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newServicesBindingFacade(cliClient *v2Client) servicesBindingFacade {
	return servicesBindingFacade{cliClient: cliClient}
}

type servicesBindingFacade struct {
	cliClient *v2Client
}

func (f servicesBindingFacade) getCommand() string {
	return "services/binding"
}

func (f servicesBindingFacade) List(ctx context.Context, subaccountId string, fieldsFilter string, labelsFilter string) ([]servicemanager.ServiceBindingResponseObject, *CommandResponse, error) {
	params := map[string]string{
		"subaccount": subaccountId,
	}

	if len(fieldsFilter) > 0 {
		params["fieldsFilter"] = fieldsFilter
	}

	if len(labelsFilter) > 0 {
		params["labelsFilter"] = labelsFilter
	}

	return doExecute[[]servicemanager.ServiceBindingResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), params))
}

func (f servicesBindingFacade) GetById(ctx context.Context, subaccountId string, bindingId string) (servicemanager.ServiceBindingResponseObject, *CommandResponse, error) {
	return doExecute[servicemanager.ServiceBindingResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"id":         bindingId,
	}))
}

func (f servicesBindingFacade) GetByName(ctx context.Context, subaccountId string, bindingName string) (servicemanager.ServiceBindingResponseObject, *CommandResponse, error) {
	return doExecute[servicemanager.ServiceBindingResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"name":       bindingName,
	}))
}

type SubaccountServiceBindingCreateInput struct {
	Subaccount        string  `btpcli:"subaccount"`
	ServiceInstanceId string  `btpcli:"serviceInstanceID"`
	Name              string  `btpcli:"name"`
	Parameters        string  `btpcli:"parameters"`
	Labels            *string `btpcli:"labels"`
}

func (f servicesBindingFacade) Create(ctx context.Context, args SubaccountServiceBindingCreateInput) (servicemanager.ServiceBindingResponseObject, *CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return servicemanager.ServiceBindingResponseObject{}, nil, err
	}

	return doExecute[servicemanager.ServiceBindingResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f servicesBindingFacade) Delete(ctx context.Context, subaccountId string, bindingId string) (servicemanager.ServiceBindingResponseObject, *CommandResponse, error) {
	return doExecute[servicemanager.ServiceBindingResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"id":         bindingId,
	}))
}
