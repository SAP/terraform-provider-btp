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

type ApiCredentialInput struct {
	SubaccountId        string `btpcli:"subaccount"`
	DirectoryId         string `btpcli:"directory"`
	Name            	string `btpcli:"name,omitempty"`
	Certificate			string `btpcli:"certificate,omitempty"`
	ReadOnly 			bool   `btpcli:"readOnly,omitempty"`
}

func (f *securityApiCredentialFacade) CreateBySubaccount(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredentialSubaccount, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredentialSubaccount{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredentialSubaccount](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *securityApiCredentialFacade) DeleteBySubaccount(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredentialSubaccount, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredentialSubaccount{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredentialSubaccount](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), params))
}

func (f *securityApiCredentialFacade) GetBySubaccount(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredentialSubaccount, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredentialSubaccount{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredentialSubaccount](f.cliClient, ctx, NewGetRequest(f.getCommand(), params))
}

func (f *securityApiCredentialFacade) CreateByDirectory(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredentialSubaccount, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredentialSubaccount{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredentialSubaccount](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *securityApiCredentialFacade) DeleteByDirectory(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredentialSubaccount, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredentialSubaccount{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredentialSubaccount](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), params))
}

func (f *securityApiCredentialFacade) GetByDirectory(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredentialSubaccount, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredentialSubaccount{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredentialSubaccount](f.cliClient, ctx, NewGetRequest(f.getCommand(), params))
}