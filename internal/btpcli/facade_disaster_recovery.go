package btpcli

func newDisasterRecoveryFacade(cliClient *v2Client) disasterRecoveryFacade {
	return disasterRecoveryFacade{
		SubaccountPair: newDisasterRecoverySubaccountPairFacade(cliClient),
	}
}

type disasterRecoveryFacade struct {
	SubaccountPair disasterRecoverySubaccountPairFacade
}
