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

	serviceInstanceResponseObject, cmdRes, err := doExecute[servicemanager.ServiceInstanceResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))

	//Workaround for NGPBUG-350117 => fix not feasible, keeping workaround
	if cmdRes.StatusCode != 202 && err == nil {
		return serviceInstanceResponseObject, cmdRes, err
	} else if cmdRes.StatusCode == 202 && err == nil {
		return f.GetByName(ctx, args.Subaccount, args.Name)
	} else if err != nil {
		// Error case
		return servicemanager.ServiceInstanceResponseObject{}, cmdRes, err
	} else {
		// Fallback for unknown errors from service manager
		err = fmt.Errorf("the backend responded with an unknown error: %d", cmdRes.StatusCode)
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
