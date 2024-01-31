package btpcli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis_entitlements"
)

const (
	subaccountEntityType = "SUBACCOUNT"
	directoryEntityType  = "DIRECTORY"
)

type DirectoryAssignmentInput struct {
	DirectoryId          string
	ServiceName          string
	ServicePlanName      string
	Amount               int
	Distribute           bool
	AutoAssign           bool
	AutoDistributeAmount int
}

func newAccountsEntitlementFacade(cliClient *v2Client) accountsEntitlementFacade {
	return accountsEntitlementFacade{cliClient: cliClient}
}

type accountsEntitlementFacade struct {
	cliClient *v2Client
}

func (f *accountsEntitlementFacade) getCommand() string {
	return "accounts/entitlement"
}

func (f *accountsEntitlementFacade) ListByGlobalAccount(ctx context.Context) (cis_entitlements.EntitledAndAssignedServicesResponseObject, CommandResponse, error) {
	return doExecute[cis_entitlements.EntitledAndAssignedServicesResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *accountsEntitlementFacade) ListBySubaccount(ctx context.Context, subaccountId string) (cis_entitlements.EntitledAndAssignedServicesResponseObject, CommandResponse, error) {
	return doExecute[cis_entitlements.EntitledAndAssignedServicesResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount":    f.cliClient.GetGlobalAccountSubdomain(),
		"subaccountFilter": subaccountId,
	}))
}

func (f *accountsEntitlementFacade) ListBySubaccountWithDirectoryParent(ctx context.Context, subaccountId string, directoryId string) (cis_entitlements.EntitledAndAssignedServicesResponseObject, CommandResponse, error) {
	return doExecute[cis_entitlements.EntitledAndAssignedServicesResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount":    f.cliClient.GetGlobalAccountSubdomain(),
		"subaccountFilter": subaccountId,
		"directory":        directoryId,
	}))
}

func (f *accountsEntitlementFacade) ListByDirectory(ctx context.Context, directoryId string) (cis_entitlements.EntitledAndAssignedServicesResponseObject, CommandResponse, error) {
	return doExecute[cis_entitlements.EntitledAndAssignedServicesResponseObject](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
		"directory":     directoryId,
	}))
}

func (f *accountsEntitlementFacade) AssignToSubaccount(ctx context.Context, directoryId string, subaccountId string, serviceName string, servicePlanName string, amount int) (CommandResponse, error) {

	params := map[string]string{
		"globalAccount":   f.cliClient.GetGlobalAccountSubdomain(),
		"subaccount":      subaccountId,
		"serviceName":     serviceName,
		"servicePlanName": servicePlanName,
		"amount":          fmt.Sprintf("%d", amount),
	}

	if len(directoryId) > 0 {
		params["directoryID"] = directoryId
	}
	_, res, err := doExecute[cis_entitlements.EntitlementAssignmentResponseObject](f.cliClient, ctx, NewAssignRequest(f.getCommand(), params))

	return res, err
}

func (f *accountsEntitlementFacade) EnableInSubaccount(ctx context.Context, directoryId string, subaccountId string, serviceName string, servicePlanName string) (CommandResponse, error) {

	params := map[string]string{
		"globalAccount":   f.cliClient.GetGlobalAccountSubdomain(),
		"subaccount":      subaccountId,
		"serviceName":     serviceName,
		"servicePlanName": servicePlanName,
		"enable":          "true",
	}

	if len(directoryId) > 0 {
		params["directoryID"] = directoryId
	}
	_, res, err := doExecute[cis_entitlements.EntitlementAssignmentResponseObject](f.cliClient, ctx, NewAssignRequest(f.getCommand(), params))

	return res, err
}

func (f *accountsEntitlementFacade) DisableInSubaccount(ctx context.Context, directoryId string, subaccountId string, serviceName string, servicePlanName string) (CommandResponse, error) {

	params := map[string]string{
		"globalAccount":   f.cliClient.GetGlobalAccountSubdomain(),
		"subaccount":      subaccountId,
		"serviceName":     serviceName,
		"servicePlanName": servicePlanName,
		"enable":          "false",
	}

	if len(directoryId) > 0 {
		params["directoryID"] = directoryId
	}

	_, res, err := doExecute[cis_entitlements.EntitlementAssignmentResponseObject](f.cliClient, ctx, NewAssignRequest(f.getCommand(), params))

	return res, err
}

type UnfoldedAssignment struct {
	Service    cis_entitlements.AssignedServiceResponseObject
	Plan       cis_entitlements.AssignedServicePlanResponseObject
	Assignment cis_entitlements.AssignedServicePlanSubaccountDto
}

type UnfoldedEntitlement struct {
	Service cis_entitlements.EntitledServicesResponseObject
	Plan    cis_entitlements.ServicePlanResponseObject
}

