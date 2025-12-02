package btpcli

func newConnectivityFacade(cliClient *v2Client) connectivityFacade {
	return connectivityFacade{
		DestinationTrust:    newConnectivityDestinationTrustFacade(cliClient),
		DestinationFragment: newConnectivityDestinationFragmentFacade(cliClient),
	}
}

type connectivityFacade struct {
	DestinationTrust    connectivityDestinationTrustFacade
	DestinationFragment connectivityDestinationFragmentFacade
}
