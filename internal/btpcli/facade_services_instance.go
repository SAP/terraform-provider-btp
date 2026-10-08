package btpcli

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

const labelRemoveOp = "remove"
const labelAddOp = "add"

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
	return f.doGet(ctx, map[string]string{
		"subaccount": subaccountId,
		"id":         instanceId,
	})
}

func (f servicesInstanceFacade) GetByName(ctx context.Context, subaccountId string, instanceName string) (servicemanager.ServiceInstanceResponseObject, CommandResponse, error) {
	return f.doGet(ctx, map[string]string{
		"subaccount": subaccountId,
		"name":       instanceName,
	})
}

func (f servicesInstanceFacade) doGet(ctx context.Context, params map[string]string) (sir servicemanager.ServiceInstanceResponseObject, cr CommandResponse, err error) {

	// Execute a call for the instance without parameters
	params["parameters"] = "false"
	sir, cr, err = doExecute[servicemanager.ServiceInstanceResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), params))

	if err != nil {
		return
	}

	// Execute a call for the parameters. We need two calls because the parameters are not returned by the first call.
	params["parameters"] = "true"

	// In addition the response format might differ depending on the service instance.
	resData, _, err_param := doExecute[servicemanager.ServiceInstanceParametersData](f.cliClient, ctx, NewGetRequest(f.getCommand(), params))

	// Case 1 - Parameters are returned as data object
	if err_param == nil && len(resData.Parameters) != 0 {
		jsonString, _ := json.Marshal(resData.Parameters)
		sir.Parameters = string(jsonString)
		return
	}

	resPlain, _, err_param := doExecute[servicemanager.ServiceInstanceParametersPlain](f.cliClient, ctx, NewGetRequest(f.getCommand(), params))

	// Case 2 - Parameters are returned as plain object
	if err_param == nil && len(resPlain.Parameters) != 0 {
		jsonString, _ := json.Marshal(resPlain.Parameters)
		sir.Parameters = string(jsonString)
		return
	}

	// Even if the service instance has parameters, the parameters are not returned by the API due to settings in the service offering
	// The service offering must have the following setting:  'instances_retrievable: TRUE'
	// In this case we return the base service instance response object without parameters
	return
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

	serviceInstanceResponseObject, cmdRes, err := doExecute[servicemanager.ServiceInstanceResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))

	//Workaround for NGPBUG-350117 => fix not feasible, keeping workaround
	if cmdRes.StatusCode != 202 && err == nil {
		return serviceInstanceResponseObject, cmdRes, err
	} else if cmdRes.StatusCode == 202 && err == nil {
		return f.GetById(ctx, args.Subaccount, serviceInstanceResponseObject.Id)
	} else {
		// Error case as default
		return servicemanager.ServiceInstanceResponseObject{}, cmdRes, err
	}
}

type ServiceInstanceUpdateInput struct {
	Id            string  `btpcli:"id"`
	NewName       string  `btpcli:"newName"`
	Subaccount    string  `btpcli:"subaccount"`
	ServicePlanId string  `btpcli:"plan"`
	Parameters    *string `btpcli:"parameters"`
	LabelsPlan    map[string][]string
	LabelsState   map[string][]string
}

func (f servicesInstanceFacade) Update(ctx context.Context, args *ServiceInstanceUpdateInput) (servicemanager.ServiceInstanceResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return servicemanager.ServiceInstanceResponseObject{}, CommandResponse{}, err
	}

	computedLabels := computeLabelParam(args.LabelsPlan, args.LabelsState)

	if computedLabels != "" {
		// Parameter must only be added to call if non-empty
		params["labels"] = computedLabels
	}

	// The CLI server returns a BTP CLI specific response - no return of doExecute possible
	// Solution:
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

type ServiceInstanceShareInput struct {
	Id         string `btpcli:"id"`
	Subaccount string `btpcli:"subaccount"`
	Name       string `btpcli:"name"`
}

func (f servicesInstanceFacade) Share(ctx context.Context, args *ServiceInstanceShareInput) (servicemanager.ServiceInstanceResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return servicemanager.ServiceInstanceResponseObject{}, CommandResponse{}, err
	}

	// The CLI server returns a BTP CLI specific response - no return of doExecute possible
	// Solution:
	// 1. Call the update directly without deserialize the response
	// 2. Do a consequent GET request to get a consistent response of the instance.

	res, err := f.cliClient.Execute(ctx, NewShareRequest(f.getCommand(), params))

	if err != nil {
		return servicemanager.ServiceInstanceResponseObject{}, res, err
	}

	if res.StatusCode == 200 {
		return f.GetById(ctx, args.Subaccount, args.Id)
	} else {
		err = fmt.Errorf("the backend responded with an unknown error: %d", res.StatusCode)
		return servicemanager.ServiceInstanceResponseObject{}, res, err
	}
}

func (f servicesInstanceFacade) Unshare(ctx context.Context, args *ServiceInstanceShareInput) (servicemanager.ServiceInstanceResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return servicemanager.ServiceInstanceResponseObject{}, CommandResponse{}, err
	}

	// The CLI server returns a BTP CLI specific response - no return of doExecute possible
	// Solution:
	// 1. Call the update directly without deserialize the response
	// 2. Do a consequent GET request to get a consistent response of the instance.

	res, err := f.cliClient.Execute(ctx, NewUnshareRequest(f.getCommand(), params))

	if err != nil {
		return servicemanager.ServiceInstanceResponseObject{}, res, err
	}

	if res.StatusCode == 200 {
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

func computeLabelParam(labelsPlan map[string][]string, labelsState map[string][]string) string {

	var labelEntry servicemanager.Label
	var labelDiff []servicemanager.Label

	for k, v := range labelsState {
		if _, ok := labelsPlan[k]; !ok {
			// Do not add not found entries
			continue
		}
		if !reflect.DeepEqual(v, labelsPlan[k]) {
			// Old label needs to be removed
			labelEntry.Op = labelRemoveOp
			labelEntry.Key = k
			labelEntry.Values = v

			labelDiff = append(labelDiff, labelEntry)

			//New label needs to be added
			labelEntry.Op = labelAddOp
			labelEntry.Key = k
			labelEntry.Values = labelsPlan[k]

			labelDiff = append(labelDiff, labelEntry)

		}
	}

	for k, v := range labelsPlan {
		if _, ok := labelsState[k]; !ok {
			// Key was added, so it needs to be put into the "add" operation
			labelEntry.Op = labelAddOp
			labelEntry.Key = k
			labelEntry.Values = v

			labelDiff = append(labelDiff, labelEntry)
		}
	}

	for k, v := range labelsState {
		if _, ok := labelsPlan[k]; !ok {
			// Key was removed, so it needs to be put into the "remove" operation
			labelEntry.Op = labelRemoveOp
			labelEntry.Key = k
			labelEntry.Values = v

			labelDiff = append(labelDiff, labelEntry)
		}
	}

	if labelDiff != nil {
		jsonLabels, _ := json.Marshal(labelDiff)
		return string(jsonLabels)
	} else {
		return ""
	}
}