func (f *accountsEntitlementFacade) GetAssignedBySubaccount(ctx context.Context, subaccountId, serviceName string, servicePlanName string, isParentGlobalAccount bool, parentId string) (*UnfoldedAssignment, CommandResponse, error) {
	var cliRes cis_entitlements.EntitledAndAssignedServicesResponseObject
	var comRes CommandResponse
	var err error

	if isParentGlobalAccount {
		cliRes, comRes, err = f.ListBySubaccount(ctx, subaccountId)
	} else {
		cliRes, comRes, err = f.ListBySubaccountWithDirectoryParent(ctx, subaccountId, parentId)
	}

	if err != nil {
		return nil, comRes, err
	}

	for _, assignedService := range cliRes.AssignedServices {
		if assignedService.Name != serviceName {
			continue
		}

		servicePlan, assignment := f.searchPlansAndAssignments(assignedService.ServicePlans, servicePlanName, subaccountEntityType, subaccountId)
		if assignment != nil {
			return &UnfoldedAssignment{
				Service:    assignedService,
				Plan:       *servicePlan,
				Assignment: *assignment,
			}, comRes, nil
		}
	}

	return nil, comRes, nil
}

func (f *accountsEntitlementFacade) searchPlansAndAssignments(servicePlans []cis_entitlements.AssignedServicePlanResponseObject, servicePlanName string, entityType string, entityId string) (*cis_entitlements.AssignedServicePlanResponseObject, *cis_entitlements.AssignedServicePlanSubaccountDto) {
	for _, servicePlan := range servicePlans {
		if servicePlan.Name != servicePlanName {
			continue
		}

		for _, assignment := range servicePlan.AssignmentInfo {
			if assignment.EntityType == entityType && assignment.EntityId == entityId {
				return &servicePlan, &assignment
			}
		}
	}
	return nil, nil
}

func (f *accountsEntitlementFacade) searchPlansForEntitlement(servicePlans []cis_entitlements.ServicePlanResponseObject, servicePlanName string, entityType string, entityId string) *cis_entitlements.ServicePlanResponseObject {
	for _, servicePlan := range servicePlans {
		if servicePlan.Name == servicePlanName {
			return &servicePlan
		}
	}
	return nil
}

func (f *accountsEntitlementFacade) AssignToDirectory(ctx context.Context, dirAssignmentInput DirectoryAssignmentInput) (CommandResponse, error) {
	_, res, err := doExecute[cis_entitlements.EntitlementAssignmentResponseObject](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"globalAccount":        f.cliClient.GetGlobalAccountSubdomain(),
		"directory":            dirAssignmentInput.DirectoryId,
		"serviceName":          dirAssignmentInput.ServiceName,
		"servicePlanName":      dirAssignmentInput.ServicePlanName,
		"amount":               fmt.Sprintf("%d", dirAssignmentInput.Amount),
		"distribute":           strconv.FormatBool(dirAssignmentInput.Distribute),
		"autoAssign":           strconv.FormatBool(dirAssignmentInput.AutoAssign),
		"autoDistributeAmount": fmt.Sprintf("%d", dirAssignmentInput.AutoDistributeAmount),
	}))

	return res, err
}

func (f *accountsEntitlementFacade) EnableInDirectory(ctx context.Context, directoryId string, serviceName string, servicePlanName string, distribute bool, autoAssign bool) (CommandResponse, error) {
	_, res, err := doExecute[cis_entitlements.EntitlementAssignmentResponseObject](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"globalAccount":   f.cliClient.GetGlobalAccountSubdomain(),
		"directory":       directoryId,
		"serviceName":     serviceName,
		"servicePlanName": servicePlanName,
		"enable":          "true",
		"distribute":      strconv.FormatBool(distribute),
		"autoAssign":      strconv.FormatBool(autoAssign),
	}))

	return res, err
}

func (f *accountsEntitlementFacade) DisableInDirectory(ctx context.Context, directoryId string, serviceName string, servicePlanName string, distribute bool, autoAssign bool) (CommandResponse, error) {
	_, res, err := doExecute[cis_entitlements.EntitlementAssignmentResponseObject](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"globalAccount":   f.cliClient.GetGlobalAccountSubdomain(),
		"directory":       directoryId,
		"serviceName":     serviceName,
		"servicePlanName": servicePlanName,
		"enable":          "false",
		"distribute":      strconv.FormatBool(distribute),
		"autoAssign":      strconv.FormatBool(autoAssign),
	}))

	return res, err
}

func (f *accountsEntitlementFacade) GetEntitledByDirectory(ctx context.Context, directoryId, serviceName string, servicePlanName string) (*UnfoldedEntitlement, CommandResponse, error) {
	cliRes, comRes, err := f.ListByDirectory(ctx, directoryId)

	if err != nil {
		return nil, comRes, err
	}

	for _, entitledService := range cliRes.EntitledServices {
		if entitledService.Name != serviceName {
			continue
		}

		servicePlan := f.searchPlansForEntitlement(entitledService.ServicePlans, servicePlanName, directoryEntityType, directoryId)
		if servicePlan != nil {
			return &UnfoldedEntitlement{
				Service: entitledService,
				Plan:    *servicePlan,
			}, comRes, nil
		}
	}

	return nil, comRes, nil
}
