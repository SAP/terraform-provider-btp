package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
)

func newSecurityRoleCollectionFacade(cliClient *v2Client) securityRoleCollectionFacade {
	return securityRoleCollectionFacade{cliClient: cliClient}
}

type securityRoleCollectionFacade struct {
	cliClient *v2Client
}

func (f *securityRoleCollectionFacade) getCommand() string {
	return "security/role-collection"
}

func (f *securityRoleCollectionFacade) ListByGlobalAccount(ctx context.Context) ([]xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[[]xsuaa_authz.RoleCollection](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *securityRoleCollectionFacade) GetByGlobalAccount(ctx context.Context, roleCollectionName string) (xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount":      f.cliClient.GetGlobalAccountSubdomain(),
		"roleCollectionName": roleCollectionName,
	}))
}

func (f *securityRoleCollectionFacade) CreateByGlobalAccount(ctx context.Context, roleCollectionName string, description string) (xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewCreateRequest(f.getCommand(), map[string]string{
		"globalAccount":      f.cliClient.GetGlobalAccountSubdomain(),
		"roleCollectionName": roleCollectionName,
		"description":        description,
	}))
}

func (f *securityRoleCollectionFacade) DeleteByGlobalAccount(ctx context.Context, roleCollectionName string) (xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"globalAccount":      f.cliClient.GetGlobalAccountSubdomain(),
		"roleCollectionName": roleCollectionName,
	}))
}

func (f *securityRoleCollectionFacade) ListBySubaccount(ctx context.Context, subaccountId string) ([]xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[[]xsuaa_authz.RoleCollection](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

func (f *securityRoleCollectionFacade) GetBySubaccount(ctx context.Context, subaccountId string, roleCollectionName string) (xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount":         subaccountId,
		"roleCollectionName": roleCollectionName,
	}))
}

func (f *securityRoleCollectionFacade) CreateBySubaccount(ctx context.Context, subaccountId string, roleCollectionName string, description string) (xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewCreateRequest(f.getCommand(), map[string]string{
		"subaccount":         subaccountId,
		"roleCollectionName": roleCollectionName,
		"description":        description,
	}))
}

func (f *securityRoleCollectionFacade) DeleteBySubaccount(ctx context.Context, subaccountId string, roleCollectionName string) (xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount":         subaccountId,
		"roleCollectionName": roleCollectionName,
	}))
}

func (f *securityRoleCollectionFacade) ListByDirectory(ctx context.Context, directoryId string) ([]xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[[]xsuaa_authz.RoleCollection](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"directory": directoryId,
	}))
}

func (f *securityRoleCollectionFacade) GetByDirectory(ctx context.Context, directoryId string, roleCollectionName string) (xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"directory":          directoryId,
		"roleCollectionName": roleCollectionName,
	}))
}

func (f *securityRoleCollectionFacade) CreateByDirectory(ctx context.Context, directoryId string, roleCollectionName string, description string) (xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewCreateRequest(f.getCommand(), map[string]string{
		"directory":          directoryId,
		"roleCollectionName": roleCollectionName,
		"description":        description,
	}))
}

func (f *securityRoleCollectionFacade) DeleteByDirectory(ctx context.Context, directoryId string, roleCollectionName string) (xsuaa_authz.RoleCollection, *CommandResponse, error) {
	return doExecute[xsuaa_authz.RoleCollection](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"directory":          directoryId,
		"roleCollectionName": roleCollectionName,
	}))
}

func (f *securityRoleCollectionFacade) AssignUserBySubaccount(ctx context.Context, subaccountId string, roleCollectionName string, username string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"subaccount":          subaccountId,
		"roleCollectionName":  roleCollectionName,
		"userName":            username,
		"origin":              origin,
		"createUserIfMissing": "true",
	}))
}

