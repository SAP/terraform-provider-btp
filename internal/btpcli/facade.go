package btpcli

func NewClientFacade(cliClient *v2Client) *ClientFacade {
	return &ClientFacade{
		v2Client:         cliClient,
		Accounts:         newAccountsFacade(cliClient),
		Services:         newServicesFacade(cliClient),
		Security:         newSecurityFacade(cliClient),
		Connectivity:     newConnectivityFacade(cliClient),
		DisasterRecovery: newDisasterRecoveryFacade(cliClient),
	}
}

type ClientFacade struct {
	*v2Client
	Accounts         accountsFacade
	Services         servicesFacade
	Security         securityFacade
	Connectivity     connectivityFacade
	DisasterRecovery disasterRecoveryFacade
}
