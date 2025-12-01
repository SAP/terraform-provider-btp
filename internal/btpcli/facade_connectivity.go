package btpcli

func newConnectivityFacade(cliClient *v2Client) connectivityFacade {
	return connectivityFacade{
		DestinationTrust:    newConnectivityDestinationTrustFacade(cliClient),
	}
}

type connectivityFacade struct {
	DestinationTrust    connectivityDestinationTrustFacade
}
