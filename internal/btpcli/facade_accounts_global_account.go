package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newAccountsGlobalAccountFacade(cliClient *v2Client) accountsGlobalAccountFacade {
	return accountsGlobalAccountFacade{cliClient: cliClient}
}

type GlobalAccountUpdateInput struct {
	DisplayName                   string `btpcli:"displayName"`
	Description                   string `btpcli:"description"`
	EnableSubaccountForceDeletion *bool  `btpcli:"enableSubaccountForceDeletion"`
	Globalaccount                 string `btpcli:"globalAccount"`
}

type accountsGlobalAccountFacade struct {
	cliClient *v2Client
}

func (f *accountsGlobalAccountFacade) getCommand() string {
	return "accounts/global-account"
}

func (f *accountsGlobalAccountFacade) Get(ctx context.Context) (cis.GlobalAccountResponseObject, CommandResponse, error) {
	return doExecute[cis.GlobalAccountResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *accountsGlobalAccountFacade) GetWithHierarchy(ctx context.Context) (cis.GlobalAccountResponseObject, CommandResponse, error) {
	return doExecute[cis.GlobalAccountResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"showHierarchy": "true",
	}))
}

func (f *accountsGlobalAccountFacade) Update(ctx context.Context, args *GlobalAccountUpdateInput) (cis.GlobalAccountResponseObject, CommandResponse, error) {
	args.Globalaccount = f.cliClient.GetGlobalAccountSubdomain()

	params, err := tfutils.ToBTPCLIParamsMap(args)
	if err != nil {
		return cis.GlobalAccountResponseObject{}, CommandResponse{}, err
	}

	return doExecute[cis.GlobalAccountResponseObject](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
}
