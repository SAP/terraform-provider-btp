package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/saas_manager_service"
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

func (f *accountsSubaccountFacade) List(ctx context.Context, labelsFilter string) (cis.ResponseCollectionSubaccountResponseObject, *CommandResponse, error) {
	params := map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}

	if len(labelsFilter) > 0 {
		params["labelsFilter"] = labelsFilter

	}

	return doExecute[cis.ResponseCollectionSubaccountResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), params))
}

func (f *accountsSubaccountFacade) Get(ctx context.Context, subaccountId string) (cis.SubaccountResponseObject, *CommandResponse, error) {
	return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

type SubaccountCreateInput struct { // TODO support all options
	//BetaEnabled       bool   `json:"betaEnabled"`
	//Description       string `json:"description"`
	//Directory         string `json:"directoryID"`
	DisplayName string `json:"displayName"`
	//Globalaccount     string `json:"globalAccount"`
	//Labels            string `json:"labels"`
	Region string `json:"region"`
	//SubaccountAdmins  string `json:"subaccountAdmins"`
	Subdomain string `json:"subdomain"`
	//UsedForProduction bool   `json:"usedForProduction"`
}

func (f *accountsSubaccountFacade) Create(ctx context.Context, displayName string, subdomain string, region string) (cis.SubaccountResponseObject, *CommandResponse, error) { // TODO switch to object
	return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), map[string]string{
		"displayName": displayName,
		"subdomain":   subdomain,
		"region":      region,
	}))
}

func (f *accountsSubaccountFacade) Update(ctx context.Context, subaccountId string, displayName string) (cis.SubaccountResponseObject, *CommandResponse, error) { // TODO switch to object
	return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), map[string]string{
		"subaccount":  subaccountId,
		"displayName": displayName,
	}))
}

func (f *accountsSubaccountFacade) Delete(ctx context.Context, subaccountId string) (cis.SubaccountResponseObject, *CommandResponse, error) {
	return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount":  subaccountId,
		"confirm":     "true",
		"forceDelete": "true",
	}))
}

func (f *accountsSubaccountFacade) Subscribe(ctx context.Context, subaccountId string, appName string, planName string, parameters string) (saas_manager_service.SubscriptionAssignmentResponseObject, *CommandResponse, error) {
	return doExecute[saas_manager_service.SubscriptionAssignmentResponseObject](f.cliClient, ctx, NewSubscribeRequest(f.getCommand(), map[string]string{
		"subaccount":         subaccountId,
		"appName":            appName,
		"planName":           planName,
		"subscriptionParams": parameters,
	}))
}

func (f *accountsSubaccountFacade) Unsubscribe(ctx context.Context, subaccountId string, appName string) (saas_manager_service.SubscriptionAssignmentResponseObject, *CommandResponse, error) {
	return doExecute[saas_manager_service.SubscriptionAssignmentResponseObject](f.cliClient, ctx, NewUnsubscribeRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"appName":    appName,
		"confirm":    "true",
	}))
}
