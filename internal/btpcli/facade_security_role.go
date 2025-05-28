package btpcli

import (
	"context"

	"github.com/SAP/terraform-provider-btp/internal/tfutils"

	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
)

func newSecurityRoleFacade(cliClient *v2Client) securityRoleFacade {
	return securityRoleFacade{cliClient: cliClient}
}

type securityRoleFacade struct {
	cliClient *v2Client
}

func (f *securityRoleFacade) getCommand() string {
	return "security/role"
}

func (f *securityRoleFacade) ListByGlobalAccount(ctx context.Context) ([]xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.Role](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"globalAccount": f.cliClient.GetGlobalAccountSubdomain(),
	}))
}

func (f *securityRoleFacade) GetByGlobalAccount(ctx context.Context, roleName string, roleTemplateAppId string, roleTemplateName string) (xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"globalAccount":    f.cliClient.GetGlobalAccountSubdomain(),
		"roleName":         roleName,
		"appId":            roleTemplateAppId,
		"roleTemplateName": roleTemplateName,
	}))
}

func (f *securityRoleFacade) ListBySubaccount(ctx context.Context, subaccountId string) ([]xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.Role](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"subaccount": subaccountId,
	}))
}

func (f *securityRoleFacade) GetBySubaccount(ctx context.Context, subaccountId string, roleName string, roleTemplateAppId string, roleTemplateName string) (xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"subaccount":       subaccountId,
		"roleName":         roleName,
		"appId":            roleTemplateAppId,
		"roleTemplateName": roleTemplateName,
	}))
}

func (f *securityRoleFacade) ListByDirectory(ctx context.Context, directoryId string) ([]xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[[]xsuaa_authz.Role](f.cliClient, ctx, NewListRequest(f.getCommand(), map[string]string{
		"directory": directoryId,
	}))
}

func (f *securityRoleFacade) GetByDirectory(ctx context.Context, directoryId string, roleName string, roleTemplateAppId string, roleTemplateName string) (xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewGetRequest(f.getCommand(), map[string]string{
		"directory":        directoryId,
		"roleName":         roleName,
		"appId":            roleTemplateAppId,
		"roleTemplateName": roleTemplateName,
	}))
}

type DirectoryRoleCreateInput struct {
	RoleName         string `btpcli:"roleName"`
	AppId            string `btpcli:"appId"`
	RoleTemplateName string `btpcli:"roleTemplateName"`
	DirectoryId      string `btpcli:"directory"`
	Description      string `btpcli:"description"`
}

func (f *securityRoleFacade) CreateByDirectory(ctx context.Context, args *DirectoryRoleCreateInput) (xsuaa_authz.Role, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_authz.Role{}, CommandResponse{}, err
	}

	_, exist := params["description"]
	if !exist {
		params["description"] = ""
	}

	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *securityRoleFacade) DeleteByDirectory(ctx context.Context, directoryId string, roleName string, roleTemplateAppId string, roleTemplateName string) (xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"directory":        directoryId,
		"roleName":         roleName,
		"appId":            roleTemplateAppId,
		"roleTemplateName": roleTemplateName,
	}))
}

type SubaccountRoleCreateInput struct {
	RoleName         string `btpcli:"roleName"`
	AppId            string `btpcli:"appId"`
	RoleTemplateName string `btpcli:"roleTemplateName"`
	SubaccountId     string `btpcli:"subaccount"`
	Description      string `btpcli:"description"`
	AttributeList    string `btpcli:"attributeList"`
}

func (f *securityRoleFacade) CreateBySubaccount(ctx context.Context, args *SubaccountRoleCreateInput) (xsuaa_authz.Role, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_authz.Role{}, CommandResponse{}, err
	}

	_, exist := params["description"]
	if !exist {
		params["description"] = ""
	}

	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *securityRoleFacade) DeleteBySubaccount(ctx context.Context, subaccountId string, roleName string, roleTemplateAppId string, roleTemplateName string) (xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"subaccount":       subaccountId,
		"roleName":         roleName,
		"appId":            roleTemplateAppId,
		"roleTemplateName": roleTemplateName,
	}))
}

type GlobalAccountRoleCreateInput struct {
	RoleName         string `btpcli:"roleName"`
	AppId            string `btpcli:"appId"`
	RoleTemplateName string `btpcli:"roleTemplateName"`
	Description      string `btpcli:"description"`
}

