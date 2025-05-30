package btpcli

import (
	"context"
	"strings"

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
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"subaccount":    subaccountId,
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
	UsedForProduction string              `btpcli:"usedForProduction"`
	Globalaccount     string              `btpcli:"globalAccount"`
	AdminDirectoryId  string              `btpcli:"adminDirectory"`
	//SubaccountAdmins  string `json:"subaccountAdmins"`
}

type SubaccountUpdateInput struct {
	BetaEnabled       bool                `btpcli:"betaEnabled"`
	Description       string              `btpcli:"description"`
	Directory         string              `btpcli:"directoryID"`
	DisplayName       string              `btpcli:"displayName"`
	Labels            map[string][]string `btpcli:"labels"`
	SubaccountId      string              `btpcli:"subaccount"`
	UsedForProduction string              `btpcli:"usedForProduction"`
	Globalaccount     string              `btpcli:"globalAccount"`
}

func (f *accountsSubaccountFacade) Create(ctx context.Context, args *SubaccountCreateInput) (cis.SubaccountResponseObject, CommandResponse, error) {

	args.Globalaccount = f.cliClient.GetGlobalAccountSubdomain()

	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return cis.SubaccountResponseObject{}, CommandResponse{}, err
	}

	return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *accountsSubaccountFacade) Update(ctx context.Context, args *SubaccountUpdateInput) (cis.SubaccountResponseObject, CommandResponse, error) { // TODO switch to object

	args.Globalaccount = f.cliClient.GetGlobalAccountSubdomain()

	// Mapping of all params except for usedForProduction
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return cis.SubaccountResponseObject{}, CommandResponse{}, err
	}

	return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
}

func (f *accountsSubaccountFacade) Delete(ctx context.Context, subaccountId string, directoryId string) (cis.SubaccountResponseObject, CommandResponse, error) {

	// We first try to delete the subaccount with force-delete enabled
	// This might fail, if the subaccount setting has not enabled this feature
	requestArgs := map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"subaccount":    subaccountId,
		"confirm":       "true",
		"forceDelete":   "true",
	}

	if len(directoryId) > 0 {
		//if the parent is a managed directory, the directoryID must be set to make sure the right authorizations are validated
		requestArgs["directoryID"] = directoryId
	}

	cisResponse, cmdResponse, err := doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), requestArgs))

	if err == nil {
		return cisResponse, cmdResponse, nil
	}
	// Check if the error is due to the fact that the force-delete option is not available
	if err != nil && (strings.Contains(err.Error(), "Subaccount cannot be deleted with forceDelete=true due to the global account settings") || strings.Contains(err.Error(), "Subaccount is marked as used for production and cannot be deleted with forceDelete=true")) {
		// Retry with force-delete disabled
		requestArgs["forceDelete"] = "false"
		return doExecute[cis.SubaccountResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), requestArgs))
	}

	return cisResponse, cmdResponse, err
}

func (f *accountsSubaccountFacade) Subscribe(ctx context.Context, subaccountId string, appName string, planName string, parameters string) (saas_manager_service.SubscriptionAssignmentResponseObject, CommandResponse, error) {
	commandOptions := map[string]string{
		"subaccount": subaccountId,
		"appName":    appName,
	}

	// Some app subscriptions cannot handle empty JSON objects as subscription parameters. In this case, the parameters must not be transferred to the API call.
	if parameters != "{}" {
		commandOptions["subscriptionParams"] = parameters
	}

	// The plan name can be empty in case of subscription hosted in the same account. In this case the plan name is not required and must not be transferred to the API call.
	if len(planName) > 0 {
		commandOptions["planName"] = planName
	}
	return doExecute[saas_manager_service.SubscriptionAssignmentResponseObject](f.cliClient, ctx, NewSubscribeRequest(f.getCommand(), commandOptions))
}

func (f *accountsSubaccountFacade) Unsubscribe(ctx context.Context, subaccountId string, appName string) (saas_manager_service.SubscriptionAssignmentResponseObject, CommandResponse, error) {
	return doExecute[saas_manager_service.SubscriptionAssignmentResponseObject](f.cliClient, ctx, NewUnsubscribeRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"appName":    appName,
		"confirm":    "true",
	}))
}
