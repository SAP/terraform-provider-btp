package btpcli

func newConnectivityFacade(cliClient *v2Client) connectivityFacade {
	return connectivityFacade{
		DestinationCertificate: newConnectivityDestinationCertificatesFacade(cliClient),
	}
}

type connectivityFacade struct {
	DestinationCertificate connectivityDestinationCertificatesFacade
}
