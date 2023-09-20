package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_trust"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newSecurityTrustFacade(cliClient *v2Client) securityTrustFacade {
	return securityTrustFacade{cliClient: cliClient}
}

type securityTrustFacade struct {
	cliClient *v2Client
}

func (f *securityTrustFacade) getCommand() string {
	return "security/trust"
}

func (f *securityTrustFacade) ListByGlobalAccount(ctx context.Context) (xsuaa_trust.TrustConfigurationResponseCollectionObject, CommandResponse, error) {
	return doExecute[xsuaa_trust.TrustConfigurationResponseCollectionObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *securityTrustFacade) GetByGlobalAccount(ctx context.Context, origin string) (xsuaa_trust.TrustConfigurationResponseObject, CommandResponse, error) {
	return doExecute[xsuaa_trust.TrustConfigurationResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"origin":        origin,
	}))
}

func (f *securityTrustFacade) ListBySubaccount(ctx context.Context, subaccountId string) (xsuaa_trust.TrustConfigurationResponseCollectionObject, CommandResponse, error) {
	return doExecute[xsuaa_trust.TrustConfigurationResponseCollectionObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

func (f *securityTrustFacade) GetBySubaccount(ctx context.Context, subaccountId string, origin string) (xsuaa_trust.TrustConfigurationResponseObject, CommandResponse, error) {
	return doExecute[xsuaa_trust.TrustConfigurationResponseObject](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"origin":     origin,
	}))
}

type TrustConfigurationCreateInput struct {
	IdentityProvider string  `btpcli:"iasTenantUrl"`
	Name             *string `btpcli:"name"`
	Description      *string `btpcli:"description"`
	Origin           *string `btpcli:"origin"`
	Domain           *string `btpcli:"domain"`
}

func (f *securityTrustFacade) CreateByGlobalAccount(ctx context.Context, args TrustConfigurationCreateInput) (xsuaa_trust.ModifyTrustConfigurationResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_trust.ModifyTrustConfigurationResponseObject{}, CommandResponse{}, err
	}

	params["globalAccount"] = f.cliClient.GetGlobalAccountSubdomain()

	return doExecute[xsuaa_trust.ModifyTrustConfigurationResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *securityTrustFacade) CreateBySubaccount(ctx context.Context, subaccountId string, args TrustConfigurationCreateInput) (xsuaa_trust.ModifyTrustConfigurationResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_trust.ModifyTrustConfigurationResponseObject{}, CommandResponse{}, err
	}

	params["subaccount"] = subaccountId

	return doExecute[xsuaa_trust.ModifyTrustConfigurationResponseObject](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

type TrustConfigurationUpdateInput struct {
	OriginKey             string  `btpcli:"originKey"`
	IdentityProvider      string  `btpcli:"iasTenantUrl"`
	Name                  *string `btpcli:"name"`
	Description           *string `btpcli:"description"`
	Domain                *string `btpcli:"domain"`
	LinkText              *string `btpcli:"linkText"`
	AvailableForUserLogon bool    `btpcli:"userLogon"`
	AutoCreateShadowUsers bool    `btpcli:"shadowUsers"`
	Status                string  `btpcli:"status"`
}

func (f *securityTrustFacade) UpdateBySubaccount(ctx context.Context, subaccountId string, args TrustConfigurationUpdateInput) (xsuaa_trust.TrustConfigurationResponseObject, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_trust.TrustConfigurationResponseObject{}, CommandResponse{}, err
	}

	params["subaccount"] = subaccountId
	params["refreshTrust"] = "true"

	return doExecute[xsuaa_trust.TrustConfigurationResponseObject](f.cliClient, ctx, NewUpdateRequest(f.getCommand(), params))
}

func (f *securityTrustFacade) DeleteByGlobalAccount(ctx context.Context, originKey string) (xsuaa_trust.ModifyTrustConfigurationResponseObject, CommandResponse, error) {
	return doExecute[xsuaa_trust.ModifyTrustConfigurationResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"originKey":     originKey,
		"confirm":       "true",
	}))
}

func (f *securityTrustFacade) DeleteBySubaccount(ctx context.Context, subaccountId string, originKey string) (xsuaa_trust.ModifyTrustConfigurationResponseObject, CommandResponse, error) {
	return doExecute[xsuaa_trust.ModifyTrustConfigurationResponseObject](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
		"originKey":  originKey,
		"confirm":    "true",
	}))
}
