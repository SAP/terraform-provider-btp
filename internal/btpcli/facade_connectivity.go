package btpcli

func newConnectivityFacade(cliClient *v2Client) connectivityFacade {
	return connectivityFacade{
		DestinationCertificate: newConnectivityDestinationCertificatesFacade(cliClient),
		DestinationTrust:    newConnectivityDestinationTrustFacade(cliClient),
		DestinationFragment: newConnectivityDestinationFragmentFacade(cliClient),
	}
}

type connectivityFacade struct {
	DestinationCertificate connectivityDestinationCertificatesFacade
	DestinationTrust    connectivityDestinationTrustFacade
	DestinationFragment connectivityDestinationFragmentFacade
}
