package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/provisioning"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newAccountsEnvironmentInstanceFacade(cliClient *v2Client) accountsEnvironmentInstanceFacade {
	return accountsEnvironmentInstanceFacade{cliClient}
}

type accountsEnvironmentInstanceFacade struct {
	cliClient *v2Client
}

func (f *accountsEnvironmentInstanceFacade) getCommand() string {
	return "accounts/environment-instance"
}

func (f *accountsEnvironmentInstanceFacade) List(ctx context.Context, subaccountId string) (provisioning.EnvironmentInstancesResponseCollection, CommandResponse, error) {
	return doExecute[provisioning.EnvironmentInstancesResponseCollection](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

func (f *accountsEnvironmentInstanceFacade) Get(ctx context.Context, subaccountId string, environmentId string) (provisioning.EnvironmentInstanceResponseObject, CommandResponse, error) {
	return doExecute[provisioning.EnvironmentInstanceResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount":    subaccountId,
		"environmentID": environmentId,
	}))
}

type SubaccountEnvironmentInstanceCreateInput struct {
	DisplayName     string `btpcli:"displayName"`
	EnvironmentType string `btpcli:"environmentType"`
	Landscape       string `btpcli:"landscapeLabel"`
	Parameters      string `btpcli:"parameters"`
	Plan            string `btpcli:"plan"`
	Service         string `btpcli:"service"`
	SubaccountID    string `btpcli:"subaccount"`
}

func (f *accountsEnvironmentInstanceFacade) Create(ctx context.Context, args *SubaccountEnvironmentInstanceCreateInput) (provisioning.EnvironmentInstanceResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return provisioning.EnvironmentInstanceResponseObject{}, CommandResponse{}, err
	}

	return doExecute[provisioning.EnvironmentInstanceResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *accountsEnvironmentInstanceFacade) Update(ctx context.Context, subaccountId string, environmentId string, plan string, parameters string) (struct{}, CommandResponse, error) {
	return doExecute[struct{}](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), map[string]string{
		"subaccount":    subaccountId,
		"environmentID": environmentId,
		"plan":          plan,
		"parameters":    parameters,
	}))
}

func (f *accountsEnvironmentInstanceFacade) Delete(ctx context.Context, subaccountId string, environmentId string) (provisioning.EnvironmentInstanceResponseObject, CommandResponse, error) {
	return doExecute[provisioning.EnvironmentInstanceResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount":    subaccountId,
		"environmentID": environmentId,
		"confirm":       "true",
	}))
}
