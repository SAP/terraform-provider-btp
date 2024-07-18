package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
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

func (f servicesBrokerFacade) List(ctx context.Context, subaccountId string, fieldsFilter string, labelsFilter string) ([]servicemanager.ServiceBrokerResponseObject, CommandResponse, error) {
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

func (f servicesBrokerFacade) GetById(ctx context.Context, subaccountId string, brokerId string) (servicemanager.ServiceBrokerResponseObject, CommandResponse, error) {
	return doExecute[servicemanager.ServiceBrokerResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"id":         brokerId,
	}))
}

func (f servicesBrokerFacade) GetByName(ctx context.Context, subaccountId string, brokerName string) (servicemanager.ServiceBrokerResponseObject, CommandResponse, error) {
	return doExecute[servicemanager.ServiceBrokerResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"name":       brokerName,
	}))
}

type SubaccountServiceBrokerRegisterInput struct {
	Subaccount  string              `btpcli:"subaccount"`
	Name        string              `btpcli:"name"`
	Description string              `btpcli:"description"`
	User        string              `btpcli:"user"`
	Password    string              `btpcli:"password"`
	URL         string              `btpcli:"url"`
	Labels      map[string][]string `btpcli:"labels"`
}

func (f servicesBrokerFacade) Register(ctx context.Context, args SubaccountServiceBrokerRegisterInput) (servicemanager.ServiceBrokerResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return servicemanager.ServiceBrokerResponseObject{}, CommandResponse{}, err
	}

	return doExecute[servicemanager.ServiceBrokerResponseObject](f.cliClient, ctx, NewRegisterRequest(f.getCommand(), params))
}

type SubaccountServiceBrokerUpdateInput struct {
	Id          string              `btpcli:"id"`
	Subaccount  string              `btpcli:"subaccount"`
	NewName     string              `btpcli:"newName"`
	Description string              `btpcli:"description"`
	User        string              `btpcli:"user"`
	Password    string              `btpcli:"password"`
	URL         string              `btpcli:"url"`
	Labels      map[string][]string `btpcli:"labels"`
}

func (f servicesBrokerFacade) Update(ctx context.Context, args SubaccountServiceBrokerUpdateInput) (servicemanager.ServiceBrokerResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return servicemanager.ServiceBrokerResponseObject{}, CommandResponse{}, err
	}

	return doExecute[servicemanager.ServiceBrokerResponseObject](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
}

// TODO Update
func (f servicesBrokerFacade) Unregister(ctx context.Context, subaccountId string, serviceId string) (CommandResponse, error) {
	res, err := f.cliClient.Execute(ctx, NewUnregisterRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"id":         serviceId,
		"confirm":    "true",
	}))
	return res, err
}
