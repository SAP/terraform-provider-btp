package btpcli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/connectivity"
)

type connectivityDestinationFragmentFacade struct {
	cliClient *v2Client
}

const errMarshalFragmentContent = "failed to marshal destination fragment content: %w"

func newConnectivityDestinationFragmentFacade(cliClient *v2Client) connectivityDestinationFragmentFacade {
	return connectivityDestinationFragmentFacade{cliClient: cliClient}
}

func (f *connectivityDestinationFragmentFacade) getCommand() string {
	return "connectivity/destination-fragment"
}

func (f *connectivityDestinationFragmentFacade) GetBySubaccount(ctx context.Context, subaccountID string, name string) (connectivity.DestinationFragment, CommandResponse, error) {
	return doExecute[connectivity.DestinationFragment](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"name":       name,
		"subaccount": subaccountID,
	}))
}

func (f *connectivityDestinationFragmentFacade) GetByServiceInstance(ctx context.Context, subaccountID string, name string, serviceInstanceID string) (connectivity.DestinationFragment, CommandResponse, error) {
	return doExecute[connectivity.DestinationFragment](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"name":            name,
		"subaccount":      subaccountID,
		"serviceInstance": serviceInstanceID,
	}))
}

func (f *connectivityDestinationFragmentFacade) ListBySubaccount(ctx context.Context, subaccountID string) ([]connectivity.DestinationFragment, CommandResponse, error) {
	return doExecute[[]connectivity.DestinationFragment](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountID,
	}))
}

func (f *connectivityDestinationFragmentFacade) ListByServiceInstance(ctx context.Context, subaccountID string, serviceInstanceID string) ([]connectivity.DestinationFragment, CommandResponse, error) {
	return doExecute[[]connectivity.DestinationFragment](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount":      subaccountID,
		"serviceInstance": serviceInstanceID,
	}))
}

func (f *connectivityDestinationFragmentFacade) CreateBySubaccount(ctx context.Context, subaccountID string, content map[string]string) (connectivity.DestinationFragment, CommandResponse, error) {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		return connectivity.DestinationFragment{}, CommandResponse{}, fmt.Errorf(errMarshalFragmentContent, err)
	}

	params := map[string]string{
		"subaccount": subaccountID,
		"content":    string(jsonBytes),
	}

	return doExecute[connectivity.DestinationFragment](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *connectivityDestinationFragmentFacade) CreateByServiceInstance(ctx context.Context, subaccountID string, serviceInstanceID string, content map[string]string) (connectivity.DestinationFragment, CommandResponse, error) {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		return connectivity.DestinationFragment{}, CommandResponse{}, fmt.Errorf(errMarshalFragmentContent, err)
	}

	params := map[string]string{
		"subaccount":      subaccountID,
		"serviceInstance": serviceInstanceID,
		"content":         string(jsonBytes),
	}

	return doExecute[connectivity.DestinationFragment](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *connectivityDestinationFragmentFacade) UpdateBySubaccount(ctx context.Context, subaccountID string, content map[string]string) (connectivity.DestinationFragment, CommandResponse, error) {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		return connectivity.DestinationFragment{}, CommandResponse{}, fmt.Errorf(errMarshalFragmentContent, err)
	}

	params := map[string]string{
		"subaccount": subaccountID,
		"content":    string(jsonBytes),
	}

	return doExecute[connectivity.DestinationFragment](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
}

func (f *connectivityDestinationFragmentFacade) UpdateByServiceInstance(ctx context.Context, subaccountID string, serviceInstanceID string, content map[string]string) (connectivity.DestinationFragment, CommandResponse, error) {
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		return connectivity.DestinationFragment{}, CommandResponse{}, fmt.Errorf(errMarshalFragmentContent, err)
	}

	params := map[string]string{
		"subaccount":      subaccountID,
		"serviceInstance": serviceInstanceID,
		"content":         string(jsonBytes),
	}

	return doExecute[connectivity.DestinationFragment](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
}

func (f *connectivityDestinationFragmentFacade) DeleteBySubaccount(ctx context.Context, subaccountID string, name string) (connectivity.DestinationFragment, CommandResponse, error) {
	return doExecute[connectivity.DestinationFragment](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountID,
		"name":       name,
	}))
}

func (f *connectivityDestinationFragmentFacade) DeleteByServiceInstance(ctx context.Context, subaccountID string, name string, serviceInstanceID string) (connectivity.DestinationFragment, CommandResponse, error) {
	return doExecute[connectivity.DestinationFragment](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount":      subaccountID,
		"name":            name,
		"serviceInstance": serviceInstanceID,
	}))
}
