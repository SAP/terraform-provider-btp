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
	Subaccount    string `btpcli:"subaccount"`
	Directory     string `btpcli:"directory"`
	GlobalAccount string `btpcli:"globalAccount"`
	Name          string `btpcli:"name,omitempty"`
	Certificate   string `btpcli:"certificate,omitempty"`
	ReadOnly      bool   `btpcli:"readOnly,omitempty"`
}

func (f *securityApiCredentialFacade) CreateByDirectoryorSubaccount(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredential, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredential{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredential](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *securityApiCredentialFacade) DeleteByDirectoryorSubaccount(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredential, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredential{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredential](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), params))
}

func (f *securityApiCredentialFacade) GetByDirectoryorSubaccount(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredential, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredential{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredential](f.cliClient, ctx, NewGetRequest(f.getCommand(), params))
}

func (f *securityApiCredentialFacade) CreateByGlobalAccount(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredential, CommandResponse, error) {

	args.GlobalAccount = f.cliClient.GetGlobalAccountSubdomain()

	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredential{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredential](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *securityApiCredentialFacade) DeleteByGlobalAccount(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredential, CommandResponse, error) {

	args.GlobalAccount = f.cliClient.GetGlobalAccountSubdomain()

	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredential{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredential](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), params))
}

func (f *securityApiCredentialFacade) GetByGlobalAccount(ctx context.Context, args *ApiCredentialInput) (xsuaa_api.ApiCredential, CommandResponse, error) {

	args.GlobalAccount = f.cliClient.GetGlobalAccountSubdomain()

	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_api.ApiCredential{}, CommandResponse{}, err
	}

	return doExecute[xsuaa_api.ApiCredential](f.cliClient, ctx, NewGetRequest(f.getCommand(), params))
}
