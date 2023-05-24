package btpcli

func newServicesFacade(cliClient *v2Client) servicesFacade {
	return servicesFacade{
		Binding:  newServicesBindingFacade(cliClient),
		Broker:   newServicesBrokerFacade(cliClient),
		Instance: newServicesInstanceFacade(cliClient),
		Offering: newServicesOfferingFacade(cliClient),
		Plan:     newServicesPlanFacade(cliClient),
		Platform: newServicesPlatformFacade(cliClient),
	}
}

type servicesFacade struct {
	Binding  servicesBindingFacade
	Broker   servicesBrokerFacade
	Instance servicesInstanceFacade
	Offering servicesOfferingFacade
	Plan     servicesPlanFacade
	Platform servicesPlatformFacade
}
