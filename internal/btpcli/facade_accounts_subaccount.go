package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/saas_manager_service"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newAccountsSubaccountFacade(cliClient *v2Client) accountsSubaccountFacade {
	return accountsSubaccountFacade{cliClient: cliClient}
}

type accountsSubaccountFacade struct {
	cliClient *v2Client
}

func (f *accountsSubaccountFacade) getCommand() string {
	return "accounts/subaccount"
}

func (f *accountsSubaccountFacade) List(ctx context.Context, labelsFilter string) (cis.ResponseCollectionSubaccountResponseObject, CommandResponse, error) {
	params := map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}

	if len(labelsFilter) > 0 {
		params["labelsFilter"] = labelsFilter

	}

	return doExecute[cis.ResponseCollectionSubaccountResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), params))
}

func (f *accountsSubaccountFacade) Get(ctx context.Context, subaccountId string) (cis.SubaccountResponseObject, CommandResponse, error) {
	return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

type SubaccountCreateInput struct { // TODO support all options
	BetaEnabled       bool                `btpcli:"betaEnabled"`
	Description       string              `btpcli:"description"`
	Directory         string              `btpcli:"directoryID"`
	DisplayName       string              `btpcli:"displayName"`
	Labels            map[string][]string `btpcli:"labels"`
	Region            string              `btpcli:"region"`
	Subdomain         string              `btpcli:"subdomain"`
	UsedForProduction bool                `btpcli:"usedForProduction"`
	//Globalaccount     string `json:"globalAccount"`
	//SubaccountAdmins  string `json:"subaccountAdmins"`
}

func (f *accountsSubaccountFacade) Create(ctx context.Context, args *SubaccountCreateInput) (cis.SubaccountResponseObject, CommandResponse, error) {

	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return cis.SubaccountResponseObject{}, CommandResponse{}, err
	}

	return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *accountsSubaccountFacade) Update(ctx context.Context, subaccountId string, displayName string) (cis.SubaccountResponseObject, CommandResponse, error) { // TODO switch to object
	return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), map[string]string{
		"subaccount":  subaccountId,
		"displayName": displayName,
	}))
}

func (f *accountsSubaccountFacade) Delete(ctx context.Context, subaccountId string) (cis.SubaccountResponseObject, CommandResponse, error) {
	return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount":  subaccountId,
		"confirm":     "true",
		"forceDelete": "true",
	}))
}

func (f *accountsSubaccountFacade) Subscribe(ctx context.Context, subaccountId string, appName string, planName string, parameters string) (saas_manager_service.SubscriptionAssignmentResponseObject, CommandResponse, error) {
	return doExecute[saas_manager_service.SubscriptionAssignmentResponseObject](f.cliClient, ctx, NewSubscribeRequest(f.getCommand(), map[string]string{
		"subaccount":         subaccountId,
		"appName":            appName,
		"planName":           planName,
		"subscriptionParams": parameters,
	}))
}

func (f *accountsSubaccountFacade) Unsubscribe(ctx context.Context, subaccountId string, appName string) (saas_manager_service.SubscriptionAssignmentResponseObject, CommandResponse, error) {
	return doExecute[saas_manager_service.SubscriptionAssignmentResponseObject](f.cliClient, ctx, NewUnsubscribeRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"appName":    appName,
		"confirm":    "true",
	}))
}