func (f *securityRoleFacade) CreateByGlobalAccount(ctx context.Context, args *GlobalAccountRoleCreateInput) (xsuaa_authz.Role, CommandResponse, error) {
	params, err := tfutils.ToBTPCLIParamsMap(args)

	if err != nil {
		return xsuaa_authz.Role{}, CommandResponse{}, err
	}

	params["globalAccount"] = f.cliClient.GetGlobalAccountSubdomain()

	_, exist := params["description"]
	if !exist {
		params["description"] = ""
	}

	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewCreateRequest(f.getCommand(), params))
}

func (f *securityRoleFacade) DeleteByGlobalAccount(ctx context.Context, roleName string, roleTemplateAppId string, roleTemplateName string) (xsuaa_authz.Role, CommandResponse, error) {
	return doExecute[xsuaa_authz.Role](f.cliClient, ctx, NewDeleteRequest(f.getCommand(), map[string]string{
		"globalAccount":    f.cliClient.GetGlobalAccountSubdomain(),
		"roleName":         roleName,
		"appId":            roleTemplateAppId,
		"roleTemplateName": roleTemplateName,
	}))
}

func (f *securityRoleFacade) AddBySubaccount(ctx context.Context, subaccountId string, targetRoleCollection string, roleName string, roleTemplateAppId string, roleTemplateName string) (CommandResponse, error) {
	return f.cliClient.Execute(ctx, NewAddRequest(f.getCommand(), map[string]string{
		"subaccount":         subaccountId,
		"roleName":           roleName,
		"roleCollectionName": targetRoleCollection,
		"roleTemplateAppID":  roleTemplateAppId,
		"roleTemplateName":   roleTemplateName,
	}))
}

func (f *securityRoleFacade) AddByDirectory(ctx context.Context, directoryId string, targetRoleCollection string, roleName string, roleTemplateAppId string, roleTemplateName string) (CommandResponse, error) {
	return f.cliClient.Execute(ctx, NewAddRequest(f.getCommand(), map[string]string{
		"directory":          directoryId,
		"roleName":           roleName,
		"roleCollectionName": targetRoleCollection,
		"roleTemplateAppID":  roleTemplateAppId,
		"roleTemplateName":   roleTemplateName,
	}))
}

func (f *securityRoleFacade) AddByGlobalAccount(ctx context.Context, targetRoleCollection string, roleName string, roleTemplateAppId string, roleTemplateName string) (CommandResponse, error) {
	return f.cliClient.Execute(ctx, NewAddRequest(f.getCommand(), map[string]string{
		"globalAccount":      f.cliClient.GetGlobalAccountSubdomain(),
		"roleName":           roleName,
		"roleCollectionName": targetRoleCollection,
		"roleTemplateAppID":  roleTemplateAppId,
		"roleTemplateName":   roleTemplateName,
	}))
}

func (f *securityRoleFacade) RemoveBySubaccount(ctx context.Context, subaccountId string, targetRoleCollection string, roleName string, roleTemplateAppId string, roleTemplateName string) (CommandResponse, error) {
	return f.cliClient.Execute(ctx, NewRemoveRequest(f.getCommand(), map[string]string{
		"subaccount":         subaccountId,
		"roleName":           roleName,
		"roleCollectionName": targetRoleCollection,
		"roleTemplateAppID":  roleTemplateAppId,
		"roleTemplateName":   roleTemplateName,
	}))
}

func (f *securityRoleFacade) RemoveByDirectory(ctx context.Context, directoryId string, targetRoleCollection string, roleName string, roleTemplateAppId string, roleTemplateName string) (CommandResponse, error) {
	return f.cliClient.Execute(ctx, NewRemoveRequest(f.getCommand(), map[string]string{
		"directory":          directoryId,
		"roleName":           roleName,
		"roleCollectionName": targetRoleCollection,
		"roleTemplateAppID":  roleTemplateAppId,
		"roleTemplateName":   roleTemplateName,
	}))
}

func (f *securityRoleFacade) RemoveByGlobalAccount(ctx context.Context, targetRoleCollection string, roleName string, roleTemplateAppId string, roleTemplateName string) (CommandResponse, error) {
	return f.cliClient.Execute(ctx, NewRemoveRequest(f.getCommand(), map[string]string{
		"globalAccount":      f.cliClient.GetGlobalAccountSubdomain(),
		"roleName":           roleName,
		"roleCollectionName": targetRoleCollection,
		"roleTemplateAppID":  roleTemplateAppId,
		"roleTemplateName":   roleTemplateName,
	}))
}
