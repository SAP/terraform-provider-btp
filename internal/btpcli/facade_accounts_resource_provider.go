package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/provisioning"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newAccountsResourceProviderFacade(cliClient *v2Client) accountsResourceProviderFacade {
	return accountsResourceProviderFacade{cliClient: cliClient}
}

type accountsResourceProviderFacade struct {
	cliClient *v2Client
}

func (f *accountsResourceProviderFacade) getCommand() string {
	return "accounts/resource-provider"
}

func (f *accountsResourceProviderFacade) List(ctx context.Context) ([]provisioning.ResourceProviderResponseObject, CommandResponse, error) {
	return doExecute[[]provisioning.ResourceProviderResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *accountsResourceProviderFacade) Get(ctx context.Context, provider string, technicalName string) (provisioning.ResourceProviderResponseObject, CommandResponse, error) {
	return doExecute[provisioning.ResourceProviderResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"provider":      provider,
		"technicalName": technicalName,
	}))
}

type GlobalaccountResourceProviderCreateUpdateInput struct {
	Provider      string `btpcli:"provider"`
	TechnicalName string `btpcli:"technicalName"`
	DisplayName   string `btpcli:"displayName"`
	Description   string `btpcli:"description"`
	Configuration string `btpcli:"configurationInfo"`
	Globalaccount string `btpcli:"globalAccount"`
}

func (f *accountsResourceProviderFacade) Create(ctx context.Context, args GlobalaccountResourceProviderCreateUpdateInput) (provisioning.ResourceProviderResponseObject, CommandResponse, error) {
	args.Globalaccount = f.cliClient.GetGlobalAccountSubdomain()
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return provisioning.ResourceProviderResponseObject{}, CommandResponse{}, err
	}

	return doExecute[provisioning.ResourceProviderResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *accountsResourceProviderFacade) Update(ctx context.Context, args GlobalaccountResourceProviderCreateUpdateInput) (provisioning.ResourceProviderResponseObject, CommandResponse, error) {
	args.Globalaccount = f.cliClient.GetGlobalAccountSubdomain()
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return provisioning.ResourceProviderResponseObject{}, CommandResponse{}, err
	}

	return doExecute[provisioning.ResourceProviderResponseObject](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
}

func (f *accountsResourceProviderFacade) Delete(ctx context.Context, provider string, technicalName string) (provisioning.ResourceProviderResponseObject, CommandResponse, error) {
	return doExecute[provisioning.ResourceProviderResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"provider":      provider,
		"technicalName": technicalName,
		"confirm":       "true",
	}))
}
