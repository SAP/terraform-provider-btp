package btpcli

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis_entitlements"
)

func newAccountsEntitlementFacade(cliClient *v2Client) accountsEntitlementFacade {
	return accountsEntitlementFacade{cliClient: cliClient}
}

type accountsEntitlementFacade struct {
	cliClient *v2Client
}

func (f *accountsEntitlementFacade) getCommand() string {
	return "accounts/entitlement"
}

func (f *accountsEntitlementFacade) ListByGlobalAccount(ctx context.Context) (cis_entitlements.EntitledAndAssignedServicesResponseObject, *CommandResponse, error) {
	return doExecute[cis_entitlements.EntitledAndAssignedServicesResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *accountsEntitlementFacade) ListBySubaccount(ctx context.Context, subaccountId string) (cis_entitlements.EntitledAndAssignedServicesResponseObject, *CommandResponse, error) {
	return doExecute[cis_entitlements.EntitledAndAssignedServicesResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccountFilter": subaccountId,
	}))
}

func (f *accountsEntitlementFacade) ListByDirectory(ctx context.Context, directoryId string) (cis_entitlements.EntitledAndAssignedServicesResponseObject, *CommandResponse, error) {
	return doExecute[cis_entitlements.EntitledAndAssignedServicesResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"directory": directoryId,
	}))
}

func (f *accountsEntitlementFacade) AssignToSubaccount(ctx context.Context, subaccountId string, serviceName string, servicePlanName string, amount int) (*CommandResponse, error) {
	_, res, err := doExecute[cis_entitlements.EntitlementAssignmentResponseObject](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"subaccount":      subaccountId,
		"serviceName":     serviceName,
		"servicePlanName": servicePlanName,
		"amount":          fmt.Sprintf("%d", amount),
	}))

	return res, err
}

func (f *accountsEntitlementFacade) EnableInSubaccount(ctx context.Context, subaccountId string, serviceName string, servicePlanName string) (*CommandResponse, error) {
	_, res, err := doExecute[cis_entitlements.EntitlementAssignmentResponseObject](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"subaccount":      subaccountId,
		"serviceName":     serviceName,
		"servicePlanName": servicePlanName,
		"enable":          "true",
	}))

	return res, err
}

func (f *accountsEntitlementFacade) DisableInSubaccount(ctx context.Context, subaccountId string, serviceName string, servicePlanName string) (*CommandResponse, error) {
	_, res, err := doExecute[cis_entitlements.EntitlementAssignmentResponseObject](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"subaccount":      subaccountId,
		"serviceName":     serviceName,
		"servicePlanName": servicePlanName,
		"enable":          "false",
	}))

	return res, err
}

type UnfoldedEntitlement struct {
	Service    cis_entitlements.AssignedServiceResponseObject
	Plan       cis_entitlements.AssignedServicePlanResponseObject
	Assignment cis_entitlements.AssignedServicePlanSubaccountDto
}

func (f *accountsEntitlementFacade) GetAssignedBySubaccount(ctx context.Context, subaccountId, serviceName string, servicePlanName string) (*UnfoldedEntitlement, *CommandResponse, error) {
	cliRes, comRes, err := f.ListBySubaccount(ctx, subaccountId)

	if err != nil {
		return nil, comRes, err
	}

	for _, assignedService := range cliRes.AssignedServices {
		if assignedService.Name != serviceName {
			continue
		}

		for _, servicePlan := range assignedService.ServicePlans {
			if servicePlan.Name != servicePlanName {
				continue
			}

			for _, assignment := range servicePlan.AssignmentInfo {
				if assignment.EntityType == "SUBACCOUNT" && assignment.EntityId == subaccountId {
					return &UnfoldedEntitlement{
						Service:    assignedService,
						Plan:       servicePlan,
						Assignment: assignment,
					}, comRes, nil
				}
			}
		}
	}

	return nil, comRes, nil
}
