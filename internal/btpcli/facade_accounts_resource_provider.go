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

func (f *accountsResourceProviderFacade) List(ctx context.Context) ([]provisioning.ResourceProviderResponseObject, *CommandResponse, error) {
	return doExecute[[]provisioning.ResourceProviderResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *accountsResourceProviderFacade) Get(ctx context.Context, resourceProvider string, resourceTechnicalName string) (provisioning.ResourceProviderResponseObject, *CommandResponse, error) {
	return doExecute[provisioning.ResourceProviderResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"provider":      resourceProvider,
		"technicalName": resourceTechnicalName,
	}))
}

type GlobalaccountResourceProviderCreateInput struct {
	Provider          string `btpcli:"provider"`
	TechnicalName     string `btpcli:"technicalName"`
	DisplayName       string `btpcli:"displayName"`
	Description       string `btpcli:"description"`
	ConfigurationInfo string `btpcli:"configurationInfo"`
}

func (f *accountsResourceProviderFacade) Create(ctx context.Context, args GlobalaccountResourceProviderCreateInput) (provisioning.ResourceProviderResponseObject, *CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return provisioning.ResourceProviderResponseObject{}, nil, err
	}

	params["globalAccount"] = f.cliClient.GetGlobalAccountSubdomain()

	return doExecute[provisioning.ResourceProviderResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *accountsResourceProviderFacade) Delete(ctx context.Context, resourceProvider string, resourceTechnicalName string) (provisioning.ResourceProviderResponseObject, *CommandResponse, error) {
	return doExecute[provisioning.ResourceProviderResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"provider":      resourceProvider,
		"technicalName": resourceTechnicalName,
		"confirm":       "true",
	}))
}