func (f *securityRoleCollectionFacade) UnassignUserBySubaccount(ctx context.Context, subaccountId string, roleCollectionName string, username string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewUnassignRequest(f.getCommand(), map[string]string{
		"subaccount":         subaccountId,
		"roleCollectionName": roleCollectionName,
		"userName":           username,
		"origin":             origin,
	}))
}

func (f *securityRoleCollectionFacade) AssignUserByDirectory(ctx context.Context, directoryId string, roleCollectionName string, username string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"directory":           directoryId,
		"roleCollectionName":  roleCollectionName,
		"userName":            username,
		"origin":              origin,
		"createUserIfMissing": "true",
	}))
}

func (f *securityRoleCollectionFacade) UnassignUserByDirectory(ctx context.Context, directoryId string, roleCollectionName string, username string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewUnassignRequest(f.getCommand(), map[string]string{
		"directory":          directoryId,
		"roleCollectionName": roleCollectionName,
		"userName":           username,
		"origin":             origin,
	}))
}

func (f *securityRoleCollectionFacade) AssignUserByGlobalaccount(ctx context.Context, roleCollectionName string, username string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"globalAccount":       f.cliClient.GetGlobalAccountSubdomain(),
		"roleCollectionName":  roleCollectionName,
		"userName":            username,
		"origin":              origin,
		"createUserIfMissing": "true",
	}))
}

func (f *securityRoleCollectionFacade) UnassignUserByGlobalaccount(ctx context.Context, roleCollectionName string, username string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewUnassignRequest(f.getCommand(), map[string]string{
		"globalAccount":      f.cliClient.GetGlobalAccountSubdomain(),
		"roleCollectionName": roleCollectionName,
		"userName":           username,
		"origin":             origin,
	}))
}

func (f *securityRoleCollectionFacade) AssignGroupBySubaccount(ctx context.Context, subaccountId string, roleCollectionName string, groupName string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"subaccount":          subaccountId,
		"roleCollectionName":  roleCollectionName,
		"group":               groupName,
		"origin":              origin,
		"createUserIfMissing": "true",
	}))
}

func (f *securityRoleCollectionFacade) UnassignGroupBySubaccount(ctx context.Context, subaccountId string, roleCollectionName string, groupName string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewUnassignRequest(f.getCommand(), map[string]string{
		"subaccount":         subaccountId,
		"roleCollectionName": roleCollectionName,
		"group":              groupName,
		"origin":             origin,
	}))
}

func (f *securityRoleCollectionFacade) AssignGroupByDirectory(ctx context.Context, directoryId string, roleCollectionName string, groupName string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"directory":           directoryId,
		"roleCollectionName":  roleCollectionName,
		"group":               groupName,
		"origin":              origin,
		"createUserIfMissing": "true",
	}))
}

func (f *securityRoleCollectionFacade) UnassignGroupByDirectory(ctx context.Context, directoryId string, roleCollectionName string, groupName string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewUnassignRequest(f.getCommand(), map[string]string{
		"directory":          directoryId,
		"roleCollectionName": roleCollectionName,
		"group":              groupName,
		"origin":             origin,
	}))
}

func (f *securityRoleCollectionFacade) AssignGroupByGlobalaccount(ctx context.Context, roleCollectionName string, groupName string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewAssignRequest(f.getCommand(), map[string]string{
		"globalAccount":       f.cliClient.GetGlobalAccountSubdomain(),
		"roleCollectionName":  roleCollectionName,
		"group":               groupName,
		"origin":              origin,
		"createUserIfMissing": "true",
	}))
}

func (f *securityRoleCollectionFacade) UnassignGroupByGlobalaccount(ctx context.Context, roleCollectionName string, groupName string, origin string) (xsuaa_authz.UserReference, *CommandResponse, error) {
	return doExecute[xsuaa_authz.UserReference](f.cliClient, ctx, NewUnassignRequest(f.getCommand(), map[string]string{
		"globalAccount":      f.cliClient.GetGlobalAccountSubdomain(),
		"roleCollectionName": roleCollectionName,
		"group":              groupName,
		"origin":             origin,
	}))
}
