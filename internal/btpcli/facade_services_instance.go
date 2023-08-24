package btpcli

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newServicesInstanceFacade(cliClient *v2Client) servicesInstanceFacade {
	return servicesInstanceFacade{cliClient: cliClient}
}

type servicesInstanceFacade struct {
	cliClient *v2Client
}

func (f servicesInstanceFacade) getCommand() string {
	return "services/instance"
}

func (f servicesInstanceFacade) List(ctx context.Context, subaccountId string, fieldsFilter string, labelsFilter string) ([]servicemanager.ServiceInstanceResponseObject, CommandResponse, error) {
	params := map[string]string{
		"subaccount": subaccountId,
	}

	if len(fieldsFilter) > 0 {
		params["fieldsFilter"] = fieldsFilter
	}

	if len(labelsFilter) > 0 {
		params["labelsFilter"] = labelsFilter
	}

	return doExecute[[]servicemanager.ServiceInstanceResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), params))
}

func (f servicesInstanceFacade) GetById(ctx context.Context, subaccountId string, instanceId string) (servicemanager.ServiceInstanceResponseObject, CommandResponse, error) {
	return doExecute[servicemanager.ServiceInstanceResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"id":         instanceId,
		"parameters": "false",
	}))
}

func (f servicesInstanceFacade) GetByName(ctx context.Context, subaccountId string, instanceName string) (servicemanager.ServiceInstanceResponseObject, CommandResponse, error) {
	return doExecute[servicemanager.ServiceInstanceResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"name":       instanceName,
		"parameters": "false",
	}))
}

type ServiceInstanceCreateInput struct {
	Name          string              `btpcli:"name"`
	Subaccount    string              `btpcli:"subaccount"`
	ServicePlanId string              `btpcli:"plan"`
	Parameters    *string             `btpcli:"parameters"`
	Labels        map[string][]string `btpcli:"labels"`
}

func (f servicesInstanceFacade) Create(ctx context.Context, args *ServiceInstanceCreateInput) (servicemanager.ServiceInstanceResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return servicemanager.ServiceInstanceResponseObject{}, CommandResponse{}, err
	}

	// return doExecute[servicemanager.ServiceInstanceResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
	serviceInstanceResponseObject, cmdRes, err := doExecute[servicemanager.ServiceInstanceResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))

	//TODO workaround for NGPBUG-350117 => needs to be structured after fix
	if cmdRes.StatusCode == 201 {
		return serviceInstanceResponseObject, cmdRes, err
	} else if cmdRes.StatusCode == 202 {
		return f.GetByName(ctx, args.Subaccount, args.Name)
	} else if err != nil {
		return servicemanager.ServiceInstanceResponseObject{}, cmdRes, err
	} else {
		err = fmt.Errorf("the backend responded with an unknown error: %d", cmdRes.StatusCode)
		return servicemanager.ServiceInstanceResponseObject{}, cmdRes, err
	}
}

type ServiceInstanceUpdateInput struct {
	Id            string              `btpcli:"id"`
	Name          string              `btpcli:"name"`
	NewName       string              `btpcli:"newName"`
	Subaccount    string              `btpcli:"subaccount"`
	ServicePlanId string              `btpcli:"plan"`
	Parameters    *string             `btpcli:"parameters"`
	Labels        map[string][]string `btpcli:"labels"`
}

func (f servicesInstanceFacade) Update(ctx context.Context, args *ServiceInstanceUpdateInput) (servicemanager.ServiceInstanceResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return servicemanager.ServiceInstanceResponseObject{}, CommandResponse{}, err
	}

	//TODO workaround for NGPBUG-359662 and NGPBUG-350117 => needs to be rebuilt after fix
	//return doExecute[servicemanager.ServiceInstanceResponseObject](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
	// 1. Call the update directly without deserialize the response
	// 2. Do a consequent GET request to get a consistent response of the instance.

	res, err := f.cliClient.Execute(ctx, NewUpdateRequest(f.getCommand(), params))

	if err != nil {
		return servicemanager.ServiceInstanceResponseObject{}, res, err
	}

	if res.StatusCode == 202 {
		return f.GetById(ctx, args.Subaccount, args.Id)
	} else {
		err = fmt.Errorf("the backend responded with an unknown error: %d", res.StatusCode)
		return servicemanager.ServiceInstanceResponseObject{}, res, err
	}

}

func (f servicesInstanceFacade) Delete(ctx context.Context, subaccountId string, serviceId string) (CommandResponse, error) {
	res, err := f.cliClient.Execute(ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"id":         serviceId,
		"confirm":    "true",
	}))
	return res, err
}
