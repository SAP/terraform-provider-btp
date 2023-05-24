package btpcli

func newAccountsFacade(cliClient *v2Client) accountsFacade {
	return accountsFacade{
		AvailableEnvironment: newAccountsAvailableEnvironmentFacade(cliClient),
		AvailableRegion:      newAccountsAvailableRegionFacade(cliClient),
		Directory:            newAccountsDirectoryFacade(cliClient),
		Entitlement:          newAccountsEntitlementFacade(cliClient),
		EnvironmentInstance:  newAccountsEnvironmentInstanceFacade(cliClient),
		GlobalAccount:        newAccountsGlobalAccountFacade(cliClient),
		Label:                newAccountsLabelFacade(cliClient),
		ResourceProvider:     newAccountsResourceProviderFacade(cliClient),
		Subaccount:           newAccountsSubaccountFacade(cliClient),
		Subscription:         newAccountsSubscriptionFacade(cliClient),
	}
}

type accountsFacade struct {
	AvailableEnvironment accountsAvailableEnvironmentFacade
	AvailableRegion      accountsAvailableRegionFacade
	Directory            accountsDirectoryFacade
	Entitlement          accountsEntitlementFacade
	EnvironmentInstance  accountsEnvironmentInstanceFacade
	GlobalAccount        accountsGlobalAccountFacade
	Label                accountsLabelFacade
	ResourceProvider     accountsResourceProviderFacade
	Subaccount           accountsSubaccountFacade
	Subscription         accountsSubscriptionFacade
}
