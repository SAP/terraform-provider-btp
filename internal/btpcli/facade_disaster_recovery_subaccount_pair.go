package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cdr"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
)

func newDisasterRecoverySubaccountPairFacade(cliClient *v2Client) disasterRecoverySubaccountPairFacade {
	return disasterRecoverySubaccountPairFacade{cliClient: cliClient}
}

type disasterRecoverySubaccountPairFacade struct {
	cliClient *v2Client
}

type SubaccountPairCreateInput struct {
	SubaccountId     string `btpcli:"subaccount"`
	WithSubaccountId string `btpcli:"with-subaccount"`
}

func (f *disasterRecoverySubaccountPairFacade) getCommand() string {
	return "disaster-recovery/subaccount-pair"
}

func (f *disasterRecoverySubaccountPairFacade) Get(ctx context.Context, subaccountId string) (cdr.GetSubaccountPairResponse, CommandResponse, error) {
	return doExecute[cdr.GetSubaccountPairResponse](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

func (f *disasterRecoverySubaccountPairFacade) Create(ctx context.Context, args *SubaccountPairCreateInput) (cdr.CreateSubaccountPairResponse, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return cdr.CreateSubaccountPairResponse{}, CommandResponse{}, err
	}

	return doExecute[cdr.CreateSubaccountPairResponse](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *disasterRecoverySubaccountPairFacade) Delete(ctx context.Context, subaccountId string) (cdr.DeleteSubaccountPairResponse, CommandResponse, error) {
	return doExecute[cdr.DeleteSubaccountPairResponse](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}
