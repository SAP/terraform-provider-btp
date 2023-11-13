package btpcli

func newSecurityFacade(cliClient *v2Client) securityFacade {
	return securityFacade{
		App:            newSecurityAppFacade(cliClient),
		Role:           newSecurityRoleFacade(cliClient),
		RoleCollection: newSecurityRoleCollectionFacade(cliClient),
		Settings:       newSecuritySettingsFacade(cliClient),
		Trust:          newSecurityTrustFacade(cliClient),
		User:           newSecurityUserFacade(cliClient),
	}
}

type securityFacade struct {
	App            securityAppFacade
	Role           securityRoleFacade
	RoleCollection securityRoleCollectionFacade
	Settings       securitySettingsFacade
	Trust          securityTrustFacade
	User           securityUserFacade
}
