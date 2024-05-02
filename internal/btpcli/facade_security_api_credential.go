package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_api"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newSecurityApiCredentialFacade(cliClient *v2Client) securityApiCredentialFacade {
	return securityApiCredentialFacade{cliClient: cliClient}
}

type securityApiCredentialFacade struct {
	cliClient *v2Client
}

func (f *securityApiCredentialFacade) getCommand() string {
	return "security/api-credential"
}

type ApiCredentialCreateInput struct {
	SubaccountId        string `btpcli:"subaccount"`
	Name            	string `btpcli:"name,omitempty"`
	Certificate			string `btpcli:"certificate,omitempty"`
	ReadOnly 			bool   `btpcli:"readOnly,omitempty"`
}

func (f *securityApiCredentialFacade) CreateBySubaccount(ctx context.Context, args *ApiCredentialCreateInput) (xsuaa_api.ApiCredentialCreateBody, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredentialCreateBody{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredentialCreateBody](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}